package plugins

import (
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
)

type SourceNewExecutionClientFunc func(context.Context, zerolog.Logger, specs.Source) (schema.ClientMeta, error)

// SourcePlugin is the base structure required to pass to sdk.serve
// We take a similar/declerative approach to API here similar to Cobra
type SourcePlugin struct {
	// Name of plugin i.e aws,gcp, azure etc'
	name string
	// Version of the plugin
	version string
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

const (
	defaultMaxGoRoutines = 500000
)

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

// Add internal columns
func addInternalColumns(tables []*schema.Table) {
	for _, table := range tables {
		cqID := schema.CqIDColumn
		if len(table.PrimaryKeys()) == 0 {
			cqID.CreationOptions.PrimaryKey = true
		}
		table.Columns = append(table.Columns, cqID, schema.CqFetchTime)
		addInternalColumns(table.Relations)
	}
}

// NewSourcePlugin returns a new plugin with a given name, version, tables, newExecutionClient
// and additional options.
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

// Tables returns all supported tables by this source plugin
func (p *SourcePlugin) Tables() schema.Tables {
	return p.tables
}

// ExampleConfig returns an example configuration for this source plugin
func (p *SourcePlugin) ExampleConfig() string {
	return p.exampleConfig
}

// Name return the name of this plugin
func (p *SourcePlugin) Name() string {
	return p.name
}

// Version returns the version of this plugin
func (p *SourcePlugin) Version() string {
	return p.version
}

// SetLogger sets the logger for this plugin which will be used in Sync and all other function calls.
func (p *SourcePlugin) SetLogger(log zerolog.Logger) {
	p.logger = log
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *SourcePlugin) Sync(ctx context.Context, spec specs.Source, res chan<- *schema.Resource) error {
	c, err := p.newExecutionClient(ctx, p.logger, spec)
	if err != nil {
		return fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
	}

	// limiter used to limit the amount of resources fetched concurrently
	concurrency := spec.Concurrency
	if concurrency == 0 {
		concurrency = defaultMaxGoRoutines
	}
	p.logger.Info().Uint64("concurrency", concurrency).Msg("starting fetch")
	goroutinesSem := semaphore.NewWeighted(helpers.Uint64ToInt64(concurrency))
	wg := sync.WaitGroup{}
	totalResources := 0
	startTime := time.Now()
	tableNames, err := p.interpolateAllResources(spec.Tables)
	if err != nil {
		return err
	}

	for _, table := range p.tables {
		table := table
		if funk.ContainsString(spec.SkipTables, table.Name) || !funk.ContainsString(tableNames, table.Name) {
			p.logger.Debug().Str("table", table.Name).Msg("skipping table")
			continue
		}
		clients := []schema.ClientMeta{c}
		if table.Multiplex != nil {
			clients = table.Multiplex(c)
		}
		for _, client := range clients {
			client := client
			wg.Add(1)
			if err := goroutinesSem.Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				return err
			}
			go func() {
				defer wg.Done()
				defer goroutinesSem.Release(1)
				// TODO: prob introduce client.Identify() to be used in logs
				tableStartTime := time.Now()
				p.logger.Info().Str("table", table.Name).Msg("fetch start")
				totalTableResources := table.Resolve(ctx, client, startTime, nil, res)
				totalResources += totalTableResources
				p.logger.Info().Str("table", table.Name).Int("total_resources", totalTableResources).TimeDiff("duration", time.Now(), tableStartTime).Msg("fetch table finished")
			}()
		}
	}
	wg.Wait()
	p.logger.Info().Int("total_resources", totalResources).TimeDiff("duration", time.Now(), startTime).Msg("fetch finished")
	return nil
}

func (p *SourcePlugin) interpolateAllResources(tables []string) ([]string, error) {
	if tables == nil {
		return make([]string, 0), nil
	}

	if funk.Equal(tables, []string{"*"}) {
		allResources := make([]string, 0, len(p.tables))
		for _, k := range p.tables {
			allResources = append(allResources, k.Name)
		}
		return allResources, nil
	}

	if funk.ContainsString(tables, "*") {
		return nil, fmt.Errorf("invalid \"*\" resource, with explicit resources")
	}

	return tables, nil
}
