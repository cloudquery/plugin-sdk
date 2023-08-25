package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
	"go.opentelemetry.io/otel"
	"golang.org/x/sync/semaphore"
)

const (
	DefaultConcurrency     = 50000
	DefaultRowsPerRecord   = 100
	DefaultMaxDepth        = 4
	minTableConcurrency    = 1
	minResourceConcurrency = 100
	minRowsPerRecord       = 1
	otelName               = "schedule"
)

const (
	StrategyDFS Strategy = iota
	StrategyRoundRobin
)

type Strategy int

func (s *Strategy) String() string {
	if s == nil {
		return ""
	}
	return AllStrategyNames[*s]
}

// MarshalJSON implements json.Marshaler.
func (s *Strategy) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	if s == nil {
		b.Write([]byte("null"))
		return b.Bytes(), nil
	}
	b.Write([]byte{'"'})
	b.Write([]byte(s.String()))
	b.Write([]byte{'"'})
	return b.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Strategy) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}
	strategy, err := StrategyForName(name)
	if err != nil {
		return err
	}
	*s = strategy
	return nil
}

func (s *Strategy) Validate() error {
	if s == nil {
		return errors.New("scheduler strategy is nil")
	}
	for _, strategy := range AllStrategies {
		if strategy == *s {
			return nil
		}
	}
	return fmt.Errorf("unknown scheduler strategy: %d", s)
}

var AllStrategies = Strategies{StrategyDFS, StrategyRoundRobin}
var AllStrategyNames = [...]string{
	StrategyDFS:        "dfs",
	StrategyRoundRobin: "round-robin",
}

func StrategyForName(s string) (Strategy, error) {
	for i, name := range AllStrategyNames {
		if name == s {
			return AllStrategies[i], nil
		}
	}
	return StrategyDFS, fmt.Errorf("unknown scheduler strategy: %s", s)
}

type Strategies []Strategy

func (s Strategies) String() string {
	var buffer bytes.Buffer
	for i, strategy := range s {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(strategy.String())
	}
	return buffer.String()
}

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

