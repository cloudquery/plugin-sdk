package plugins

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type BenchmarkScenario struct {
	Clients           int
	Tables            int
	ChildrenPerTable  int
	Columns           int
	ColumnResolvers   int // number of columns with custom resolvers
	ResourcesPerTable int
	ResourcesPerPage  int
	ResolverMin       time.Duration
	ResolverStdDev    time.Duration
	ResolverMean      time.Duration
	Concurrency       uint64
}

func (s *BenchmarkScenario) SetDefaults() {
	if s.Clients == 0 {
		s.Clients = 1
	}
	if s.Tables == 0 {
		s.Tables = 1
	}
	if s.Columns == 0 {
		s.Columns = 10
	}
	if s.ResourcesPerTable == 0 {
		s.ResourcesPerTable = 100
	}
	if s.ResourcesPerPage == 0 {
		s.ResourcesPerPage = 10
	}
	if s.ResolverMin == 0 {
		s.ResolverMin = time.Millisecond
	}
	if s.ResolverStdDev == 0 {
		s.ResolverStdDev = 100 * time.Millisecond
	}
	if s.ResolverMean == 0 {
		s.ResolverMean = 10 * time.Millisecond
	}
}

type SourceBenchmark struct {
	*BenchmarkScenario

	b      *testing.B
	tables []*schema.Table
	plugin *SourcePlugin

	apiCalls atomic.Int64
}

func NewSourceBenchmark(b *testing.B, scenario BenchmarkScenario) *SourceBenchmark {
	scenario.SetDefaults()
	sb := &SourceBenchmark{
		BenchmarkScenario: &scenario,
		b:                 b,
		tables:            nil,
		plugin:            nil,
	}
	sb.setup(b)
	return sb
}

func (s *SourceBenchmark) setup(b *testing.B) {
	tableResolver := func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
		s.simulateAPICall(s.ResolverMin, s.ResolverStdDev, s.ResolverMean)
		total := 0
		for total < s.ResourcesPerTable {
			num := min(s.ResourcesPerPage, s.ResourcesPerTable-total)
			resources := make([]struct {
				Column1 string
			}, num)
			for i := 0; i < num; i++ {
				resources[i] = struct {
					Column1 string
				}{
					Column1: "test-table",
				}
			}
			res <- resources
			total += num
		}
		return nil
	}
	preResourceResolver := func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
		s.simulateAPICall(s.ResolverMin, s.ResolverStdDev, s.ResolverMean)
		resource.Item = struct {
			Column1 string
		}{
			Column1: "test-pre",
		}
		return nil
	}
	columnResolver := func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
		s.simulateAPICall(s.ResolverMin, s.ResolverStdDev, s.ResolverMean)
		return resource.Set(c.Name, "test")
	}
	s.tables = make([]*schema.Table, s.Tables)
	for i := 0; i < s.Tables; i++ {
		columns := make([]schema.Column, s.Columns)
		for u := 0; u < s.Columns; u++ {
			columns[u] = schema.Column{
				Name: fmt.Sprintf("column%d", u),
				Type: schema.TypeString,
			}
			if u < s.ColumnResolvers {
				columns[u].Resolver = columnResolver
			}
		}
		relations := make([]*schema.Table, s.ChildrenPerTable)
		for u := 0; u < s.ChildrenPerTable; u++ {
			relations[u] = &schema.Table{
				Name:                fmt.Sprintf("table%d_child%d", i, u),
				Columns:             columns,
				Resolver:            tableResolver,
				PreResourceResolver: preResourceResolver,
			}
		}
		s.tables[i] = &schema.Table{
			Name:                fmt.Sprintf("table%d", i),
			Columns:             columns,
			Relations:           relations,
			Resolver:            tableResolver,
			Multiplex:           nMultiplexer(s.Clients),
			PreResourceResolver: preResourceResolver,
		}
		for u := range relations {
			relations[u].Parent = s.tables[i]
		}
	}

	plugin := NewSourcePlugin(
		"testSourcePlugin",
		"1.0.0",
		s.tables,
		newTestExecutionClient,
	)
	plugin.SetLogger(zerolog.New(zerolog.NewTestWriter(b)).Level(zerolog.WarnLevel))
	s.plugin = plugin
	s.b = b
}

