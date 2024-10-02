package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/metrics"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/semaphore"
)

const (
	DefaultSingleResourceMaxConcurrency    = 5
	DefaultSingleNestedTableMaxConcurrency = 5
	DefaultConcurrency                     = 50000
	DefaultMaxDepth                        = 4
	minTableConcurrency                    = 1
	minResourceConcurrency                 = 100
)

var ErrNoTables = errors.New("no tables specified for syncing, review `tables` and `skip_tables` in your config and specify at least one table to sync")

const (
	StrategyDFS Strategy = iota
	StrategyRoundRobin
	StrategyShuffle
	StrategyRandomQueue
)

type Option func(*Scheduler)

func WithLogger(logger zerolog.Logger) Option {
	return func(s *Scheduler) {
		s.logger = logger
	}
}

func WithConcurrency(concurrency int) Option {
	return func(s *Scheduler) {
		s.concurrency = concurrency
	}
}

func WithMaxDepth(maxDepth uint64) Option {
	return func(s *Scheduler) {
		s.maxDepth = maxDepth
	}
}

func WithStrategy(strategy Strategy) Option {
	return func(s *Scheduler) {
		s.strategy = strategy
	}
}

func WithSingleNestedTableMaxConcurrency(concurrency int64) Option {
	return func(s *Scheduler) {
		s.singleNestedTableMaxConcurrency = concurrency
	}
}

func WithSingleResourceMaxConcurrency(concurrency int64) Option {
	return func(s *Scheduler) {
		s.singleResourceMaxConcurrency = concurrency
	}
}

type SyncOption func(*syncClient)

func WithSyncDeterministicCQID(deterministicCQID bool) SyncOption {
	return func(s *syncClient) {
		s.deterministicCQID = deterministicCQID
	}
}

func WithInvocationID(invocationID string) Option {
	return func(s *Scheduler) {
		s.invocationID = invocationID
	}
}

func WithShard(num int32, total int32) SyncOption {
	return func(s *syncClient) {
		s.shard = &shard{num: num, total: total}
	}
}

type Client interface {
	ID() string
}

type Scheduler struct {
	caser    *caser.Caser
	strategy Strategy
	maxDepth uint64
	// resourceSem is a semaphore that limits the number of concurrent resources being fetched
	resourceSem *semaphore.Weighted
	// tableSem is a semaphore that limits the number of concurrent tables being fetched
	tableSems []*semaphore.Weighted
	// Logger to call, this logger is passed to the serve.Serve Client, if not defined Serve will create one instead.
	logger      zerolog.Logger
	concurrency int
	// This Map holds all of the concurrency semaphores for each table+client pair.
	singleTableConcurrency sync.Map
	// The maximum number of go routines that can be spawned for a single table+client pair
	singleNestedTableMaxConcurrency int64

	// The maximum number of go routines that can be spawned for a specific resource
	singleResourceMaxConcurrency int64

	// Controls how records are constructed on the source side.
	batchSettings *BatchSettings

	invocationID string
}

type shard struct {
	num   int32
	total int32
}

type syncClient struct {
	tables            schema.Tables
	client            schema.ClientMeta
	scheduler         *Scheduler
	deterministicCQID bool
	// status sync metrics
	metrics      *metrics.Metrics
	logger       zerolog.Logger
	invocationID string

	shard *shard
}

func NewScheduler(opts ...Option) *Scheduler {
	s := Scheduler{
		caser:                           caser.New(),
		concurrency:                     DefaultConcurrency,
		maxDepth:                        DefaultMaxDepth,
		singleResourceMaxConcurrency:    DefaultSingleResourceMaxConcurrency,
		singleNestedTableMaxConcurrency: DefaultSingleNestedTableMaxConcurrency,
		batchSettings: &BatchSettings{
			MaxRows: DefaultBatchMaxRows,
			Timeout: DefaultBatchTimeout,
		},
	}
	for _, opt := range opts {
		opt(&s)
	}

	actualMinResourceConcurrency := minResourceConcurrency
	if s.concurrency <= minResourceConcurrency {
		actualMinResourceConcurrency = max(s.concurrency/10, 1)
	}

	// This is very similar to the concurrent web crawler problem with some minor changes.
	// We are using DFS/Round-Robin to make sure memory usage is capped at O(h) where h is the height of the tree.
	tableConcurrency := max(s.concurrency/actualMinResourceConcurrency, minTableConcurrency)
	resourceConcurrency := tableConcurrency * actualMinResourceConcurrency
	s.tableSems = make([]*semaphore.Weighted, s.maxDepth)
	for i := uint64(0); i < s.maxDepth; i++ {
		s.tableSems[i] = semaphore.NewWeighted(int64(tableConcurrency))
		// reduce table concurrency logarithmically for every depth level
		tableConcurrency = max(tableConcurrency/2, minTableConcurrency)
	}
	s.resourceSem = semaphore.NewWeighted(int64(resourceConcurrency))

	return &s
}

