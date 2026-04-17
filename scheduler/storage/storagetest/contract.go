// Package storagetest provides a contract test suite that every Storage
// backend must pass. Each backend's test file calls TestContract with a
// factory for a fresh instance.
package storagetest

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
	"github.com/stretchr/testify/require"
)

// TestContract runs every contract assertion against the Storage returned by
// newStorage. The factory must return an empty, independent instance on each
// call (contract tests mutate state).
func TestContract(t *testing.T, newStorage func(t *testing.T) storage.Storage) {
	t.Helper()
	t.Run("push_pop_roundtrip", func(t *testing.T) { testPushPopRoundtrip(t, newStorage(t)) })
	t.Run("pop_empty_returns_nil", func(t *testing.T) { testPopEmptyReturnsNil(t, newStorage(t)) })
	t.Run("push_batch", func(t *testing.T) { testPushBatch(t, newStorage(t)) })
	t.Run("push_batch_empty", func(t *testing.T) { testPushBatchEmpty(t, newStorage(t)) })
	t.Run("work_len", func(t *testing.T) { testWorkLen(t, newStorage(t)) })
	t.Run("resource_put_get", func(t *testing.T) { testResourcePutGet(t, newStorage(t)) })
	t.Run("resource_put_rejects_zero_refcount", func(t *testing.T) { testPutResourceRejectsZeroRefcount(t, newStorage(t)) })
	t.Run("resource_refcount_delete_on_zero", func(t *testing.T) { testRefcountDeleteOnZero(t, newStorage(t)) })
	t.Run("resource_get_missing_errors", func(t *testing.T) { testGetMissingErrors(t, newStorage(t)) })
	t.Run("resource_dec_missing_errors", func(t *testing.T) { testDecMissingErrors(t, newStorage(t)) })
	t.Run("concurrent_push_pop_no_loss", func(t *testing.T) { testConcurrentPushPopNoLoss(t, newStorage(t)) })
	t.Run("concurrent_refcount_no_double_delete", func(t *testing.T) { testConcurrentRefcountNoDoubleDelete(t, newStorage(t)) })
	t.Run("close_is_idempotent", func(t *testing.T) { testCloseIsIdempotent(t, newStorage(t)) })
	t.Run("cascade_dec_deletes_ancestors", func(t *testing.T) { testCascadeDecDeletesAncestors(t, newStorage(t)) })
	t.Run("put_resource_rejects_unknown_parent", func(t *testing.T) { testPutResourceRejectsUnknownParent(t, newStorage(t)) })
}

func testPushPopRoundtrip(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	wu := storage.SerializedWorkUnit{TableName: "t1", ClientID: "c1", ParentID: "p1"}
	require.NoError(t, s.PushWork(ctx, wu))

	got, err := s.PopWork(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, wu, *got)

	// Second pop drains the queue.
	got, err = s.PopWork(ctx)
	require.NoError(t, err)
	require.Nil(t, got)
}

func testPopEmptyReturnsNil(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	got, err := s.PopWork(ctx)
	require.NoError(t, err)
	require.Nil(t, got)
}

func testPushBatch(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	batch := []storage.SerializedWorkUnit{
		{TableName: "a"}, {TableName: "b"}, {TableName: "c"},
	}
	require.NoError(t, s.PushWorkBatch(ctx, batch))

	seen := map[string]bool{}
	for i := 0; i < 3; i++ {
		got, err := s.PopWork(ctx)
		require.NoError(t, err)
		require.NotNil(t, got)
		seen[got.TableName] = true
	}
	require.Equal(t, map[string]bool{"a": true, "b": true, "c": true}, seen)
}

func testWorkLen(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	n, err := s.WorkLen(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, n)

	require.NoError(t, s.PushWork(ctx, storage.SerializedWorkUnit{TableName: "t"}))
	require.NoError(t, s.PushWork(ctx, storage.SerializedWorkUnit{TableName: "t"}))

	n, err = s.WorkLen(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, n)
}

func testResourcePutGet(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	data := []byte("hello")
	require.NoError(t, s.PutResource(ctx, "id-1", data, 1, ""))

	got, err := s.GetResource(ctx, "id-1")
	require.NoError(t, err)
	require.Equal(t, data, got)
}

func testRefcountDeleteOnZero(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	require.NoError(t, s.PutResource(ctx, "id-1", []byte("x"), 2, ""))

	// First dec: resource still exists.
	require.NoError(t, s.DecResourceRefcount(ctx, "id-1"))
	got, err := s.GetResource(ctx, "id-1")
	require.NoError(t, err)
	require.Equal(t, []byte("x"), got)

	// Second dec: resource deleted.
	require.NoError(t, s.DecResourceRefcount(ctx, "id-1"))
	_, err = s.GetResource(ctx, "id-1")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
}

func testGetMissingErrors(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	_, err := s.GetResource(ctx, "missing")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
}