func (s *SourceBenchmark) simulateAPICall(min, stdDev, mean time.Duration) {
	s.apiCalls.Add(1)
	sample := int(rand.NormFloat64()*float64(stdDev) + float64(mean))
	duration := time.Duration(sample)
	if duration < min {
		time.Sleep(min)
		return
	}
	time.Sleep(duration)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *SourceBenchmark) Run() {
	ctx := context.Background()
	spec := specs.Source{
		Name:         "testSource",
		Path:         "cloudquery/testSource",
		Tables:       []string{"*"},
		Version:      "v1.0.0",
		Destinations: []string{"test"},
		Concurrency:  s.Concurrency,
	}
	resources := make(chan *schema.Resource)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(resources)
		return s.plugin.Sync(ctx,
			spec,
			resources)
	})
	s.b.ResetTimer()
	start := time.Now()
	totalResources := 0
	for range resources {
		// read resources channel until empty
		totalResources++
	}
	if err := g.Wait(); err != nil {
		s.b.Fatal(err)
	}
	end := time.Now()
	s.b.StopTimer()
	s.b.ReportMetric(0, "ns/op") // drop default ns/op output
	s.b.ReportMetric(float64(totalResources)/(end.Sub(start).Seconds()), "resources/s")
	s.b.ReportMetric(float64(totalResources)/s.lowerBound().Seconds(), "targetResources/s")

	// Enable the below metrics for more verbose information about the scenario:
	//s.b.ReportMetric(float64(totalResources), "resources")
	//s.b.ReportMetric(float64(s.apiCalls.Load()), "apiCalls")
}

// lowerBound calculates a rough lower bound on the sync time so that we know how
// much room there is for optimization. This does not currently take the "concurrency"
// value into account. Use this number only as a rough guide.
func (s *SourceBenchmark) lowerBound() time.Duration {
	// we require one API call per page
	pages := s.ResourcesPerTable / s.ResourcesPerPage
	if s.ResourcesPerTable%s.ResourcesPerPage == 0 {
		pages++
	}

	// Use the mean time + stdDev for now, but that's not 100% accurate for many reasons. One
	// is that samples are rounded up to the minimum time. Use only as a rough guide.
	longestLoad := s.ResolverMean + s.ResolverStdDev
	minTime := longestLoad * time.Duration(pages)

	// double this because PreResourceResolver requires an additional call
	factor := 2
	if s.ColumnResolvers > 0 {
		// column resolvers also take time, and have to be called after PreResourceResolver
		factor++
	}
	minTime *= time.Duration(factor)

	// every child table will need to make one API call per resource in the
	// parent table. Theoretically these calls can be done in parallel, and child tables
	// can all be resolved in parallel too.
	if s.ChildrenPerTable > 0 {
		minTime += longestLoad * time.Duration(pages) * time.Duration(factor)
	}

	return minTime
}

type benchmarkClient struct {
	num int
}

func (b benchmarkClient) ID() string {
	return fmt.Sprintf("client%d", b.num)
}

func nMultiplexer(n int) schema.Multiplexer {
	return func(meta schema.ClientMeta) []schema.ClientMeta {
		clients := make([]schema.ClientMeta, n)
		for i := 0; i < n; i++ {
			clients[i] = benchmarkClient{
				num: i,
			}
		}
		return clients
	}
}

func BenchmarkDefaultConcurrency(b *testing.B) {
	benchmarkWithConcurrency(b, 0)
}

func benchmarkWithConcurrency(b *testing.B, concurrency uint64) {
	b.ReportAllocs()
	bs := BenchmarkScenario{
		Clients:           25,
		Tables:            5,
		Columns:           10,
		ColumnResolvers:   1,
		ResourcesPerTable: 100,
		ResourcesPerPage:  50,
		ResolverMin:       1 * time.Millisecond,
		ResolverMean:      10 * time.Millisecond,
		ResolverStdDev:    100 * time.Millisecond,
		Concurrency:       concurrency,
	}
	sb := NewSourceBenchmark(b, bs)
	sb.Run()
}

func BenchmarkTablesWithChildrenDefaultConcurrency(b *testing.B) {
	benchmarkTablesWithChildrenConcurrency(b, 0)
}

func benchmarkTablesWithChildrenConcurrency(b *testing.B, concurrency uint64) {
	b.ReportAllocs()
	bs := BenchmarkScenario{
		Clients:           2,
		Tables:            2,
		ChildrenPerTable:  2,
		Columns:           10,
		ColumnResolvers:   1,
		ResourcesPerTable: 100,
		ResourcesPerPage:  50,
		ResolverMin:       1 * time.Millisecond,
		ResolverMean:      10 * time.Millisecond,
		ResolverStdDev:    100 * time.Millisecond,
		Concurrency:       concurrency,
	}
	sb := NewSourceBenchmark(b, bs)
	sb.Run()
}
