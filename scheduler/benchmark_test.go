package scheduler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"

	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type BenchmarkScenario struct {
	Client                 *sync.Map
	ClientInit             func() Client
	Scheduler              scheduler.Strategy
	Clients                int
	Tables                 int
	ChildrenPerTable       int
	Columns                int
	ColumnResolvers        int // number of columns with custom resolvers
	ResourcesPerTable      int
	ResourcesPerPage       int
	NoPreResourceResolver  bool
	Concurrency            uint64
	SingleTableConcurrency int
	MaxRetries             int
	GlobalRateLimiter      bool
}

func defaultBenchmarkScenario() BenchmarkScenario {
	return BenchmarkScenario{
		Client:                 &sync.Map{},
		Clients:                1,
		Tables:                 1,
		Columns:                10,
		ColumnResolvers:        1,
		ResourcesPerTable:      50,
		ResourcesPerPage:       10,
		MaxRetries:             5,
		Concurrency:            50000,
		SingleTableConcurrency: 50000,
	}
}

type Client interface {
	Call(clientID, tableName string) error
}

type Benchmark struct {
	*BenchmarkScenario

	b                 *testing.B
	tables            []*schema.Table
	plugin            *plugin.Plugin
	failedApiCalls    atomic.Int64
	timeSpentSleeping atomic.Int64
	succeededApiCalls atomic.Int64
	throttles         atomic.Int64
}

func NewBenchmark(b *testing.B, scenario BenchmarkScenario) *Benchmark {
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

	plugin := plugin.NewPlugin(
		"testPlugin",
		"1.0.0",
		s.Configure,
	)
	plugin.SetLogger(zerolog.New(zerolog.NewTestWriter(b)).Level(zerolog.WarnLevel))
	s.plugin = plugin
	s.b = b
}

type PluginClient struct {
	plugin.UnimplementedDestination
	scheduler *scheduler.Scheduler
	logger    zerolog.Logger
	options   plugin.NewClientOptions
	allTables schema.Tables
}

func (*PluginClient) Close(_ context.Context) error {
	return nil
}

