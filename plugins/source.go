package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
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

const defaultConcurrency = 20

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
	c, err := p.newExecutionClient(ctx, p.logger, spec)
	if err != nil {
		return fmt.Errorf("failed to create execution client for source plugin %s: %w", p.name, err)
	}

	// limiter used to limit the amount of resources fetched concurrently
	concurrency := spec.Concurrency
	if concurrency == 0 {
		concurrency = defaultConcurrency
	}
	p.logger.Info().Uint64("concurrency", concurrency).Msg("starting fetch")

	startTime := time.Now()
	tableNames, err := p.interpolateAllResources(spec.Tables)
	if err != nil {
		return err
	}

	// this is the same fetchtime for all resources
	syncTime := time.Now()

	tableJobs := make([]workerJob, 0)
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
			tableJobs = append(tableJobs, workerJob{
				ctx:      ctx,
				logger:   &p.logger,
				client:   client,
				table:    table,
				syncTime: syncTime,
				res:      res,
			})
		}
	}

	jobsCount := len(tableJobs)
	jobs := make(chan workerJob, jobsCount)
	results := make(chan workerResult, jobsCount)

	for w := uint64(1); w <= concurrency; w++ {
		go worker(jobs, results)
	}

	for i := 0; i < jobsCount; i++ {
		jobs <- tableJobs[i]
	}
	close(jobs)

	totalResources := 0
	for i := 0; i < jobsCount; i++ {
		result := <-results
		totalResources += result.totalResources
	}
	close(results)

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

type workerJob struct {
	ctx      context.Context
	logger   *zerolog.Logger
	client   schema.ClientMeta
	table    *schema.Table
	syncTime time.Time
	parent   *schema.Resource
	res      chan<- *schema.Resource
}

type workerResult struct {
	totalResources int
}

func worker(jobs chan workerJob, results chan<- workerResult) {
	for job := range jobs {
		tableStartTime := time.Now()
		job.logger.Info().Str("table", job.table.Name).Msg("fetch table started")
		resources := job.table.Resolve(job.ctx, job.client, job.syncTime, job.parent, job.res)
		job.logger.Info().Str("table", job.table.Name).Int("total_resources", resources).TimeDiff("duration", time.Now(), tableStartTime).Msg("fetch table finished")
		results <- workerResult{totalResources: resources}
	}
}
