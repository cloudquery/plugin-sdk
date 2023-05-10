package destination

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/util"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/internal/pk"
	"github.com/cloudquery/plugin-sdk/v3/schema"
)

type worker struct {
	count int
	wg    *sync.WaitGroup
	ch    chan arrow.Record
	flush chan chan bool

	openError bool
}

func (*Plugin) dummyWorker(_ context.Context, ch <-chan arrow.Record, flush <-chan chan bool) {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case done := <-flush:
			done <- true
		}
	}
}

func (p *Plugin) worker(ctx context.Context, sourceSpec specs.Source, metrics *Metrics, table *schema.Table, ch <-chan arrow.Record, flush <-chan chan bool) {
	sizeBytes := int64(0)
	resources := make([]arrow.Record, 0)
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if len(resources) > 0 {
					p.flush(ctx, sourceSpec, metrics, table, resources)
				}
				return
			}
			if len(resources) == p.spec.BatchSize || sizeBytes+util.TotalRecordSize(r) > int64(p.spec.BatchSizeBytes) {
				p.flush(ctx, sourceSpec, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
			resources = append(resources, r)
			sizeBytes += util.TotalRecordSize(r)
		case <-time.After(p.batchTimeout):
			if len(resources) > 0 {
				p.flush(ctx, sourceSpec, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
		case done := <-flush:
			if len(resources) > 0 {
				p.flush(ctx, sourceSpec, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func (p *Plugin) flush(ctx context.Context, sourceSpec specs.Source, metrics *Metrics, table *schema.Table, resources []arrow.Record) {
	resources = p.removeDuplicatesByPK(table, resources)
	start := time.Now()
	batchSize := len(resources)
	if err := p.client.WriteTableBatch(ctx, sourceSpec, table, resources); err != nil {
		p.logger.Err(err).Str("table", table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
		// we don't return an error as we need to continue until channel is closed otherwise there will be a deadlock
		atomic.AddUint64(&metrics.Errors, uint64(batchSize))
	} else {
		p.logger.Info().Str("table", table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
		atomic.AddUint64(&metrics.Writes, uint64(batchSize))
	}
}

func (*Plugin) removeDuplicatesByPK(table *schema.Table, resources []arrow.Record) []arrow.Record {
	pkIndices := table.PrimaryKeysIndexes()
	// special case where there's no PK at all
	if len(pkIndices) == 0 {
		return resources
	}

	pks := make(map[string]struct{}, len(resources))
	res := make([]arrow.Record, 0, len(resources))
	for _, r := range resources {
		if r.NumRows() > 1 {
			panic(fmt.Sprintf("record with more than 1 row: %d", r.NumRows()))
		}
		key := pk.String(r)
		_, ok := pks[key]
		if !ok {
			pks[key] = struct{}{}
			res = append(res, r)
			continue
		}
		// duplicate, release
		r.Release()
	}

	return res
}

func (p *Plugin) writeManagedTableBatch(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, _ time.Time, res <-chan arrow.Record) error {
	workers := make(map[string]*worker, len(tables))
	metrics := &Metrics{}

	p.workersLock.Lock()
	for _, table := range tables {
		table := table
		if p.workers[table.Name] == nil {
			ch := make(chan arrow.Record)
			flush := make(chan chan bool)
			wg := &sync.WaitGroup{}
			var errored bool
			if p.clientOCW != nil {
				if err := p.clientOCW.OpenTable(ctx, sourceSpec, table); err != nil {
					p.logger.Err(err).Str("table", table.Name).Msg("OpenTable failed")
					// we don't return an error as we need to continue until channel is closed otherwise there will be a deadlock
					atomic.AddUint64(&metrics.Errors, 1)
					errored = true
				}
			}
			p.workers[table.Name] = &worker{
				count: 1,
				ch:    ch,
				flush: flush,
				wg:    wg,

				openError: errored,
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				if errored {
					p.dummyWorker(ctx, ch, flush)
					return
				}
				p.worker(ctx, sourceSpec, metrics, table, ch, flush)
			}()
		} else {
			p.workers[table.Name].count++
		}
		// we save this locally because we don't want to access the map after that so we can
		// keep the workersLock for as short as possible
		workers[table.Name] = p.workers[table.Name]
	}
	p.workersLock.Unlock()

	for r := range res {
		tableName, ok := r.Schema().Metadata().GetValue(schema.MetadataTableName)
		if !ok {
			return fmt.Errorf("missing table name in record metadata")
		}
		if _, ok := workers[tableName]; !ok {
			return fmt.Errorf("table %s not found in destination", tableName)
		}
		workers[tableName].ch <- r
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

			if p.clientOCW != nil && !p.workers[tableName].openError {
				if err := p.clientOCW.CloseTable(ctx, sourceSpec, tables.Get(tableName)); err != nil {
					p.logger.Err(err).Str("table", tableName).Msg("CloseTable failed")
					atomic.AddUint64(&metrics.Errors, 1)
				}
			}

			delete(p.workers, tableName)
		}
	}
	p.workersLock.Unlock()
	return nil
}
