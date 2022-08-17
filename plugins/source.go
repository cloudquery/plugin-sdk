package plugins

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/helpers"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/semaphore"
	"gopkg.in/yaml.v3"

	_ "embed"
)

const ExampleSourceConfig = `
# max_goroutines to use when fetching. 0 means default and calculated by CloudQuery
# max_goroutines: 0
# By default cloudquery will fetch all tables in the source plugin
# tables: ["*"]
# skip_tables specify which tables to skip. especially useful when using "*" for tables
# skip_tables: []
`

const sourcePluginExampleConfigTemplate = `kind: source
spec:
  name: {{.Name}}
  version: {{.Version}}
  # path: Path to the plugin. by default it is the same as the name of the plugin.
  # registry can be local, github, grpc (default github)
  # registry: github
  # max_goroutines used for sync, by default calculated automatically depending on
  # memory and cpu avaialble
  # max_goroutines: 0
  tables: ["*"]
  # skip_tables is useful if you want to fetch all tables apart from specific ones
  # skip_tables: []
  # name of destinations to sync the data
  destinations: []
  configuration:
  {{.PluginExampleConfig | indent 4}}
`

type SourceNewExecutionClientFunc func(context.Context, *SourcePlugin, specs.Source) (schema.ClientMeta, error)

type SourceNewSpecFunc func() interface{}

// SourcePlugin is the base structure required to pass to sdk.serve
// We take a similar/declerative approach to API here similar to Cobra
type SourcePlugin struct {
	// Name of plugin i.e aws,gcp, azure etc'
	name string
	// Version of the plugin
	version string
	// Classify error and return it's severity and type
	ignoreError schema.IgnoreErrorFunc
	// Called upon configure call to validate and init configuration
	newExecutionClient SourceNewExecutionClientFunc
	// Tables is all tables supported by this source plugin
	tables schema.Tables
	// newSpec return a new struct to be pupolated by the passed configuration
	newSpec SourceNewSpecFunc
	// JsonSchema for specific source plugin spec
	jsonSchema string
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger
}

type SourceOption func(*SourcePlugin)

func WithSourceJsonSchema(jsonSchema string) SourceOption {
	return func(p *SourcePlugin) {
		p.jsonSchema = jsonSchema
	}
}

func WithSourceLogger(logger zerolog.Logger) SourceOption {
	return func(p *SourcePlugin) {
		p.logger = logger
	}
}

func WithClassifyError(ignoreError schema.IgnoreErrorFunc) SourceOption {
	return func(p *SourcePlugin) {
		p.ignoreError = ignoreError
	}
}

