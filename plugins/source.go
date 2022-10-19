package plugins

import (
	"context"
	"fmt"
	"strings"
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

	if err := p.tables.ValidateNames(); err != nil {
		return fmt.Errorf("found table with invalid name in source plugin: %s: %w", p.name, err)
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
	spec.SetDefaults()
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("invalid spec: %w", err)
	}
	tableNames, err := p.listAndValidateTables(spec.Tables, spec.SkipTables)
	if err != nil {
		return nil, err
	}

	c, err := p.newExecutionClient(ctx, logger, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
	}
	logger.Info().Interface("spec", spec).Msg("starting sync")
	tableSem := semaphore.NewWeighted(helpers.Uint64ToInt64(spec.TableConcurrency))
	wg := sync.WaitGroup{}
	resourceSem := semaphore.NewWeighted(helpers.Uint64ToInt64(spec.ResourceConcurrency))
	summary := schema.SyncSummary{}
	startTime := time.Now()

	logger.Debug().Interface("tables", tableNames).Msg("got table names")

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
		for _, client := range clients {
			client := client
			wg.Add(1)
			if err := tableSem.Acquire(ctx, 1); err != nil {
				// This means context was cancelled
				return nil, err
			}
			go func() {
				defer wg.Done()
				defer tableSem.Release(1)
				// TODO: prob introduce client.Identify() to be used in logs

				tableSummary := table.Resolve(ctx, client, nil, resourceSem, res)
				atomic.AddUint64(&summary.Resources, tableSummary.Resources)
				atomic.AddUint64(&summary.Errors, tableSummary.Errors)
				atomic.AddUint64(&summary.Panics, tableSummary.Panics)
			}()
		}
	}
	wg.Wait()
	logger.Info().Uint64("total_resources", summary.Resources).TimeDiff("duration", time.Now(), startTime).Msg("sync finished")
	return &summary, nil
}

func (p *SourcePlugin) listAndValidateTables(tables, skipTables []string) ([]string, error) {
	if len(tables) == 0 {
		return nil, fmt.Errorf("list of tables is empty")
	}

	// return an error if skip tables contains a wildcard or glob pattern
	for _, t := range skipTables {
		if strings.Contains(t, "*") {
			return nil, fmt.Errorf("glob matching in skipped table name %q is not supported", t)
		}
	}

	// handle wildcard entry
	if funk.Equal(tables, []string{"*"}) {
		allResources := make([]string, 0, len(p.tables))
		for _, k := range p.tables {
			if funk.ContainsString(skipTables, k.Name) {
				continue
			}
			allResources = append(allResources, k.Name)
		}
		return allResources, nil
	}

	// wildcard should not be combined with other tables
	if funk.ContainsString(tables, "*") {
		return nil, fmt.Errorf("wildcard \"*\" table not allowed with explicit tables")
	}

	// return an error if other kinds of glob-matching is detected
	for _, t := range tables {
		if strings.Contains(t, "*") {
			return nil, fmt.Errorf("glob matching in table name %q is not supported", t)
		}
	}

	// return an error if a table is both explicitly included and skipped
	for _, t := range tables {
		if funk.ContainsString(skipTables, t) {
			return nil, fmt.Errorf("table %s cannot be both included and skipped", t)
		}
	}

	// return an error if a given table name doesn't match any known tables
	for _, t := range tables {
		if !funk.ContainsString(p.tables.TableNames(), t) {
			return nil, fmt.Errorf("name %s does not match any known table names", t)
		}
	}

	// return an error if child table is included, but not its parent table
	selectedTables := map[string]bool{}
	for _, t := range tables {
		selectedTables[t] = true
	}
	for _, t := range tables {
		for _, tt := range p.tables {
			if tt.Name != t {
				continue
			}
			if tt.Parent != nil && !selectedTables[tt.Parent.Name] {
				return nil, fmt.Errorf("table %s is a child table, and requires its parent table %s to also be synced", t, tt.Parent.Name)
			}
		}
	}

	return tables, nil
}
