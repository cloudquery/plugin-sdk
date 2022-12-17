package plugins

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

type NewDestinationClientFunc func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error)

type DestinationClient interface {
	schema.CQTypeTransformer
	ReverseTransformValues(table *schema.Table, values []interface{}) (schema.CQTypes, error)
	Migrate(ctx context.Context, tables schema.Tables) error
	Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- []interface{}) error
	PreWrite(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) (interface{}, error)
	WriteTableBatch(ctx context.Context, writeClient interface{}, table *schema.Table, data [][]interface{}) error
	PostWrite(ctx context.Context, writeClient interface{}, tables schema.Tables, sourceName string, syncTime time.Time) error
	DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error
	Close(ctx context.Context) error
}

type destinationWorker struct {
	count int
	wg *sync.WaitGroup
	ch chan schema.CQTypes
}

type DestinationPlugin struct {
	// Name of destination plugin i.e postgresql,snowflake
	name string
	// Version of the destination plugin
	version string
	// Called upon configure call to validate and init configuration
	newDestinationClient NewDestinationClientFunc
	// initialized destination client
	client DestinationClient
	// spec the client was initialized with
	spec specs.Destination
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger
	// destination plugin metrics
	metrics map[string]*DestinationMetrics
	metricsLock *sync.RWMutex
	
	workers map[string]*destinationWorker
	workersLock *sync.Mutex
}

// this is the buffer per worker so we can continue writing
// while the previous batch is being sent over the network
// const bufferSize = 1000
// batchTimeout is the timeout for a batch to be sent to the destination if no resources are received
const batchTimeout = 20 * time.Second

func NewDestinationPlugin(name string, version string, newDestinationClient NewDestinationClientFunc) *DestinationPlugin {
	p := &DestinationPlugin{
		name:                 name,
		version:              version,
		newDestinationClient: newDestinationClient,
		metrics:              make(map[string]*DestinationMetrics),
		metricsLock:          &sync.RWMutex{},
		workers:              make(map[string]*destinationWorker),
		workersLock:          &sync.Mutex{},
	}
	return p
}

func (p *DestinationPlugin) Name() string {
	return p.name
}

func (p *DestinationPlugin) Version() string {
	return p.version
}

func (p *DestinationPlugin) Metrics() DestinationMetrics {
	metrics := DestinationMetrics{}
	p.metricsLock.RLock()
	for _, m := range p.metrics {
		metrics.Errors += m.Errors
		metrics.Writes += m.Writes
	}
	p.metricsLock.RUnlock()
	
	return metrics
}

// we need lazy loading because we want to be able to initialize after
func (p *DestinationPlugin) Init(ctx context.Context, logger zerolog.Logger, spec specs.Destination) error {
	var err error
	p.logger = logger
	p.spec = spec
	p.spec.SetDefaults()
	p.client, err = p.newDestinationClient(ctx, logger, spec)
	if err != nil {
		return err
	}
	if p.client == nil {
		return fmt.Errorf("destination client is nil")
	}
	return nil
}

// we implement all DestinationClient functions so we can hook into pre-post behavior
func (p *DestinationPlugin) Migrate(ctx context.Context, tables schema.Tables) error {
	SetDestinationManagedCqColumns(tables)
	return p.client.Migrate(ctx, tables)
}

func (p *DestinationPlugin) readAll(ctx context.Context, table *schema.Table, sourceName string) ([]schema.CQTypes, error) {
	var readErr error
	ch := make(chan schema.CQTypes)
	go func() {
		defer close(ch)
		readErr = p.Read(ctx, table, sourceName, ch)
	}()
	//nolint:prealloc
	var resources []schema.CQTypes
	for resource := range ch {
		resources = append(resources, resource)
	}
	return resources, readErr
}

func (p *DestinationPlugin) Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- schema.CQTypes) error {
	SetDestinationManagedCqColumns(schema.Tables{table})
	ch := make(chan []interface{})
	var err error
	go func() {
		defer close(ch)
		err = p.client.Read(ctx, table, sourceName, ch)
	}()
	for resource := range ch {
		r, err := p.client.ReverseTransformValues(table, resource)
		if err != nil {
			return err
		}
		res <- r
	}
	return err
}

// this function is currently used mostly for testing so it's not a public api
func (p *DestinationPlugin) writeOne(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time, resource schema.DestinationResource){
	resources := []schema.DestinationResource{resource}
	p.writeAll(ctx, tables, sourceName, syncTime, resources)
}

// this function is currently used mostly for testing so it's not a public api
func (p *DestinationPlugin) writeAll(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time, resources []schema.DestinationResource){
	ch := make(chan schema.DestinationResource, len(resources))
	for _, resource := range resources {
		ch <- resource
	}
	close(ch)
	p.Write(ctx, tables, sourceName, syncTime, ch)
}