func NewSourcePlugin(name string, version string, tables []*schema.Table, newExecutionClient SourceNewExecutionClientFunc, newSpec SourceNewSpecFunc, opts ...SourceOption) *SourcePlugin {
	p := SourcePlugin{
		name:               name,
		version:            version,
		tables:             tables,
		newExecutionClient: newExecutionClient,
		newSpec:            newSpec,
	}
	if newExecutionClient == nil {
		panic("newExecutionClient function not defined for source plugin:" + name)
	}
	if newSpec == nil {
		panic("newConfig function not defined for source plugin:" + name)
	}
	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

func (p *SourcePlugin) Tables() schema.Tables {
	return p.tables
}

func (p *SourcePlugin) ExampleConfig() (string, error) {
	spec := specs.Spec{
		Kind: specs.KindSource,
		Spec: specs.Source{
			Name:         p.name,
			Version:      p.version,
			Tables:       []string{"*"},
			Destinations: []string{},
			Spec:         p.newSpec(),
		},
	}
	bytes := bytes.NewBuffer([]byte(""))
	enc := yaml.NewEncoder(bytes)
	enc.SetIndent(2)
	if err := enc.Encode(spec); err != nil {
		return "", err
	}
	return bytes.String(), nil
}

func (p *SourcePlugin) GetJsonSchema() string {
	return p.jsonSchema
}

func (p *SourcePlugin) Name() string {
	return p.name
}

func (p *SourcePlugin) Version() string {
	return p.version
}

func (p *SourcePlugin) SetLogger(log zerolog.Logger) {
	p.logger = log
}

const minGoRoutines = 5

// Sync data from source to the given channel
func (p *SourcePlugin) Sync(ctx context.Context, spec specs.Source, res chan<- *schema.Resource) error {
	c, err := p.newExecutionClient(ctx, p, spec)
	if err != nil {
		return fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
	}

	// limiter used to limit the amount of resources fetched concurrently
	maxGoroutines := spec.MaxGoRoutines
	if maxGoroutines == 0 {
		maxGoroutines = minGoRoutines
	}
	p.logger.Info().Uint64("max_goroutines", maxGoroutines).Msg("starting fetch")

	goroutinesSem := semaphore.NewWeighted(helpers.Uint64ToInt64(maxGoroutines))

	w := sync.WaitGroup{}
	totalResources := 0
	startTime := time.Now()
	tableNames, err := p.interpolateAllResources(p.tables.TableNames())
	if err != nil {
		return err
	}
	for _, table := range p.tables {
		table := table
		if funk.ContainsString(spec.SkipTables, table.Name) || !funk.ContainsString(tableNames, table.Name) {
			p.logger.Info().Str("table", table.Name).Msg("skipping table")
			continue
		}
		clients := []schema.ClientMeta{c}
		if table.Multiplex != nil {
			clients = table.Multiplex(c)
		}
		// because table can't import sourceplugin we need to set classifyError if it is not set by table
		if table.IgnoreError == nil {
			table.IgnoreError = p.ignoreError
		}
		// we call this here because we dont know when the following goroutine will be called and we do want an order
		// of table by table
		totalClients := len(clients)
		newN := helpers.TryAcquireMax(goroutinesSem, int64(totalClients))
		// goroutinesSem.TryAcquire()
		w.Add(1)
		go func() {
			defer w.Done()
			defer goroutinesSem.Release(int64(totalClients) - newN)
			wg := sync.WaitGroup{}
			p.logger.Info().Str("table", table.Name).Msg("fetch start")
			tableStartTime := time.Now()
			totalTableResources := 0
			for i, client := range clients {
				client := client
				i := i
				// acquire semaphore only if we couldn't acquire it earlier
				if newN > 0 && i >= (totalClients-int(newN)) {
					if err := goroutinesSem.Acquire(ctx, 1); err != nil {
						// this can happen if context was cancelled so we just break out of the loop
						p.logger.Error().Err(err).Msg("failed to acquire semaphore")
						return
					}
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					if newN > 0 && i >= (totalClients-int(newN)) {
						defer goroutinesSem.Release(1)
					}

					totalTableResources += table.Resolve(ctx, client, nil, res)
				}()
			}
			wg.Wait()
			totalResources += totalTableResources
			p.logger.Info().Str("table", table.Name).Int("total_resources", totalTableResources).TimeDiff("duration", time.Now(), tableStartTime).Msg("fetch table finished")
		}()
	}
	w.Wait()
	p.logger.Info().Int("total_resources", totalResources).TimeDiff("duration", time.Now(), startTime).Msg("fetch finished")
	return nil
}

func (p *SourcePlugin) interpolateAllResources(tables []string) ([]string, error) {
	if tables == nil {
		return make([]string, 0), nil
	}
	if !funk.ContainsString(tables, "*") {
		return tables, nil
	}

	if len(tables) > 1 {
		return nil, fmt.Errorf("invalid \"*\" resource, with explicit resources")
	}

	allResources := make([]string, 0, len(p.tables))
	for _, k := range p.tables {
		allResources = append(allResources, k.Name)
	}
	return allResources, nil
}

// func getTableDuplicates(resource string, table *schema.Table, tableNames map[string]string) error {
// 	for _, r := range table.Relations {
// 		if err := getTableDuplicates(resource, r, tableNames); err != nil {
// 			return err
// 		}
// 	}
// 	if existing, ok := tableNames[table.Name]; ok {
// 		return fmt.Errorf("table name %s used more than once, duplicates are in %s and %s", table.Name, existing, resource)
// 	}
// 	tableNames[table.Name] = resource
// 	return nil
// }