// SyncAll is mostly used for testing as it will sync all tables and can run out of memory
// in the real world. Should use Sync for production.
func (s *Scheduler) SyncAll(ctx context.Context, client schema.ClientMeta, tables schema.Tables) (message.SyncMessages, error) {
	res := make(chan message.SyncMessage)
	var err error
	go func() {
		defer close(res)
		err = s.Sync(ctx, client, tables, res)
	}()
	// nolint:prealloc
	var messages message.SyncMessages
	for msg := range res {
		messages = append(messages, msg)
	}
	return messages, err
}

func (s *Scheduler) Sync(ctx context.Context, client schema.ClientMeta, tables schema.Tables, res chan<- message.SyncMessage, opts ...SyncOption) error {
	ctx, span := otel.Tracer(metrics.OtelName).Start(ctx,
		"sync",
		trace.WithAttributes(attribute.Key("sync.invocation.id").String(s.invocationID)),
	)
	defer span.End()
	if len(tables) == 0 {
		return ErrNoTables
	}

	syncClient := &syncClient{
		metrics:      &metrics.Metrics{TableClient: make(map[string]map[string]*metrics.TableClientMetrics)},
		tables:       tables,
		client:       client,
		scheduler:    s,
		logger:       s.logger,
		invocationID: s.invocationID,
	}
	for _, opt := range opts {
		opt(syncClient)
	}

	if maxDepth(tables) > s.maxDepth {
		return fmt.Errorf("max depth exceeded, max depth is %d", s.maxDepth)
	}

	// send migrate messages first
	for _, tableOriginal := range tables.FlattenTables() {
		migrateMessage := &message.SyncMigrateTable{
			Table: tableOriginal.Copy(tableOriginal.Parent),
		}
		if syncClient.deterministicCQID {
			schema.CqIDAsPK(migrateMessage.Table)
		}
		res <- migrateMessage
	}

	resources := make(chan *schema.Resource)
	go func() {
		defer close(resources)
		testMultiplier, err := getTestMultiplier()
		if err != nil {
			panic(err)
		}
		if testMultiplier > 0 {
			syncClient.syncTest(ctx, testMultiplier, resources)
			return
		}
		switch s.strategy {
		case StrategyDFS:
			syncClient.syncDfs(ctx, resources)
		case StrategyRoundRobin:
			syncClient.syncRoundRobin(ctx, resources)
		case StrategyShuffle:
			syncClient.syncShuffle(ctx, resources)
		case StrategyRandomQueue:
			syncClient.syncRandomQueue(ctx, resources)
		default:
			panic(fmt.Errorf("unknown scheduler %s", s.strategy.String()))
		}
	}()

	b := s.batchSettings.getBatcher(ctx, res, s.logger)
	defer b.close()    // wait for all resources to be processed
	done := ctx.Done() // no need to do the lookups in loop
	for resource := range resources {
		select {
		case <-done:
			s.logger.Debug().Msg("sync context cancelled")
			return context.Cause(ctx)
		default:
			b.process(resource)
		}
	}
	return context.Cause(ctx)
}

func (s *syncClient) logTablesMetrics(tables schema.Tables, client Client) {
	clientName := client.ID()
	for _, table := range tables {
		m := s.metrics.TableClient[table.Name][clientName]
		duration := m.Duration.Load()
		if duration == nil {
			// This can happen for a relation when there are no resources to resolve from the parent
			duration = new(time.Duration)
		}
		s.logger.Info().Str("table", table.Name).Str("client", clientName).Uint64("resources", m.Resources).Dur("duration_ms", *duration).Uint64("errors", m.Errors).Msg("table sync finished")
		s.logTablesMetrics(table.Relations, client)
	}
}

func maxDepth(tables schema.Tables) uint64 {
	var depth uint64
	if len(tables) == 0 {
		return 0
	}
	for _, table := range tables {
		newDepth := 1 + maxDepth(table.Relations)
		if newDepth > depth {
			depth = newDepth
		}
	}
	return depth
}

func shardTableClients(tableClients []tableClient, shard *shard) []tableClient {
	// For sharding to work as expected, tableClients must be deterministic between different shards.
	if shard == nil || len(tableClients) == 0 {
		return tableClients
	}
	num := int(shard.num)
	total := int(shard.total)
	chunkSize := len(tableClients) / total
	if chunkSize == 0 {
		chunkSize = 1
	}
	chunks := lo.Chunk(tableClients, chunkSize)
	if num > len(chunks) {
		return nil
	}
	if len(chunks) > total && num == total {
		return append(chunks[num-1], chunks[num]...)
	}
	return chunks[num-1]
}