func (c *PluginClient) Tables(_ context.Context, options plugin.TableOptions) (schema.Tables, error) {
	return c.allTables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
}
func (c *PluginClient) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	tt, err := c.allTables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return err
	}

	return c.scheduler.Sync(ctx, nil, tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
}
func (s *Benchmark) Configure(ctx context.Context, logger zerolog.Logger, specBytes []byte, options plugin.NewClientOptions) (plugin.Client, error) {
	c := &PluginClient{
		options: options,
		logger:  logger,
	}

	c.scheduler = scheduler.NewScheduler(
		scheduler.WithConcurrency(int(s.Concurrency)),
		scheduler.WithLogger(c.logger),
		scheduler.WithStrategy(s.Scheduler),
		// scheduler.WithSingleTableMaxConcurrency(s.SingleTableConcurrency),
	)

	createResolvers := func(tableName string, depth int) (schema.TableResolver, schema.RowResolver, schema.ColumnResolver) {
		tableResolver := func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
			total := 0
			ResourcesPerPage := s.ResourcesPerPage
			if depth > 0 {
				ResourcesPerPage = s.ResourcesPerTable
			}
			for total < s.ResourcesPerTable {
				s.simulateAPICall(meta.ID(), tableName)
				num := min(ResourcesPerPage, s.ResourcesPerTable-total)
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
		tableResolver, preResourceResolver, columnResolver := createResolvers(fmt.Sprintf("table%d", i), i)
		columns := make([]schema.Column, s.Columns)
		for u := 0; u < s.Columns; u++ {
			columns[u] = schema.Column{
				Name: fmt.Sprintf("column%d", u),
				Type: arrow.BinaryTypes.String,
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
	c.allTables = s.tables
	return c, nil
}

func (s *Benchmark) simulateAPICall(clientID, tableName string) {
	retries := 0
	for {
		if retries > s.MaxRetries {
			s.failedApiCalls.Add(1)
			return
		}
		key := clientID + "-" + tableName
		if s.GlobalRateLimiter {
			key = "global"
		}

		client, _ := s.Client.LoadOrStore(key, s.ClientInit())
		err := client.(Client).Call(clientID, tableName)
		if err == nil {
			// if no error, we are done
			s.succeededApiCalls.Add(1)
			break
		}
		s.throttles.Add(1)
		retries++
		// if error, we have to retry
		// we simulate a random backoff

		sleepDur := s.calculateBackoff(retries)
		s.timeSpentSleeping.Add(int64(sleepDur.Seconds()))
		time.Sleep(sleepDur)

	}
}

func (s *Benchmark) calculateBackoff(retry int) time.Duration {
	backoffDuration := time.Duration(float64(1.2)*math.Pow(float64(1.5), float64(retry))) * time.Second
	if backoffDuration > time.Duration(15*time.Second) {
		backoffDuration = 15 * time.Second
	}
	return backoffDuration
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
		}
		// Marshal spec into []byte
		specBytes, _ := json.Marshal(spec)

		if err := s.plugin.Init(ctx, specBytes, plugin.NewClientOptions{}); err != nil {
			s.b.Fatal(err)
		}
		resources := make(chan message.SyncMessage)
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			defer close(resources)

			return s.plugin.Sync(ctx,
				plugin.SyncOptions{
					Tables: []string{"*"},
				},
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
		s.b.ReportMetric(0, "ns/op")     // drop default ns/op output
		s.b.ReportMetric(0, "B/op")      // drop default B/op output
		s.b.ReportMetric(0, "allocs/op") // drop default allocs/op output

		s.b.ReportMetric(end.Sub(start).Seconds(), "duration(seconds)")
		s.b.ReportMetric(float64(totalResources)/(end.Sub(start).Seconds()), "resources/s")

		// Enable the below metrics for more verbose information about the scenario:
		s.b.ReportMetric((float64(s.succeededApiCalls.Load()))/(end.Sub(start).Seconds()), "successful-api-calls/s")
		s.b.ReportMetric(float64(totalResources), "resources")
		s.b.ReportMetric(float64(s.succeededApiCalls.Load()), "succeededApiCalls")
		s.b.ReportMetric(float64(s.failedApiCalls.Load()), "failedApiCalls")
		s.b.ReportMetric(float64(s.throttles.Load()), "throttles")
		s.b.ReportMetric(float64(s.timeSpentSleeping.Load()), "total-time-spent-sleeping")
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

func runBenchmark(b *testing.B, options ...TestOptions) {
	// b.ReportAllocs()

	bs := defaultBenchmarkScenario()

	for _, option := range options {
		option(&bs)
	}
	sb := NewBenchmark(b, bs)
	sb.Run()
}

func benchmarkTablesWithChildrenScheduler(b *testing.B, scheduler scheduler.Strategy, options ...TestOptions) {
	// b.ReportAllocs()
	minTime := 1 * time.Millisecond
	mean := 10 * time.Millisecond
	stdDev := 100 * time.Millisecond
	bs := defaultBenchmarkScenario()
	bs.ClientInit = func() Client { return NewDefaultClient(minTime, mean, stdDev) }
	bs.ChildrenPerTable = 2
	bs.Scheduler = scheduler
	for _, option := range options {
		option(&bs)
	}
	sb := NewBenchmark(b, bs)
	sb.Run()
}

type DefaultClient struct {
	min, stdDev, mean time.Duration
}

func NewDefaultClient(min, mean, stdDev time.Duration) *DefaultClient {
	if min == 0 {
		min = 1 * time.Millisecond
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

func (c *DefaultClient) Call(clientID, Table string) error {
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
	global            bool
}

func NewGlobalRateLimitClient(min, mean, stdDev time.Duration, maxCallsPerWindow int, window time.Duration) *RateLimitClient {
	return &RateLimitClient{
		DefaultClient:     NewDefaultClient(min, mean, stdDev),
		calls:             map[string][]time.Time{},
		window:            window,
		maxCallsPerWindow: maxCallsPerWindow,
		global:            true,
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

func BenchmarkTablesWithGlobalRateLimiting(b *testing.B) {
	for _, strategy := range scheduler.AllStrategies {
		for _, concurrency := range []int{10000, 1000, 500} {
			b.Run(fmt.Sprintf("%s-%d", strategy.String(), concurrency), func(b *testing.B) {
				runBenchmark(b,
					WithScheduler(strategy),
					WithClientInit(func() Client {
						return NewGlobalRateLimitClient(50*time.Millisecond, 250*time.Millisecond, 50*time.Millisecond, 3, 500*time.Millisecond)
					}),
					WithGlobalRateLimiting(true),
					WithClients(10),
					WithTables(100),
					WithScheduler(strategy),
					WithColumnResolvers(0),
					WithChildTables(0),
					WithNoPreResourceResolver(),
					WithConcurrency(uint64(concurrency)),
				)
			})
		}
	}
}

// BenchmarkTablesWithTableClientRateLimiting represents a benchmark scenario where rate limiting is applied on a
// by the cloud provider. In this rate limiter, the limit is applied on a per table + client basis. It makes the assumption that each client + table pair have separate rate limits
// This mirrors the behavior of AWS, where rate limiting is applied per account, region and table. This will help test nested tables

// In this benchmark, we set up a scenario where each table has a rate limit of 6 call per second
// Each API call takes around 250ms to resolve and backoff/retries are calculated at `(1.2 * 1.5 ^ retry) * seconds`.
// It takes 255 API requests to fully resolve this scenario. At a theoretical rate of 6 calls per second the fastest this could ever resolve is 42.5 seconds, but this assumes no throttling or time waiting for backoff
// The rate limiter is applied

// Strategy is not included in this benchmark because it is not relevant as there is only a single top level table
func BenchmarkTablesWithTableClientRateLimiting(b *testing.B) {
	for _, concurrency := range []int{10000, 1000, 500} {
		b.Run(fmt.Sprintf("concurrency-%d", concurrency), func(b *testing.B) {
			runBenchmark(b,
				WithClientInit(func() Client {
					return NewGlobalRateLimitClient(50*time.Millisecond, 250*time.Millisecond, 50*time.Millisecond, 3, 500*time.Millisecond)
				}),
				WithGlobalRateLimiting(false),
				WithClients(1),
				WithTables(1),
				WithColumnResolvers(0),
				WithChildTables(1),
				WithNoPreResourceResolver(),
				WithConcurrency(uint64(concurrency)),
				// WithSingleTableMaxConcurrency(3),

			)
		})
	}
}

func BenchmarkDefaultConcurrency(b *testing.B) {
	for _, strategy := range scheduler.AllStrategies {
		b.Run(strategy.String(), func(b *testing.B) {
			runBenchmark(b,
				WithScheduler(strategy),
				WithClientInit(func() Client { return NewDefaultClient(1*time.Millisecond, 10*time.Millisecond, 100*time.Millisecond) }))
		})
	}
}

func BenchmarkTablesWithChildren(b *testing.B) {
	for _, strategy := range scheduler.AllStrategies {
		for _, concurrency := range []int{1000, 10, 1} {
			b.Run(fmt.Sprintf("%s-%d", strategy.String(), concurrency), func(b *testing.B) {
				benchmarkTablesWithChildrenScheduler(b, strategy,
					WithClientInit(func() Client { return NewDefaultClient(1*time.Millisecond, 10*time.Millisecond, 100*time.Millisecond) }))
				// WithSingleTableMaxConcurrency(concurrency),
			})
		}
	}
}

type TestOptions func(*BenchmarkScenario)

func WithConcurrency(concurrency uint64) TestOptions {
	return func(s *BenchmarkScenario) {
		s.Concurrency = concurrency
	}
}

// func WithSingleTableMaxConcurrency(concurrency int) TestOptions {
// 	return func(s *BenchmarkScenario) {
// 		s.SingleTableConcurrency = concurrency
// 	}
// }

func WithGlobalRateLimiting(global bool) TestOptions {
	return func(s *BenchmarkScenario) {
		s.GlobalRateLimiter = global
	}
}

func WithTables(tables int) TestOptions {
	return func(s *BenchmarkScenario) {
		s.Tables = tables
	}
}

func WithScheduler(scheduler scheduler.Strategy) TestOptions {
	return func(s *BenchmarkScenario) {
		s.Scheduler = scheduler
	}
}

func WithClientInit(init func() Client) TestOptions {
	return func(s *BenchmarkScenario) {
		s.ClientInit = init
	}
}

func WithClients(count int) TestOptions {
	return func(s *BenchmarkScenario) {
		s.Clients = count
	}
}

func WithColumnResolvers(count int) TestOptions {
	return func(s *BenchmarkScenario) {
		s.ColumnResolvers = count
	}
}

func WithChildTables(count int) TestOptions {
	return func(s *BenchmarkScenario) {
		s.ChildrenPerTable = count
	}
}

func WithNoPreResourceResolver() TestOptions {
	return func(s *BenchmarkScenario) {
		s.NoPreResourceResolver = true
	}
}
