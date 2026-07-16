package scheduler

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var errTableResolverBoom = errors.New("table resolver boom")

// TestSchedulerErrorClassifier verifies that a table resolver error is raised as a
// SyncError message by default, and suppressed (no SyncError) when the configured
// ErrorClassifier returns true for it.
func TestSchedulerErrorClassifier(t *testing.T) {
	for _, strategy := range AllStrategies {
		t.Run(strategy.String(), func(t *testing.T) {
			table := &schema.Table{
				Name: "test_table",
				Resolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, _ chan<- any) error {
					return errTableResolverBoom
				},
				Columns: []schema.Column{
					{Name: "name", Type: arrow.BinaryTypes.String},
				},
			}

			t.Run("raised by default", func(t *testing.T) {
				errs := syncErrors(t, strategy, table, nil)
				require.Len(t, errs, 1, "table resolver error should be raised as a SyncError")
				require.Equal(t, "test_table", errs[0].TableName)
				require.Contains(t, errs[0].Error, errTableResolverBoom.Error())
			})

			t.Run("suppressed by classifier", func(t *testing.T) {
				var seen []schema.ErrorEvent
				var mu sync.Mutex
				classifier := func(_ context.Context, err error, event schema.ErrorEvent) bool {
					mu.Lock()
					seen = append(seen, event)
					mu.Unlock()
					return errors.Is(err, errTableResolverBoom)
				}
				errs := syncErrors(t, strategy, table, classifier)
				require.Empty(t, errs, "suppressed table resolver error should not emit a SyncError")

				mu.Lock()
				defer mu.Unlock()
				require.Len(t, seen, 1)
				require.Equal(t, schema.ErrorPhaseTableResolver, seen[0].Phase)
				require.Equal(t, table, seen[0].Table)
				require.NotNil(t, seen[0].Client)
			})

			t.Run("not suppressed when classifier returns false", func(t *testing.T) {
				classifier := func(_ context.Context, _ error, _ schema.ErrorEvent) bool {
					return false
				}
				errs := syncErrors(t, strategy, table, classifier)
				require.Len(t, errs, 1, "a classifier returning false must not suppress the error")
			})
		})
	}
}

func syncErrors(t *testing.T, strategy Strategy, table *schema.Table, classifier schema.ErrorClassifier) []*message.SyncError {
	t.Helper()
	opts := []Option{
		WithLogger(zerolog.New(zerolog.NewTestWriter(t)).Level(zerolog.DebugLevel)),
		WithStrategy(strategy),
	}
	if classifier != nil {
		opts = append(opts, WithErrorClassifier(classifier))
	}
	sc := NewScheduler(opts...)

	msgs := make(chan message.SyncMessage, 100)
	require.NoError(t, sc.Sync(context.Background(), &testExecutionClient{}, schema.Tables{table}, msgs))
	close(msgs)

	var syncErrs []*message.SyncError
	for msg := range msgs {
		if e, ok := msg.(*message.SyncError); ok {
			syncErrs = append(syncErrs, e)
		}
	}
	return syncErrs
}
