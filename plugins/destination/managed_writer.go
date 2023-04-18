package destination

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/util"
	"github.com/cloudquery/plugin-sdk/v2/internal/pk"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/specs"
)

type worker struct {
	count int
	wg    *sync.WaitGroup
	ch    chan arrow.Record
	flush chan chan bool
}

func (p *Plugin) worker(ctx context.Context, metrics *Metrics, table *arrow.Schema, ch <-chan arrow.Record, flush <-chan chan bool) {
	sizeBytes := int64(0)
	resources := make([]arrow.Record, 0)
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				if len(resources) > 0 {
					p.flush(ctx, metrics, table, resources)
				}
				return
			}
			if len(resources) == p.spec.BatchSize || sizeBytes+util.TotalRecordSize(r) > int64(p.spec.BatchSizeBytes) {
				p.flush(ctx, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
			resources = append(resources, r)
			sizeBytes += util.TotalRecordSize(r)
		case <-time.After(p.batchTimeout):
			if len(resources) > 0 {
				p.flush(ctx, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
		case done := <-flush:
			if len(resources) > 0 {
				p.flush(ctx, metrics, table, resources)
				resources = make([]arrow.Record, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func (p *Plugin) flush(ctx context.Context, metrics *Metrics, table *arrow.Schema, resources []arrow.Record) {
	resources = p.removeDuplicatesByPK(table, resources)
	tableName := schema.TableName(table)
	start := time.Now()
	batchSize := len(resources)
	if err := p.client.WriteTableBatch(ctx, table, resources); err != nil {
		p.logger.Err(err).Str("table", tableName).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
		// we don't return an error as we need to continue until channel is closed otherwise there will be a deadlock
		atomic.AddUint64(&metrics.Errors, uint64(batchSize))
	} else {
		p.logger.Info().Str("table", tableName).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
		atomic.AddUint64(&metrics.Writes, uint64(batchSize))
	}
	for _, r := range resources {
		r.Release()
	}
}

func (*Plugin) removeDuplicatesByPK(table *arrow.Schema, resources []arrow.Record) []arrow.Record {
	pkIndices := schema.PrimaryKeyIndices(table)
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
		// duplicate, release early
		r.Release()
	}

	return res
}

func (p *Plugin) writeManagedTableBatch(ctx context.Context, _ specs.Source, tables schema.Schemas, _ time.Time, res <-chan arrow.Record) error {
	workers := make(map[string]*worker, len(tables))
	metrics := &Metrics{}

	p.workersLock.Lock()
	for _, table := range tables {
		table := table
		tableName := schema.TableName(table)
		if p.workers[schema.TableName(table)] == nil {
			ch := make(chan arrow.Record)
			flush := make(chan chan bool)
			wg := &sync.WaitGroup{}
			p.workers[schema.TableName(table)] = &worker{
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
			p.workers[tableName].count++
		}
		// we save this locally because we don't want to access the map after that so we can
		// keep the workersLock for as short as possible
		workers[tableName] = p.workers[tableName]
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
		r.Retain()
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
			delete(p.workers, tableName)
		}
	}
	p.workersLock.Unlock()
	return nil
}
