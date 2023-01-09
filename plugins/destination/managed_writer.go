package destination

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

type worker struct {
	count int
	wg    *sync.WaitGroup
	ch    chan schema.CQTypes
	flush chan chan bool
}

func (p *Plugin) worker(ctx context.Context, metrics *Metrics, table *schema.Table, ch <-chan schema.CQTypes, flush <-chan chan bool) {
	resources := make([][]any, 0)
	sizeBytes := 0
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if len(resources) > 0 {
					p.flush(ctx, metrics, table, resources)
				}
				return
			}
			if len(resources) == p.spec.BatchSize || sizeBytes+r.Size() > p.spec.BatchSizeBytes {
				p.flush(ctx, metrics, table, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
			resources = append(resources, schema.TransformWithTransformer(p.client, r))
			sizeBytes += r.Size()
		case <-time.After(p.batchTimeout):
			if len(resources) > 0 {
				p.flush(ctx, metrics, table, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
		case done := <-flush:
			if len(resources) > 0 {
				p.flush(ctx, metrics, table, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func (p *Plugin) flush(ctx context.Context, metrics *Metrics, table *schema.Table, resources [][]any) {
	start := time.Now()
	if err := p.client.WriteTableBatch(ctx, table, resources); err != nil {
		p.logger.Err(err).Str("table", table.Name).Int("len", p.spec.BatchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
		// we don't return an error as we need to continue until channel is closed otherwise there will be a deadlock
		atomic.AddUint64(&metrics.Errors, uint64(p.spec.BatchSize))
	} else {
		p.logger.Info().Str("table", table.Name).Int("len", p.spec.BatchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
		atomic.AddUint64(&metrics.Writes, uint64(p.spec.BatchSize))
	}
}

func (p *Plugin) writeManagedTableBatch(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time, res <-chan schema.DestinationResource) error {
	syncTime = syncTime.UTC()
	SetDestinationManagedCqColumns(tables)

	workers := make(map[string]*worker, len(tables))
	metrics := &Metrics{}

	p.workersLock.Lock()
	for _, table := range tables.FlattenTables() {
		table := table
		if p.workers[table.Name] == nil {
			ch := make(chan schema.CQTypes)
			flush := make(chan chan bool)
			wg := &sync.WaitGroup{}
			p.workers[table.Name] = &worker{
				count: 1,
				ch:    ch,
				flush: flush,
				wg:    wg,
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.worker(ctx, metrics, table, ch, flush)
			}()
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
		// this is a check to keep backward compatible for sources that are not adding
		// source and sync time
		if len(r.Data) < len(tables.Get(r.TableName).Columns) {
			r.Data = append([]schema.CQType{sourceColumn, syncTimeColumn}, r.Data...)
		}
		workers[r.TableName].ch <- r.Data
	}

	// flush and wait for all workers to finish flush before finish and calling delete stale
	// This is because destinations can be longed lived and called from multiple sources
	flushChannels := make(map[string]chan bool, len(workers))
	for tableName, w := range workers {
		flushCh := make(chan bool)
		flushChannels[tableName] = flushCh
		w.flush <- flushCh
	}
	for tableName := range flushChannels {
		<-flushChannels[tableName]
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

	if p.spec.WriteMode == specs.WriteModeOverwriteDeleteStale {
		if err := p.DeleteStale(ctx, tables, sourceName, syncTime); err != nil {
			return err
		}
	}
	return nil
}
