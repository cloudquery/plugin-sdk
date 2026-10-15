package queue

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

func workUnit(name string) WorkUnit {
	return WorkUnit{Table: &schema.Table{Name: name}}
}

// startNext calls Next in a goroutine and returns after Next blocks or returns.
func startNext(ctx context.Context, q *workQueue) <-chan *WorkUnit {
	next := make(chan *WorkUnit, 1)
	go func() { next <- q.Next(ctx) }()
	synctest.Wait()
	return next
}

func requireNextReturned(t *testing.T, next <-chan *WorkUnit) *WorkUnit {
	t.Helper()
	synctest.Wait()
	select {
	case item := <-next:
		return item
	default:
		require.FailNow(t, "Next is still blocked")
		return nil
	}
}

func TestWorkQueue(t *testing.T) {
	t.Run("returns nil when nothing is queued or in progress", func(t *testing.T) {
		require.Nil(t, newWorkQueue(0, 0).Next(t.Context()))
	})

	t.Run("wakes up when a work unit in progress queues a relation", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			q := newWorkQueue(0, 1)
			q.Push(workUnit("parent"))
			require.Equal(t, "parent", q.Next(t.Context()).Table.Name)

			next := startNext(t.Context(), q)
			q.Push(workUnit("relation"))

			require.Equal(t, "relation", requireNextReturned(t, next).Table.Name)
		})
	})

	t.Run("returns a relation that the last work unit in progress queued before it finished", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			q := newWorkQueue(0, 1)
			q.Push(workUnit("parent"))
			require.Equal(t, "parent", q.Next(t.Context()).Table.Name)

			next := startNext(t.Context(), q)
			q.Push(workUnit("relation"))
			q.Done()

			require.Equal(t, "relation", requireNextReturned(t, next).Table.Name)
		})
	})

	t.Run("returns nil when the last work unit in progress finishes", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			q := newWorkQueue(0, 1)
			q.Push(workUnit("parent"))
			require.Equal(t, "parent", q.Next(t.Context()).Table.Name)

			next := startNext(t.Context(), q)
			q.Done()

			require.Nil(t, requireNextReturned(t, next))
		})
	})

	t.Run("returns nil when the context is done", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			q := newWorkQueue(0, 1)
			q.Push(workUnit("parent"))
			require.Equal(t, "parent", q.Next(ctx).Table.Name)

			next := startNext(ctx, q)
			cancel()

			require.Nil(t, requireNextReturned(t, next))
		})
	})
}
