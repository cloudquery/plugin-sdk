package source

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type BenchmarkScenario struct {
	Client                Client
	Scheduler             specs.Scheduler
	Clients               int
	Tables                int
	ChildrenPerTable      int
	Columns               int
	ColumnResolvers       int // number of columns with custom resolvers
	ResourcesPerTable     int
	ResourcesPerPage      int
	NoPreResourceResolver bool
	Concurrency           uint64
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
}

type Client interface {
	Call(clientID, tableName string) error
}

type Benchmark struct {
	*BenchmarkScenario

	b      *testing.B
	tables []*schema.Table
	plugin *Plugin

	apiCalls atomic.Int64
}

func NewBenchmark(b *testing.B, scenario BenchmarkScenario) *Benchmark {
	scenario.SetDefaults()
	sb := &Benchmark{
		BenchmarkScenario: &scenario,
		b:                 b,
		tables:            nil,
		plugin:            nil,
	}
	sb.setup(b)
	return sb
}

func (s *Benchmark) setup(b *testing.B) {
	createResolvers := func(tableName string) (schema.TableResolver, schema.RowResolver, schema.ColumnResolver) {
		tableResolver := func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
			total := 0
			for total < s.ResourcesPerTable {
				s.simulateAPICall(meta.ID(), tableName)
				num := min(s.ResourcesPerPage, s.ResourcesPerTable-total)
				resources := make([]struct {
					Column1 string
				}, num)
				for i := 0; i < num; i++ {
					resources[i] = struct {
						Column1 string
					}{
						Column1: "test-column",
					}
				}
				res <- resources
				total += num
			}
			return nil
		}
		preResourceResolver := func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
			s.simulateAPICall(meta.ID(), tableName)
			resource.Item = struct {
				Column1 string
			}{
				Column1: "test-pre",
			}
			return nil
		}
		columnResolver := func(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
			s.simulateAPICall(meta.ID(), tableName)
			return resource.Set(c.Name, "test")
		}
		return tableResolver, preResourceResolver, columnResolver
	}

	s.tables = make([]*schema.Table, s.Tables)
	for i := 0; i < s.Tables; i++ {
		tableResolver, preResourceResolver, columnResolver := createResolvers(fmt.Sprintf("table%d", i))
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
				Name:     fmt.Sprintf("table%d_child%d", i, u),
				Columns:  columns,
				Resolver: tableResolver,
			}
			if !s.NoPreResourceResolver {
				relations[u].PreResourceResolver = preResourceResolver
			}
		}
		s.tables[i] = &schema.Table{
			Name:      fmt.Sprintf("table%d", i),
			Columns:   columns,
			Relations: relations,
			Resolver:  tableResolver,
			Multiplex: nMultiplexer(s.Clients),
		}
		if !s.NoPreResourceResolver {
			s.tables[i].PreResourceResolver = preResourceResolver
		}
		for u := range relations {
			relations[u].Parent = s.tables[i]
		}
	}

	plugin := NewPlugin(
		"testPlugin",
		"1.0.0",
		s.tables,
		newTestExecutionClient,
	)
	plugin.SetLogger(zerolog.New(zerolog.NewTestWriter(b)).Level(zerolog.WarnLevel))
	s.plugin = plugin
	s.b = b
}

func (s *Benchmark) simulateAPICall(clientID, tableName string) {
	for {
		s.apiCalls.Add(1)
		err := s.Client.Call(clientID, tableName)
		if err == nil {
			// if no error, we are done
			break
		}
		// if error, we have to retry
		// we simulate a random backoff
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *Benchmark) Run() {
	for n := 0; n < s.b.N; n++ {
		s.b.StopTimer()
		ctx := context.Background()
		spec := specs.Source{
			Name:         "testSource",
			Path:         "cloudquery/testSource",
			Tables:       []string{"*"},
			Version:      "v1.0.0",
			Destinations: []string{"test"},
			Concurrency:  s.Concurrency,
			Scheduler:    s.Scheduler,
		}
		if err := s.plugin.Init(ctx, spec); err != nil {
			s.b.Fatal(err)
		}
		resources := make(chan *schema.Resource)
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			defer close(resources)
			return s.plugin.Sync(ctx,
				resources)
		})
		s.b.StartTimer()
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
		s.b.ReportMetric(0, "ns/op") // drop default ns/op output
		s.b.ReportMetric(float64(totalResources)/(end.Sub(start).Seconds()), "resources/s")

		// Enable the below metrics for more verbose information about the scenario:
		// s.b.ReportMetric(float64(s.apiCalls.Load())/(end.Sub(start).Seconds()), "api-calls/s")
		//  s.b.ReportMetric(float64(totalResources), "resources")
		// s.b.ReportMetric(float64(s.apiCalls.Load()), "apiCalls")
	}
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

func BenchmarkDefaultConcurrencyDFS(b *testing.B) {
	benchmarkWithScheduler(b, specs.SchedulerDFS)
}

func BenchmarkDefaultConcurrencyRoundRobin(b *testing.B) {
	benchmarkWithScheduler(b, specs.SchedulerRoundRobin)
}

func benchmarkWithScheduler(b *testing.B, scheduler specs.Scheduler) {
	b.ReportAllocs()
	minTime := 1 * time.Millisecond
	mean := 10 * time.Millisecond
	stdDev := 100 * time.Millisecond
	client := NewDefaultClient(minTime, mean, stdDev)
	bs := BenchmarkScenario{
		Client:            client,
		Clients:           25,
		Tables:            5,
		Columns:           10,
		ColumnResolvers:   1,
		ResourcesPerTable: 100,
		ResourcesPerPage:  50,
		Scheduler:         scheduler,
	}
	sb := NewBenchmark(b, bs)
	sb.Run()
}

