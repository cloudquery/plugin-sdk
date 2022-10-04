package plugins

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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
}

const (
	defaultConcurrency = 500000
)

// Add internal columns
func addInternalColumns(tables []*schema.Table) {
	for _, table := range tables {
		cqID := schema.CqIDColumn
		if len(table.PrimaryKeys()) == 0 {
			cqID.CreationOptions.PrimaryKey = true
		}
		table.Columns = append([]schema.Column{cqID, schema.CqParentIDColumn, schema.CqSourceName, schema.CqSyncTime}, table.Columns...)
		addInternalColumns(table.Relations)
	}
}

// Set parent links on relational tables
func setParents(tables []*schema.Table) {
	for _, table := range tables {
		for i := range table.Relations {
			table.Relations[i].Parent = table
		}
	}
}

// NewSourcePlugin returns a new plugin with a given name, version, tables, newExecutionClient
// and additional options.
func NewSourcePlugin(name string, version string, tables []*schema.Table, newExecutionClient SourceNewExecutionClientFunc) *SourcePlugin {
	p := SourcePlugin{
		name:               name,
		version:            version,
		tables:             tables,
		newExecutionClient: newExecutionClient,
	}
	addInternalColumns(p.tables)
	setParents(p.tables)
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

// Name return the name of this plugin
func (p *SourcePlugin) Name() string {
	return p.name
}

// Version returns the version of this plugin
func (p *SourcePlugin) Version() string {
	return p.version
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *SourcePlugin) Sync(ctx context.Context, logger zerolog.Logger, spec specs.Source, res chan<- *schema.Resource) (*schema.SyncSummary, error) {
	c, err := p.newExecutionClient(ctx, logger, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
	}

	// limiter used to limit the amount of resources fetched concurrently
	concurrency := spec.Concurrency
	if concurrency == 0 {
		concurrency = defaultConcurrency
	}
	logger.Info().Uint64("concurrency", concurrency).Msg("starting fetch")
	goroutinesSem := semaphore.NewWeighted(helpers.Uint64ToInt64(concurrency))
	wg := sync.WaitGroup{}
	summary := schema.SyncSummary{}
	startTime := time.Now()
	tableNames, err := p.interpolateAllResources(spec.Tables)
	if err != nil {
		return nil, err
	}

	for _, table := range p.tables {
		table := table
		if funk.ContainsString(spec.SkipTables, table.Name) || !funk.ContainsString(tableNames, table.Name) {
			logger.Debug().Str("table", table.Name).Msg("skipping table")
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
				return nil, err
			}
			go func() {
				defer wg.Done()
				defer goroutinesSem.Release(1)
				// TODO: prob introduce client.Identify() to be used in logs

				tableSummary := table.Resolve(ctx, client, nil, res)
				atomic.AddUint64(&summary.Resources, tableSummary.Resources)
				atomic.AddUint64(&summary.Errors, tableSummary.Errors)
				atomic.AddUint64(&summary.Panics, tableSummary.Panics)
			}()
		}
	}
	wg.Wait()
	logger.Info().Uint64("total_resources", summary.Resources).TimeDiff("duration", time.Now(), startTime).Msg("fetch finished")
	return &summary, nil
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
