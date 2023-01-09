package destination

import (
	"context"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

type worker struct {
	count int
	wg    sync.WaitGroup // implies usage by pointer
	ch    chan schema.CQTypes
	flush chan chan bool
	table *schema.Table
}

func newWorker(table *schema.Table) *worker {
	return &worker{
		count: 1,
		ch:    make(chan schema.CQTypes),
		flush: make(chan chan bool),
		table: table,
	}
}

func (p *Plugin) worker(ctx context.Context, w *worker) {
	w.wg.Add(1)
	defer w.wg.Done()

	resources := make([][]any, 0)
	sizeBytes := 0
	for {
		select {
		case r, ok := <-w.ch:
			if !ok {
				if len(resources) > 0 {
					p.flush(ctx, w, resources)
				}
				return
			}
			if len(resources) == p.spec.BatchSize || sizeBytes+r.Size() > p.spec.BatchSizeBytes {
				p.flush(ctx, w, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
			resources = append(resources, schema.TransformWithTransformer(p.client, r))
			sizeBytes += r.Size()
		case <-time.After(p.batchTimeout):
			if len(resources) > 0 {
				p.flush(ctx, w, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
		case done := <-w.flush:
			if len(resources) > 0 {
				p.flush(ctx, w, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func (p *Plugin) flush(ctx context.Context, w *worker, resources [][]any) {
	start := time.Now()
	err := p.client.WriteTableBatch(ctx, w.table, resources)
	dur := time.Since(start)
	batchSize := len(resources)

	if err != nil {
		p.errors.Add(uint64(batchSize))
		p.logger.Err(err).
			Str("table", w.table.Name).
			Int("len", batchSize).
			Dur("duration", dur).
			Msg("failed to write batch")
		return
	}

	p.writes.Add(uint64(batchSize))
	p.logger.Info().
		Str("table", w.table.Name).
		Int("len", batchSize).
		Dur("duration", dur).
		Msg("batch written successfully")
}

func (p *Plugin) writeManagedTableBatch(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, res <-chan schema.DestinationResource) error {
	syncTime = syncTime.UTC()
	SetDestinationManagedCqColumns(tables)

	workers := make(map[string]*worker, len(tables))

	p.workersLock.Lock()
	for _, table := range tables.FlattenTables() {
		if w, ok := p.workers[table.Name]; ok {
			w.count++
			continue
		}
		w := newWorker(table)
		p.workers[table.Name] = w
		// we save this locally because we don't want to access the map after that so we can
		// keep the workersLock for as short as possible
		workers[table.Name] = w
		go p.worker(ctx, w)
	}
	p.workersLock.Unlock()

	sourceColumn := &schema.Text{}
	_ = sourceColumn.Set(sourceSpec.Name)
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
	flushChannels := make([]chan bool, 0, len(workers))
	for _, w := range workers {
		flushCh := make(chan bool)
		flushChannels = append(flushChannels, flushCh)
		w.flush <- flushCh
	}
	for _, ch := range flushChannels {
		<-ch
	}

	p.workersLock.Lock()
	for tableName, w := range p.workers {
		if w.count--; w.count == 0 {
			close(w.ch)
			w.wg.Wait()
			delete(p.workers, tableName)
		}
	}
	p.workersLock.Unlock()
	return nil
}