func (p *DestinationPlugin) worker(ctx context.Context, writeClient interface{}, metrics *DestinationMetrics, table *schema.Table, ch <-chan schema.CQTypes) {
	resources := make([][]interface{}, 0)
	totalSize := 0
	for {
		select {
		case r, ok := <-ch:
			if ok {
				totalSize += r.Size()
				resources = append(resources, schema.TransformWithTransformer(p.client, r))
				if len(resources) == p.spec.BatchSize {
					start := time.Now()
					if err := p.client.WriteTableBatch(ctx, writeClient, table,  resources); err != nil {
						p.logger.Err(err).Str("table", table.Name).Int("len", p.spec.BatchSize).Int("size", totalSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
						// we don't return as we need to continue until channel is closed otherwise there will be a deadlock
						atomic.AddUint64(&metrics.Errors, uint64(p.spec.BatchSize))
					} else {
						p.logger.Info().Str("table", table.Name).Int("len", p.spec.BatchSize).Int("size", totalSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
						atomic.AddUint64(&metrics.Writes, uint64(p.spec.BatchSize))
					}
					resources = make([][]interface{}, 0)
					totalSize = 0
				}
			} else {
				if len(resources) > 0 {
					start := time.Now()
					if err := p.client.WriteTableBatch(ctx, writeClient, table, resources); err != nil {
						p.logger.Err(err).Str("table", table.Name).Int("len", len(resources)).Int("size", totalSize).Dur("duration", time.Since(start)).Msg("failed to write last batch")
						atomic.AddUint64(&metrics.Errors, uint64(len(resources)))
					} else {
						p.logger.Info().Str("table", table.Name).Int("len", len(resources)).Int("size", totalSize).Dur("duration", time.Since(start)).Msg("last batch written successfully")
						atomic.AddUint64(&metrics.Writes, uint64(len(resources)))
					}
				}
				return
			}
		case <-time.After(batchTimeout):
			if len(resources) > 0 {
				start := time.Now()
				if err := p.client.WriteTableBatch(ctx, writeClient, table,  resources); err != nil {
					p.logger.Err(err).Str("table", table.Name).Int("len", len(resources)).Int("size", totalSize).Dur("time", time.Since(start)).Msg("failed to write batch on timeout")
					// we don't return as we need to continue until channel is closed otherwise there will be a deadlock
					atomic.AddUint64(&metrics.Errors, uint64(len(resources)))
				} else {
					p.logger.Info().Str("table", table.Name).Int("len", len(resources)).Int("size", totalSize).Dur("time", time.Since(start)).Msg("batch written successfully on timeout")
					atomic.AddUint64(&metrics.Writes, uint64(len(resources)))
				}
				resources = make([][]interface{}, 0)
				totalSize = 0
			}
		}
	}
}

func (p *DestinationPlugin) Write(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time, res <-chan schema.DestinationResource) error {
	syncTime = syncTime.UTC()
	SetDestinationManagedCqColumns(tables)
	workers := make(map[string]*destinationWorker, len(tables))
	metrics := &DestinationMetrics{}
	writeClient, err := p.client.PreWrite(ctx, tables, sourceName, syncTime);
	if err != nil {
		return fmt.Errorf("pre write failed: %w", err)
	}

	p.workersLock.Lock()
	for _, table := range tables.FlattenTables() {
		table := table
		if p.workers[table.Name] == nil {
			ch := make(chan schema.CQTypes)
			wg := &sync.WaitGroup{}
			p.workers[table.Name] = &destinationWorker{
				count: 1,
				ch: ch,
				wg: wg,
			}
			for i := 0; i < p.spec.Workers; i++ {
				wg.Add(1)
				go func () {
					defer wg.Done()
					p.worker(ctx, writeClient, metrics, table, ch)
				}()
			}
		} else {
			p.workers[table.Name].count++
		}
		// we save this locally because we don't want to access the map after that so we can
		// keep the workersLock for as short as possible
		workers[table.Name] = p.workers[table.Name]
	}
	p.workersLock.Unlock()

	sourceColumn := &schema.Text{}
	_ = sourceColumn.Set(sourceName)
	syncTimeColumn := &schema.Timestamptz{}
	_ = syncTimeColumn.Set(syncTime)
	for r := range res {
		r.Data = append([]schema.CQType{sourceColumn, syncTimeColumn}, r.Data...)
		workers[r.TableName].ch <- r.Data
	}
	p.workersLock.Lock()
	for tableName := range workers {
		p.workers[tableName].count--
		if p.workers[tableName].count == 0 {
			close(p.workers[tableName].ch)
			p.workers[tableName].wg.Wait()
			delete(p.workers, tableName)
		}
	}
	p.workersLock.Unlock()

	if err := p.client.PostWrite(ctx, writeClient, tables, sourceName, syncTime); err != nil {
		return fmt.Errorf("post write failed: %w", err)
	}

	if p.spec.WriteMode == specs.WriteModeOverwriteDeleteStale {
		if err := p.DeleteStale(ctx, tables, sourceName, syncTime); err != nil {
			p.logger.Err(err).Msg("failed to delete stale resources")
			// TODO: Update metrics
		}
	}
	p.logger.Info().Uint64("writes", metrics.Writes).Uint64("errors", metrics.Errors).Msg("finished writing")
	p.metricsLock.Lock()
	p.metrics[sourceName] = metrics
	p.metricsLock.Unlock()
	return nil
}

func (p *DestinationPlugin) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	syncTime = syncTime.UTC()
	return p.client.DeleteStale(ctx, tables, sourceName, syncTime)
}

func (p *DestinationPlugin) Close(ctx context.Context) error {
	return p.client.Close(ctx)
}

// Overwrites or adds the CQ columns that are managed by the destination plugins (_cq_sync_time, _cq_source_name).
func SetDestinationManagedCqColumns(tables []*schema.Table) {
	for _, table := range tables {
		table.OverwriteOrAddColumn(&schema.CqSyncTimeColumn)
		table.OverwriteOrAddColumn(&schema.CqSourceNameColumn)

		SetDestinationManagedCqColumns(table.Relations)
	}
}