func BenchmarkTablesWithChildrenDFS(b *testing.B) {
	benchmarkTablesWithChildrenScheduler(b, specs.SchedulerDFS)
}

func BenchmarkTablesWithChildrenRoundRobin(b *testing.B) {
	benchmarkTablesWithChildrenScheduler(b, specs.SchedulerRoundRobin)
}

func benchmarkTablesWithChildrenScheduler(b *testing.B, scheduler specs.Scheduler) {
	b.ReportAllocs()
	minTime := 1 * time.Millisecond
	mean := 10 * time.Millisecond
	stdDev := 100 * time.Millisecond
	client := NewDefaultClient(minTime, mean, stdDev)
	bs := BenchmarkScenario{
		Client:            client,
		Clients:           2,
		Tables:            2,
		ChildrenPerTable:  2,
		Columns:           10,
		ColumnResolvers:   1,
		ResourcesPerTable: 100,
		ResourcesPerPage:  50,
		Scheduler:         scheduler,
	}
	sb := NewBenchmark(b, bs)
	sb.Run()
}

type DefaultClient struct {
	min, stdDev, mean time.Duration
}

func NewDefaultClient(min, mean, stdDev time.Duration) *DefaultClient {
	if min == 0 {
		min = time.Millisecond
	}
	if mean == 0 {
		mean = 10 * time.Millisecond
	}
	if stdDev == 0 {
		stdDev = 100 * time.Millisecond
	}
	return &DefaultClient{
		min:    min,
		mean:   mean,
		stdDev: stdDev,
	}
}

func (c *DefaultClient) Call(_, _ string) error {
	sample := int(rand.NormFloat64()*float64(c.stdDev) + float64(c.mean))
	duration := time.Duration(sample)
	if duration < c.min {
		duration = c.min
	}
	time.Sleep(duration)
	return nil
}

type RateLimitClient struct {
	*DefaultClient
	calls             map[string][]time.Time
	callsLock         sync.Mutex
	window            time.Duration
	maxCallsPerWindow int
}

func NewRateLimitClient(min, mean, stdDev time.Duration, maxCallsPerWindow int, window time.Duration) *RateLimitClient {
	return &RateLimitClient{
		DefaultClient:     NewDefaultClient(min, mean, stdDev),
		calls:             map[string][]time.Time{},
		window:            window,
		maxCallsPerWindow: maxCallsPerWindow,
	}
}

func (r *RateLimitClient) Call(clientID, table string) error {
	// this will sleep for the appropriate amount of time before responding
	err := r.DefaultClient.Call(clientID, table)
	if err != nil {
		return err
	}

	r.callsLock.Lock()
	defer r.callsLock.Unlock()

	// limit the number of calls per window by table
	key := table

	// remove calls from outside the call window
	updated := make([]time.Time, 0, len(r.calls[key]))
	for i := range r.calls[key] {
		if time.Since(r.calls[key][i]) < r.window {
			updated = append(updated, r.calls[key][i])
		}
	}

	// return error if we've exceeded the max calls in the time window
	if len(updated) >= r.maxCallsPerWindow {
		return fmt.Errorf("rate limit exceeded")
	}

	r.calls[key] = append(r.calls[key], time.Now())
	return nil
}

// BenchmarkDefaultConcurrency represents a benchmark scenario where rate limiting is applied
// by the cloud provider. In this rate limiter, the limit is applied globally per table.
// This mirrors the behavior of GCP, where rate limiting is applied per project *token*, not
// per project. A good scheduler should spread the load across tables so that other tables can make
// progress while waiting for the rate limit to reset.
func BenchmarkTablesWithRateLimitingDFS(b *testing.B) {
	benchmarkTablesWithRateLimitingScheduler(b, specs.SchedulerDFS)
}

func BenchmarkTablesWithRateLimitingRoundRobin(b *testing.B) {
	benchmarkTablesWithRateLimitingScheduler(b, specs.SchedulerRoundRobin)
}

// In this benchmark, we set up a scenario where each table has a global rate limit of 1 call per 100ms.
// Every table requires 1 call to resolve, and has 10 clients. This means, at best, each table can resolve in 1 second.
// We have 100 such tables and a concurrency that allows 1000 calls at a time. A good scheduler for this scenario
// should be able to resolve all tables in a bit more than 1 second.
func benchmarkTablesWithRateLimitingScheduler(b *testing.B, scheduler specs.Scheduler) {
	b.ReportAllocs()
	minTime := 1 * time.Millisecond
	mean := 1 * time.Millisecond
	stdDev := 1 * time.Millisecond
	maxCallsPerWindow := 1
	window := 100 * time.Millisecond
	c := NewRateLimitClient(minTime, mean, stdDev, maxCallsPerWindow, window)

	bs := BenchmarkScenario{
		Client:                c,
		Scheduler:             scheduler,
		Clients:               10,
		Tables:                100,
		ChildrenPerTable:      0,
		Columns:               10,
		ColumnResolvers:       0,
		ResourcesPerTable:     1,
		ResourcesPerPage:      1,
		Concurrency:           1000,
		NoPreResourceResolver: true,
	}
	sb := NewBenchmark(b, bs)
	sb.Run()
}