func WithRowsPerRecord(rowsPerRecord int) Option {
	return func(s *Scheduler) {
		s.rowsPerRecord = max(rowsPerRecord, minRowsPerRecord)
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

type SyncOption func(*syncClient)

func WithSyncDeterministicCQID(deterministicCQID bool) SyncOption {
	return func(s *syncClient) {
		s.deterministicCQID = deterministicCQID
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
	logger        zerolog.Logger
	concurrency   int
	rowsPerRecord int
}

type syncClient struct {
	tables            schema.Tables
	client            schema.ClientMeta
	scheduler         *Scheduler
	deterministicCQID bool
	// status sync metrics
	metrics *Metrics
	logger  zerolog.Logger
}

func NewScheduler(opts ...Option) *Scheduler {
	s := Scheduler{
		caser:         caser.New(),
		concurrency:   DefaultConcurrency,
		maxDepth:      DefaultMaxDepth,
		rowsPerRecord: DefaultRowsPerRecord,
	}
	for _, opt := range opts {
		opt(&s)
	}
	// This is very similar to the concurrent web crawler problem with some minor changes.
	// We are using DFS/Round-Robin to make sure memory usage is capped at O(h) where h is the height of the tree.
	tableConcurrency := max(s.concurrency/minResourceConcurrency, minTableConcurrency)
	resourceConcurrency := tableConcurrency * minResourceConcurrency
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
	ctx, span := otel.Tracer(otelName).Start(ctx, "Sync")
	defer span.End()
	if len(tables) == 0 {
		return nil
	}

	syncClient := &syncClient{
		metrics:   &Metrics{TableClient: make(map[string]map[string]*TableClientMetrics)},
		tables:    tables,
		client:    client,
		scheduler: s,
		logger:    s.logger,
	}
	for _, opt := range opts {
		opt(syncClient)
	}

	if maxDepth(tables) > s.maxDepth {
		return fmt.Errorf("max depth exceeded, max depth is %d", s.maxDepth)
	}

	// send migrate messages first
	for _, table := range tables.FlattenTables() {
		res <- &message.SyncMigrateTable{
			Table: table,
		}
	}

	resources := make(chan *schema.Resource)
	go func() {
		defer close(resources)
		switch s.strategy {
		case StrategyDFS:
			syncClient.syncDfs(ctx, resources)
		case StrategyRoundRobin:
			syncClient.syncRoundRobin(ctx, resources)
		default:
			panic(fmt.Errorf("unknown scheduler %s", s.strategy.String()))
		}
	}()
	s.collectSend(resources, res)
	return nil
}

func (s *syncClient) logTablesMetrics(tables schema.Tables, client Client) {
	clientName := client.ID()
	for _, table := range tables {
		metrics := s.metrics.TableClient[table.Name][clientName]
		s.logger.Info().Str("table", table.Name).Str("client", clientName).Uint64("resources", metrics.Resources).Uint64("errors", metrics.Errors).Msg("table sync finished")
		s.logTablesMetrics(table.Relations, client)
	}
}

func (s *syncClient) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, item any) *schema.Resource {
	var validationErr *schema.ValidationError
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	resource := schema.NewResourceData(table, parent, item)
	objectStartTime := time.Now()
	clientID := client.ID()
	tableMetrics := s.metrics.TableClient[table.Name][clientID]
	logger := s.logger.With().Str("table", table.Name).Str("client", clientID).Logger()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Interface("error", err).TimeDiff("duration", time.Now(), objectStartTime).Str("stack", stack).Msg("resource resolver finished with panic")
			atomic.AddUint64(&tableMetrics.Panics, 1)
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("table", table.Name)
				sentry.CurrentHub().CaptureMessage(stack)
			})
		}
	}()
	if table.PreResourceResolver != nil {
		if err := table.PreResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Err(err).Msg("pre resource resolver failed")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			if errors.As(err, &validationErr) {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", table.Name)
					sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
				})
			}
			return nil
		}
	}

	for _, c := range table.Columns {
		s.resolveColumn(ctx, logger, tableMetrics, client, resource, c)
	}

	if table.PostResourceResolver != nil {
		if err := table.PostResourceResolver(ctx, client, resource); err != nil {
			logger.Error().Stack().Err(err).Msg("post resource resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			if errors.As(err, &validationErr) {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", table.Name)
					sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
				})
			}
		}
	}
	atomic.AddUint64(&tableMetrics.Resources, 1)
	return resource
}

func (s *syncClient) resolveColumn(ctx context.Context, logger zerolog.Logger, tableMetrics *TableClientMetrics, client schema.ClientMeta, resource *schema.Resource, c schema.Column) {
	var validationErr *schema.ValidationError
	columnStartTime := time.Now()
	defer func() {
		if err := recover(); err != nil {
			stack := fmt.Sprintf("%s\n%s", err, string(debug.Stack()))
			logger.Error().Str("column", c.Name).Interface("error", err).TimeDiff("duration", time.Now(), columnStartTime).Str("stack", stack).Msg("column resolver finished with panic")
			atomic.AddUint64(&tableMetrics.Panics, 1)
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("table", resource.Table.Name)
				scope.SetTag("column", c.Name)
				sentry.CurrentHub().CaptureMessage(stack)
			})
		}
	}()

	if c.Resolver != nil {
		if err := c.Resolver(ctx, client, resource, c); err != nil {
			logger.Error().Err(err).Msg("column resolver finished with error")
			atomic.AddUint64(&tableMetrics.Errors, 1)
			if errors.As(err, &validationErr) {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("table", resource.Table.Name)
					scope.SetTag("column", c.Name)
					sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
				})
			}
		}
	} else {
		// base use case: try to get column with CamelCase name
		v := funk.Get(resource.GetItem(), s.scheduler.caser.ToPascal(c.Name), funk.WithAllowZero())
		if v != nil {
			err := resource.Set(c.Name, v)
			if err != nil {
				logger.Error().Err(err).Msg("column resolver finished with error")
				atomic.AddUint64(&tableMetrics.Errors, 1)
				if errors.As(err, &validationErr) {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetTag("table", resource.Table.Name)
						scope.SetTag("column", c.Name)
						sentry.CurrentHub().CaptureMessage(validationErr.MaskedError())
					})
				}
			}
		}
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

// unparam's suggestion to remove the second parameter is not good advice here.
// nolint:unparam
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
