package destination

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/internal/pk"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/getsentry/sentry-go"
)

type worker struct {
	count int
	wg    sync.WaitGroup // implies usage by pointer only

	ch    chan schema.CQTypes
	flush chan chan bool

	table   *schema.Table
	srcSpec *specs.Source
}

func (p *Plugin) worker(ctx context.Context, metrics *Metrics, w *worker) {
	defer w.wg.Done()
	resources := make([][]any, 0)
	sizeBytes := 0
	for {
		select {
		case r, ok := <-w.ch:
			if !ok {
				if len(resources) > 0 {
					p.flush(ctx, metrics, w, resources)
				}
				return
			}
			if len(resources) == p.spec.BatchSize || sizeBytes+r.Size() > p.spec.BatchSizeBytes {
				p.flush(ctx, metrics, w, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
			resources = append(resources, schema.TransformWithTransformer(p.client, r))
			sizeBytes += r.Size()
		case <-time.After(p.batchTimeout):
			if len(resources) > 0 {
				p.flush(ctx, metrics, w, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
		case done := <-w.flush:
			if len(resources) > 0 {
				p.flush(ctx, metrics, w, resources)
				resources = make([][]any, 0)
				sizeBytes = 0
			}
			done <- true
		}
	}
}

func (p *Plugin) flush(ctx context.Context, metrics *Metrics, w *worker, resources [][]any) {
	resources = p.removeDuplicatesByPK(w, resources)

	start := time.Now()
	batchSize := len(resources)
	if err := p.client.WriteTableBatch(ctx, w.table, resources); err != nil {
		p.logger.Err(err).Str("table", w.table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("failed to write batch")
		// we don't return an error as we need to continue until channel is closed otherwise there will be a deadlock
		atomic.AddUint64(&metrics.Errors, uint64(batchSize))
	} else {
		p.logger.Info().Str("table", w.table.Name).Int("len", batchSize).Dur("duration", time.Since(start)).Msg("batch written successfully")
		atomic.AddUint64(&metrics.Writes, uint64(batchSize))
	}
}

func (p *Plugin) removeDuplicatesByPK(w *worker, resources [][]any) [][]any {
	pks := make(map[string]struct{}, len(resources))
	res := make([][]any, 0, len(resources))
	var reported bool
	for _, r := range resources {
		key := pk.String(w.table, r)
		_, ok := pks[key]
		switch {
		case !ok:
			pks[key] = struct{}{}
			res = append(res, r)
			continue
		case reported:
			continue
		}

		reported = true
		pkSpec := "(" + strings.Join(w.table.PrimaryKeys(), ",") + ")"

		// log err
		p.logger.Error().
			Str("source_plugin", w.srcSpec.Name).
			Str("source_version", w.srcSpec.Version).
			Str("table", w.table.Name).
			Str("pk", pkSpec).
			Str("value", key).
			Msg("duplicate primary key")

		// send to Sentry only once per table,
		// to avoid sending too many duplicate messages
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("source_plugin", w.srcSpec.Name)
			scope.SetTag("source_version", w.srcSpec.Version)
			scope.SetTag("table", w.table.Name)
			scope.SetExtra("pk", pkSpec)
			sentry.CurrentHub().CaptureMessage("duplicate primary key")
		})
	}

	return res
}

func (p *Plugin) writeManagedTableBatch(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, res <-chan schema.DestinationResource) error {
	syncTime = syncTime.UTC()
	SetDestinationManagedCqColumns(tables)

	workers := make(map[string]*worker, len(tables))
	metrics := &Metrics{}

	p.workersLock.Lock()
	for _, table := range tables.FlattenTables() {
		table := table
		if p.workers[table.Name] == nil {
			w := &worker{
				count:   1,
				ch:      make(chan schema.CQTypes),
				flush:   make(chan chan bool),
				table:   table,
				srcSpec: &sourceSpec,
			}
			p.workers[table.Name] = w
			w.wg.Add(1)
			go p.worker(ctx, metrics, w)
		} else {
			p.workers[table.Name].count++
		}
		// we save this locally because we don't want to access the map after that so we can
		// keep the workersLock for as short as possible
		workers[table.Name] = p.workers[table.Name]
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
