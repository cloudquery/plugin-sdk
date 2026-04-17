// Package storagetest provides a contract test suite that every Storage
// backend must pass. Each backend's test file calls TestContract with a
// factory for a fresh instance.
package storagetest

import (
	"context"
	"errors"
	"sync"
	"testing"

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
	t.Run("work_len", func(t *testing.T) { testWorkLen(t, newStorage(t)) })
	t.Run("resource_put_get", func(t *testing.T) { testResourcePutGet(t, newStorage(t)) })
	t.Run("resource_refcount_delete_on_zero", func(t *testing.T) { testRefcountDeleteOnZero(t, newStorage(t)) })
	t.Run("resource_get_missing_errors", func(t *testing.T) { testGetMissingErrors(t, newStorage(t)) })
	t.Run("resource_dec_missing_errors", func(t *testing.T) { testDecMissingErrors(t, newStorage(t)) })
	t.Run("concurrent_push_pop_no_loss", func(t *testing.T) { testConcurrentPushPopNoLoss(t, newStorage(t)) })
	t.Run("concurrent_refcount_no_double_delete", func(t *testing.T) { testConcurrentRefcountNoDoubleDelete(t, newStorage(t)) })
	t.Run("close_is_idempotent", func(t *testing.T) { testCloseIsIdempotent(t, newStorage(t)) })
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
	require.NoError(t, s.PutResource(ctx, "id-1", data, 1))

	got, err := s.GetResource(ctx, "id-1")
	require.NoError(t, err)
	require.Equal(t, data, got)
}

func testRefcountDeleteOnZero(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	require.NoError(t, s.PutResource(ctx, "id-1", []byte("x"), 2))

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
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			_ = s.PushWork(ctx, storage.SerializedWorkUnit{TableName: "t"})
		}
	}()

	popped := 0
	for popped < n {
		got, err := s.PopWork(ctx)
		require.NoError(t, err)
		if got != nil {
			popped++
		}
	}
	wg.Wait()

	n2, err := s.WorkLen(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, n2)
}

func testConcurrentRefcountNoDoubleDelete(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	const n = 100
	require.NoError(t, s.PutResource(ctx, "shared", []byte("x"), n))

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_ = s.DecResourceRefcount(ctx, "shared")
		}()
	}
	wg.Wait()

	_, err := s.GetResource(ctx, "shared")
	require.True(t, errors.Is(err, storage.ErrResourceNotFound), "resource should be deleted after all refs drained, got err=%v", err)
}

func testCloseIsIdempotent(t *testing.T, s storage.Storage) {
	ctx := context.Background()
	require.NoError(t, s.Close(ctx))
	require.NoError(t, s.Close(ctx))
}