func testDecMissingErrors(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	err := s.DecResourceRefcount(ctx, "missing")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
}

func testConcurrentPushPopNoLoss(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	const n = 500
	var wg sync.WaitGroup
	var pushErr error
	var pushErrMu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			if err := s.PushWork(ctx, storage.SerializedWorkUnit{TableName: "t"}); err != nil {
				pushErrMu.Lock()
				if pushErr == nil {
					pushErr = err
				}
				pushErrMu.Unlock()
				return
			}
		}
	}()

	deadline := time.Now().Add(10 * time.Second)
	popped := 0
	for popped < n {
		if time.Now().After(deadline) {
			t.Fatalf("lost work: popped %d of %d after deadline", popped, n)
		}
		got, err := s.PopWork(ctx)
		require.NoError(t, err)
		if got != nil {
			popped++
		}
	}
	wg.Wait()

	pushErrMu.Lock()
	require.NoError(t, pushErr, "push error during concurrent test")
	pushErrMu.Unlock()

	n2, err := s.WorkLen(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, n2)
}

func testConcurrentRefcountNoDoubleDelete(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	const n = 100
	require.NoError(t, s.PutResource(ctx, "shared", []byte("x"), n, ""))

	var wg sync.WaitGroup
	errs := make(chan error, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := s.DecResourceRefcount(ctx, "shared"); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)

	// A correct backend gives us exactly n successful decs. Any error here
	// indicates a double-delete bug (over-decrement) or worse.
	var got []error
	for e := range errs {
		got = append(got, e)
	}
	require.Empty(t, got, "expected no errors from %d concurrent decs, got: %v", n, got)

	_, err := s.GetResource(ctx, "shared")
	require.ErrorIs(t, err, storage.ErrResourceNotFound, "resource should be deleted after all refs drained")
}

func testCloseIsIdempotent(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	require.NoError(t, s.Close(ctx))
	require.NoError(t, s.Close(ctx))
}

func testPutResourceRejectsZeroRefcount(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	err := s.PutResource(ctx, "id", []byte("x"), 0, "")
	require.Error(t, err, "PutResource with refcount=0 must return an error")

	err = s.PutResource(ctx, "id", []byte("x"), -1, "")
	require.Error(t, err, "PutResource with refcount<0 must return an error")

	_, err = s.GetResource(ctx, "id")
	require.ErrorIs(t, err, storage.ErrResourceNotFound, "failed Put must not create the resource")
}

func testPushBatchEmpty(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	require.NoError(t, s.PushWorkBatch(ctx, nil), "empty batch should be a no-op")
	require.NoError(t, s.PushWorkBatch(ctx, []storage.SerializedWorkUnit{}), "empty batch should be a no-op")

	n, err := s.WorkLen(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, n)
}

func testCascadeDecDeletesAncestors(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	// Chain: grandparent (refcount 1) ← parent (refcount 1) ← child (refcount 2)
	require.NoError(t, s.PutResource(ctx, "gp", []byte("g"), 1, ""))
	require.NoError(t, s.PutResource(ctx, "p", []byte("p"), 1, "gp"))
	// After the above, gp.refcount should be 2 (its own 1 + 1 for p).
	require.NoError(t, s.PutResource(ctx, "c", []byte("c"), 2, "p"))
	// After above, p.refcount should be 2 (its own 1 + 1 for c).

	// Drain child refcount → delete c → cascade dec p → now p.refcount = 1 (still there).
	require.NoError(t, s.DecResourceRefcount(ctx, "c"))
	require.NoError(t, s.DecResourceRefcount(ctx, "c"))
	_, err := s.GetResource(ctx, "c")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
	_, err = s.GetResource(ctx, "p")
	require.NoError(t, err, "p should still exist with refcount 1 after c cascade")
	_, err = s.GetResource(ctx, "gp")
	require.NoError(t, err, "gp should still exist with refcount 2 after p cascade")

	// Drain parent → delete p → cascade dec gp → now gp.refcount = 1 (still there because of its own initial 1).
	require.NoError(t, s.DecResourceRefcount(ctx, "p"))
	_, err = s.GetResource(ctx, "p")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
	_, err = s.GetResource(ctx, "gp")
	require.NoError(t, err, "gp should still exist after p cascade: its own refcount of 1 keeps it alive")

	// Drain gp's own refcount → delete gp, no further cascade (parentID=="").
	require.NoError(t, s.DecResourceRefcount(ctx, "gp"))
	_, err = s.GetResource(ctx, "gp")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
}

func testPutResourceRejectsUnknownParent(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	err := s.PutResource(ctx, "child", []byte("x"), 1, "does-not-exist")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
	_, err = s.GetResource(ctx, "child")
	require.ErrorIs(t, err, storage.ErrResourceNotFound, "child must not exist after failed Put")
}
