package plugins

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

type SourceNewExecutionClientFunc func(context.Context, *SourcePlugin, specs.Source) (schema.ClientMeta, error)

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
	// exampleConfig
	exampleConfig string
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger
}

type SourceOption func(*SourcePlugin)

const minGoRoutines = 5

func WithSourceExampleConfig(exampleConfig string) SourceOption {
	return func(p *SourcePlugin) {
		p.exampleConfig = exampleConfig
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

// Add internal columns
func addInternalColumns(tables []*schema.Table) {
	for _, table := range tables {
		cqId := schema.CqIdColumn
		if len(table.PrimaryKeys()) == 0 {
			cqId.CreationOptions.PrimaryKey = true
		}
		table.Columns = append(table.Columns, cqId, schema.CqFetchTime)
		addInternalColumns(table.Relations)
	}
}

func NewSourcePlugin(name string, version string, tables []*schema.Table, newExecutionClient SourceNewExecutionClientFunc, opts ...SourceOption) *SourcePlugin {
	p := SourcePlugin{
		name:               name,
		version:            version,
		tables:             tables,
		newExecutionClient: newExecutionClient,
	}
	for _, opt := range opts {
		opt(&p)
	}
	addInternalColumns(p.tables)
	if err := p.validate(); err != nil {
		panic(err)
	}
	return &p
}

func (p *SourcePlugin) validate() error {
	if p.newExecutionClient == nil {
		return fmt.Errorf("newExecutionClient function not defined for source plugin: " + p.name)
	}

	if err := p.tables.ValidateDuplicateColumns(); err != nil {
		return fmt.Errorf("found duplicate columns in source plugin: %s: %w", p.name, err)
	}

	if err := p.tables.ValidateDuplicateTables(); err != nil {
		return fmt.Errorf("found duplicate tables in source plugin: %s: %w", p.name, err)
	}
	return nil
}

func (p *SourcePlugin) Tables() schema.Tables {
	return p.tables
}

func (p *SourcePlugin) ExampleConfig() string {
	return p.exampleConfig
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

// Sync data from source to the given channel
func (p *SourcePlugin) Sync(ctx context.Context, spec specs.Source, res chan<- *schema.Resource) error {
	c, err := p.newExecutionClient(ctx, p, spec)
	if err != nil {
		return fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
	}

	// limiter used to limit the amount of resources fetched concurrently
	maxGoroutines := spec.MaxGoRoutines
	if maxGoroutines < minGoRoutines {
		maxGoroutines = minGoRoutines
	}
	p.logger.Info().Uint64("max_goroutines", maxGoroutines).Msg("starting fetch")

	// goroutinesSem := semaphore.NewWeighted(helpers.Uint64ToInt64(maxGoroutines))

	w := sync.WaitGroup{}
	totalResources := 0
	startTime := time.Now()
	tableNames, err := p.interpolateAllResources(p.tables.TableNames())
	if err != nil {
		return err
	}

	// this is the same fetchtime for all resources
	fetchTime := time.Now()

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
		// totalClients := len(clients)
		// newN, err := helpers.TryAcquireMax(ctx, goroutinesSem, int64(totalClients))
		// if err != nil {
		// 	p.logger.Error().Err(err).Msg("failed to TryAcquireMax semaphore. exiting")
		// 	break
		// }
		// goroutinesSem.TryAcquire()
		w.Add(1)
		go func() {
			defer w.Done()
			wg := sync.WaitGroup{}
			p.logger.Info().Str("table", table.Name).Msg("fetch start")
			tableStartTime := time.Now()
			totalTableResources := 0
			for _, client := range clients {
				client := client

				// i := i
				wg.Add(1)
				go func() {
					defer wg.Done()
					// defer goroutinesSem.Release(1)
					totalTableResources += table.Resolve(ctx, client, fetchTime, nil, res)
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
