package scheduler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/semaphore"
)

const (
	minTableConcurrency    = 1
	minResourceConcurrency = 100
	defaultConcurrency     = 200000
	defaultMaxDepth        = 4
)

type Strategy int

const (
	StrategyDFS Strategy = iota
	StrategyRoundRobin
)

var AllSchedulers = Strategies{StrategyDFS, StrategyRoundRobin}
var AllSchedulerNames = [...]string{
	StrategyDFS:        "dfs",
	StrategyRoundRobin: "round-robin",
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

func (s Strategy) String() string {
	return AllSchedulerNames[s]
}

type Option func(*Scheduler)

func WithLogger(logger zerolog.Logger) Option {
	return func(s *Scheduler) {
		s.logger = logger
	}
}

func WithDeterministicCQId(deterministicCQId bool) Option {
	return func(s *Scheduler) {
		s.deterministicCQId = deterministicCQId
	}
}

func WithConcurrency(concurrency uint64) Option {
	return func(s *Scheduler) {
		s.concurrency = concurrency
	}
}

func WithMaxDepth(maxDepth uint64) Option {
	return func(s *Scheduler) {
		s.maxDepth = maxDepth
	}
}

func WithSchedulerStrategy(strategy Strategy) Option {
	return func(s *Scheduler) {
		s.strategy = strategy
	}
}

type Client interface {
	ID() string
}

type Scheduler struct {
	tables   schema.Tables
	client   schema.ClientMeta
	caser    *caser.Caser
	strategy Strategy
	// status sync metrics
	metrics  *Metrics
	maxDepth uint64
	// resourceSem is a semaphore that limits the number of concurrent resources being fetched
	resourceSem *semaphore.Weighted
	// tableSem is a semaphore that limits the number of concurrent tables being fetched
	tableSems []*semaphore.Weighted
	// Logger to call, this logger is passed to the serve.Serve Client, if not defined Serve will create one instead.
	logger            zerolog.Logger
	deterministicCQId bool
	concurrency       uint64
}

func NewScheduler(client schema.ClientMeta, opts ...Option) *Scheduler {
	s := Scheduler{
		client:      client,
		metrics:     &Metrics{TableClient: make(map[string]map[string]*TableClientMetrics)},
		caser:       caser.New(),
		concurrency: defaultConcurrency,
		maxDepth:    defaultMaxDepth,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

// SyncAll is mostly used for testing as it will sync all tables and can run out of memory
// in the real world. Should use Sync for production.
func (s *Scheduler) SyncAll(ctx context.Context, tables schema.Tables) (message.Messages, error) {
	res := make(chan message.Message)
	go func() {
		defer close(res)
		s.Sync(ctx, tables, res)
	}()
	var messages []message.Message
	for msg := range res {
		messages = append(messages, msg)
	}
	return messages, nil
}

func (s *Scheduler) Sync(ctx context.Context, tables schema.Tables, res chan<- message.Message) error {
	if len(tables) == 0 {
		return nil
	}

	if maxDepth(tables) > s.maxDepth {
		return fmt.Errorf("max depth exceeded, max depth is %d", s.maxDepth)
	}
	s.tables = tables

	// send migrate messages first
	for _, table := range tables.FlattenTables() {
		res <- &message.MigrateTable{
			Table: table,
		}
	}

	resources := make(chan *schema.Resource)
	go func() {
		defer close(resources)
		switch s.strategy {
		case StrategyDFS:
			s.syncDfs(ctx, resources)
		case StrategyRoundRobin:
			s.syncRoundRobin(ctx, resources)
		default:
			panic(fmt.Errorf("unknown scheduler %s", s.strategy))
		}
	}()
	for resource := range resources {
		vector := resource.GetValues()
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, resource.Table.ToArrowSchema())
		scalar.AppendToRecordBuilder(bldr, vector)
		rec := bldr.NewRecord()
		res <- &message.Insert{Record: rec}
	}
	return nil
}

func (s *Scheduler) logTablesMetrics(tables schema.Tables, client Client) {
	clientName := client.ID()
	for _, table := range tables {
		metrics := s.metrics.TableClient[table.Name][clientName]
		s.logger.Info().Str("table", table.Name).Str("client", clientName).Uint64("resources", metrics.Resources).Uint64("errors", metrics.Errors).Msg("table sync finished")
		s.logTablesMetrics(table.Relations, client)
	}
}

func (s *Scheduler) resolveResource(ctx context.Context, table *schema.Table, client schema.ClientMeta, parent *schema.Resource, item any) *schema.Resource {
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

func (s *Scheduler) resolveColumn(ctx context.Context, logger zerolog.Logger, tableMetrics *TableClientMetrics, client schema.ClientMeta, resource *schema.Resource, c schema.Column) {
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
		v := funk.Get(resource.GetItem(), s.caser.ToPascal(c.Name), funk.WithAllowZero())
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
func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
