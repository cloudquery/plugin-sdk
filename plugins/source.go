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

type TableClientStats struct {
	Resources uint64
	Errors    uint64
	Panics    uint64
	StartTime time.Time
	EndTime   time.Time
}

type SourceStats struct {
	TableClient map[string]map[string]*TableClientStats
}

// SourcePlugin is the base structure required to pass to sdk.serve
// We take a declarative approach to API here similar to Cobra
type SourcePlugin struct {
	// Name of plugin i.e aws,gcp, azure etc'
	name string
	// Version of the plugin
	version string
	// Called upon configure call to validate and init configuration
	newExecutionClient SourceNewExecutionClientFunc
	// Tables is all tables supported by this source plugin
	tables schema.Tables
	// status sync stats
	stats SourceStats
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger
	// resourceSem is a semaphore that limits the number of concurrent resources being fetched
	resourceSem *semaphore.Weighted
	// tableSem is a semaphore that limits the number of concurrent tables being fetched
	tableSem *semaphore.Weighted
}

func (s *SourceStats) initWithTables(tables schema.Tables) {
	for _, table := range tables {
		if _, ok := s.TableClient[table.Name]; !ok {
			s.TableClient[table.Name] = make(map[string]*TableClientStats)
		}
		s.initWithTables(table.Relations)
	}
}

func (s *SourceStats) initWithClients(clients []schema.ClientMeta) {
	for _, client := range clients {
		for table := range s.TableClient {
			if _, ok := s.TableClient[table][client.Name()]; !ok {
				s.TableClient[table][client.Name()] = &TableClientStats{}
			}
		}
	}
}

// Add internal columns
func addInternalColumns(tables []*schema.Table) {
	for _, table := range tables {
		cqID := schema.CqIDColumn
		if len(table.PrimaryKeys()) == 0 {
			cqID.CreationOptions.PrimaryKey = true
		}
		table.Columns = append([]schema.Column{cqID, schema.CqParentIDColumn}, table.Columns...)
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
		stats:              SourceStats{TableClient: make(map[string]map[string]*TableClientStats)},
	}
	addInternalColumns(p.tables)
	setParents(p.tables)
	if err := p.validate(); err != nil {
		panic(err)
	}
	p.stats.initWithTables(p.tables)
	return &p
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

func (p *SourcePlugin) Stats() SourceStats {
	return p.stats
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *SourcePlugin) Sync(ctx context.Context, logger zerolog.Logger, spec specs.Source, res chan<- *schema.Resource) error {
	spec.SetDefaults()
	if err := spec.Validate(); err != nil {
		return fmt.Errorf("invalid spec: %w", err)
	}
	tableNames, err := p.listAndValidateTables(spec.Tables, spec.SkipTables)
	if err != nil {
		return err
	}

	c, err := p.newExecutionClient(ctx, logger, spec)
	if err != nil {
		return fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
	}

	p.tableSem = semaphore.NewWeighted(helpers.Uint64ToInt64(spec.TableConcurrency))
	p.resourceSem = semaphore.NewWeighted(helpers.Uint64ToInt64(spec.ResourceConcurrency))
	wg := sync.WaitGroup{}
	startTime := time.Now()
	logger.Info().Interface("spec", spec).Strs("tables", tableNames).Msg("starting sync")
	for _, table := range p.tables {
		table := table
		if !funk.ContainsString(tableNames, table.Name) {
			logger.Debug().Str("table", table.Name).Msg("skipping table")
			continue
		}
		clients := []schema.ClientMeta{c}
		if table.Multiplex != nil {
			clients = table.Multiplex(c)
		}
		p.stats.initWithClients(clients)

		for _, client := range clients {
			client := client
			if err := p.tableSem.Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				wg.Wait()
				return err
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer p.tableSem.Release(1)
				// not checking for error here as nothing much todo.
				// the error is logged and this happens when context is cancelled
				p.resolveTable(ctx, table, client, nil, res)
			}()
		}
	}
	wg.Wait()
	logger.Info().Interface("stats", p.stats).TimeDiff("duration", time.Now(), startTime).Msg("sync finished")
	return nil
}


