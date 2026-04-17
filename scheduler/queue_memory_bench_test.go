package scheduler

import (
	"context"
	"runtime"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
)

type benchClient struct{}

func (benchClient) ID() string { return "c" }

type benchItem struct {
	ID      string
	Payload [1024]byte // make Items sizeable so memory differences are visible
}

func buildDeepTables(rootFanout int) schema.Tables {
	deep := &schema.Table{
		Name: "deep",
		Resolver: func(ctx context.Context, _ schema.ClientMeta, p *schema.Resource, res chan<- any) error {
			items := make([]any, 10)
			for i := range items {
				items[i] = benchItem{ID: "d"}
			}
			res <- items
			return nil
		},
		Transform: transformers.TransformWithStruct(&benchItem{}),
	}
	mid := &schema.Table{
		Name: "mid",
		Resolver: func(ctx context.Context, _ schema.ClientMeta, p *schema.Resource, res chan<- any) error {
			items := make([]any, 10)
			for i := range items {
				items[i] = benchItem{ID: "m"}
			}
			res <- items
			return nil
		},
		Transform: transformers.TransformWithStruct(&benchItem{}),
		Relations: []*schema.Table{deep},
	}
	root := &schema.Table{
		Name: "root",
		Resolver: func(ctx context.Context, _ schema.ClientMeta, p *schema.Resource, res chan<- any) error {
			items := make([]any, rootFanout)
			for i := range items {
				items[i] = benchItem{ID: "r"}
			}
			res <- items
			return nil
		},
		Transform: transformers.TransformWithStruct(&benchItem{}),
		Relations: []*schema.Table{mid},
	}
	tables := schema.Tables{root}
	var apply func([]*schema.Table)
	apply = func(ts []*schema.Table) {
		for _, t := range ts {
			if t.Transform != nil {
				_ = t.Transform(t)
			}
			apply(t.Relations)
		}
	}
	apply(tables)
	return tables
}

func peakHeapMB() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapInuse / (1024 * 1024)
}

func BenchmarkQueue_InMemory(b *testing.B) {
	tables := buildDeepTables(50) // 50 × 10 × 10 = 5000 deep items
	for i := 0; i < b.N; i++ {
		s := NewScheduler(WithStrategy(StrategyShuffleQueue))
		_, err := s.SyncAll(context.Background(), benchClient{}, tables)
		if err != nil {
			b.Fatal(err)
		}
		b.ReportMetric(float64(peakHeapMB()), "peakheap_mb")
	}
}

func BenchmarkQueue_Badger(b *testing.B) {
	tables := buildDeepTables(50)
	for i := 0; i < b.N; i++ {
		dir := b.TempDir()
		store, err := NewStorageFromConfig(&QueueConfig{Type: QueueTypeBadger, Path: dir}, 1, "inv-1")
		if err != nil {
			b.Fatal(err)
		}
		s := NewScheduler(WithStrategy(StrategyShuffleQueue), WithStorage(store))
		_, err = s.SyncAll(context.Background(), benchClient{}, tables)
		if err != nil {
			b.Fatal(err)
		}
		_ = store.Close(context.Background())
		b.ReportMetric(float64(peakHeapMB()), "peakheap_mb")
	}
}
