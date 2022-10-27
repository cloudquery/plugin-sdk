package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
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
	// maxDepth is the max depth of tables
	maxDepth uint64
}

const (
	maxAllowedDepth = 3
)

func (s *TableClientStats) Equal(other *TableClientStats) bool {
	return s.Resources == other.Resources && s.Errors == other.Errors && s.Panics == other.Panics
}

// Equal compares to stats. Mostly useful in testing
func (s *SourceStats) Equal(other *SourceStats) bool {
	for table, clientStats := range s.TableClient {
		for client, stats := range clientStats {
			if _, ok := other.TableClient[table]; !ok {
				return false
			}
			if _, ok := other.TableClient[table][client]; !ok {
				return false
			}
			if !stats.Equal(other.TableClient[table][client]) {
				return false
			}
		}
	}
	for table, clientStats := range other.TableClient {
		for client, stats := range clientStats {
			if _, ok := s.TableClient[table]; !ok {
				return false
			}
			if _, ok := s.TableClient[table][client]; !ok {
				return false
			}
			if !stats.Equal(s.TableClient[table][client]) {
				return false
			}
		}
	}
	return true
}

func (s *SourceStats) initWithTables(tables schema.Tables) {
	for _, table := range tables {
		if _, ok := s.TableClient[table.Name]; !ok {
			s.TableClient[table.Name] = make(map[string]*TableClientStats)
		}
		s.initWithTables(table.Relations)
	}
}

func (s *SourceStats) initWithClients(table *schema.Table, clients []schema.ClientMeta) {
	for _, client := range clients {
		if _, ok := s.TableClient[table.Name][client.Name()]; !ok {
			s.TableClient[table.Name][client.Name()] = &TableClientStats{}
		}
		for _, relation := range table.Relations {
			s.initWithClients(relation, clients)
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
func setParents(tables schema.Tables, parent *schema.Table) {
	for _, table := range tables {
		table.Parent = parent
		setParents(table.Relations, table)
	}
}

func maxDepth(tables schema.Tables) uint64 {
	var depth uint64
	if len(tables) == 0 {
		return 0
	}
	for _, table := range tables {
		newDepth := 1 + maxDepth(table.Relations)
		if newDepth > depth {
			depth = newDepth
		}
	}
	return depth
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
	setParents(p.tables, nil)
	if err := p.validate(); err != nil {
		panic(err)
	}
	p.stats.initWithTables(p.tables)
	p.maxDepth = maxDepth(p.tables)
	if p.maxDepth > maxAllowedDepth {
		panic(fmt.Errorf("max depth of tables is %d, max allowed is %d", p.maxDepth, maxAllowedDepth))
	}

	return &p
}

func (p *SourcePlugin) SetLogger(logger zerolog.Logger) {
	p.logger = logger
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

func filterParentTables(tables schema.Tables, filter []string) schema.Tables {
	var res schema.Tables
	if tables == nil {
		return nil
	}
	if len(filter) == 0 {
		return tables
	}
	for _, name := range filter {
		if tables.Get(name) != nil {
			res = append(res, tables.Get(name))
		}
	}
	return res
}

// Sync is syncing data from the requested tables in spec to the given channel
func (p *SourcePlugin) Sync(ctx context.Context, spec specs.Source, res chan<- *schema.Resource) error {
	spec.SetDefaults()
	if err := spec.Validate(); err != nil {
		return fmt.Errorf("invalid spec: %w", err)
	}
	tableNames, err := p.listAndValidateTables(spec.Tables, spec.SkipTables)
	if err != nil {
		return err
	}
	tables := filterParentTables(p.tables, tableNames)

	c, err := p.newExecutionClient(ctx, p.logger, spec)
	if err != nil {
		return fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
	}
	startTime := time.Now()
	p.syncDfs(ctx, spec, c, tables, res)

	p.logger.Info().Interface("stats", p.stats).TimeDiff("duration", time.Now(), startTime).Msg("sync finished")
	return nil
}
