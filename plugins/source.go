package plugins

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/helpers"
	"github.com/cloudquery/plugin-sdk/helpers/limit"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/sync/semaphore"

	_ "embed"
)

//go:embed source_schema.json
var sourceSchema string

const ExampleSourceConfig = `
# max_goroutines to use when fetching. 0 means default and calculated by CloudQuery
# max_goroutines: 0
# By default cloudquery will fetch all tables in the source plugin
# tables: ["*"]
# skip_tables specify which tables to skip. especially useful when using "*" for tables
# skip_tables: []
`

type SourceConfigureFunc func(context.Context, *SourcePlugin, specs.SourceSpec) (schema.ClientMeta, error)

// SourcePlugin is the base structure required to pass to sdk.serve
// We take a similar/declerative approach to API here similar to Cobra
type SourcePlugin struct {
	// Name of plugin i.e aws,gcp, azure etc'
	Name string
	// Version of the plugin
	Version string
	// Classify error and return it's severity and type
	ClassifyError schema.ClassifyErrorFunc
	// Called upon configure call to validate and init configuration
	configure SourceConfigureFunc
	// Tables is all tables supported by this source plugin
	Tables schema.Tables
	// JsonSchema for specific source plugin spec
	JsonSchema string
	// ExampleConfig is the example configuration for this plugin
	ExampleConfig string
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	Logger zerolog.Logger

	// Internal fields set by configure
	clientMeta schema.ClientMeta
	spec       *specs.SourceSpec
	m          sync.Mutex
}

type SourceOption func(*SourcePlugin)

func WithSourceExampleConfig(exampleConfig string) SourceOption {
	return func(p *SourcePlugin) {
		p.ExampleConfig = exampleConfig
	}
}

func WithSourceJsonSchema(jsonSchema string) SourceOption {
	return func(p *SourcePlugin) {
		p.JsonSchema = jsonSchema
	}
}

func WithSourceLogger(logger zerolog.Logger) SourceOption {
	return func(p *SourcePlugin) {
		p.Logger = logger
	}
}

func WithClassifyError(classifyError schema.ClassifyErrorFunc) SourceOption {
	return func(p *SourcePlugin) {
		p.ClassifyError = classifyError
	}
}

func NewSourcePlugin(name string, version string, tables []*schema.Table, configure SourceConfigureFunc, opts ...SourceOption) *SourcePlugin {
	p := SourcePlugin{
		Name:      name,
		Version:   version,
		Tables:    tables,
		configure: configure,
	}
	if configure == nil {
		panic("configure function not defined for source plugin:" + name)
	}
	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

func (p *SourcePlugin) Configure(ctx context.Context, spec specs.SourceSpec) (*gojsonschema.Result, error) {
	// we permit only one configure per source plugin for security reasons.
	// in the grpc layer this will behave similarly and for every new configuration/creds the cli will have to create a new process.
	p.m.Lock()
	defer p.m.Unlock()
	if p.spec != nil {
		return nil, fmt.Errorf("source plugin %s already configured", p.Name)
	}
	res, err := specs.ValidateSpec(sourceSchema, spec)
	if err != nil {
		return nil, err
	}
	if !res.Valid() {
		return res, nil
	}

	// if resources ["*"] is requested we will fetch all resources
	p.spec.Tables, err = p.interpolateAllResources(p.spec.Tables)
	if err != nil {
		return res, fmt.Errorf("failed to interpolate resources: %w", err)
	}

	p.clientMeta, err = p.configure(ctx, p, spec)
	if err != nil {
		return res, fmt.Errorf("failed to configure source plugin: %w", err)
	}
	p.spec = &spec
	return res, nil
}

// Fetch fetches data according to source configuration and
func (p *SourcePlugin) Fetch(ctx context.Context, res chan<- *schema.Resource) error {
	if p.spec == nil {
		return fmt.Errorf("source plugin not configured")
	}

	// limiter used to limit the amount of resources fetched concurrently
	maxGoroutines := p.spec.MaxGoRoutines
	if maxGoroutines == 0 {
		maxGoroutines = limit.GetMaxGoRoutines()
	}
	p.Logger.Info().Uint64("max_goroutines", maxGoroutines).Msg("starting fetch")
	goroutinesSem := semaphore.NewWeighted(helpers.Uint64ToInt64(maxGoroutines))

	w := sync.WaitGroup{}
	totalResources := 0
	startTime := time.Now()
	tableNames := p.Tables.TableNames()
	for _, table := range p.Tables {
		table := table
		if funk.ContainsString(p.spec.SkipTables, table.Name) || !funk.ContainsString(tableNames, table.Name) {
			p.Logger.Info().Str("table", table.Name).Msg("skipping table")
			continue
		}
		clients := []schema.ClientMeta{p.clientMeta}
		if table.Multiplex != nil {
			clients = table.Multiplex(p.clientMeta)
		}
		// because table can't import sourceplugin we need to set classifyError if it is not set by table
		if table.ClassifyError == nil {
			table.ClassifyError = p.ClassifyError
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
			p.Logger.Info().Str("table", table.Name).Msg("fetch start")
			tableStartTime := time.Now()
			totalTableResources := 0
			for i, client := range clients {
				client := client
				i := i
				// acquire semaphore only if we couldn't acquire it earlier
				if newN > 0 && i >= (totalClients-int(newN)) {
					if err := goroutinesSem.Acquire(ctx, 1); err != nil {
						// this can happen if context was cancelled so we just break out of the loop
						p.Logger.Error().Err(err).Msg("failed to acquire semaphore")
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
			p.Logger.Info().Str("table", table.Name).Int("total_resources", totalTableResources).TimeDiff("duration", time.Now(), tableStartTime).Msg("fetch table finished")
		}()
	}
	w.Wait()
	p.Logger.Info().Int("total_resources", totalResources).TimeDiff("duration", time.Now(), startTime).Msg("fetch finished")
	return nil
}

func (p *SourcePlugin) interpolateAllResources(tables []string) ([]string, error) {
	if !funk.ContainsString(tables, "*") {
		return tables, nil
	}

	if len(tables) > 1 {
		return nil, fmt.Errorf("invalid \"*\" resource, with explicit resources")
	}

	allResources := make([]string, 0, len(p.Tables))
	for _, k := range p.Tables {
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
