package badger_test

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
	bstore "github.com/cloudquery/plugin-sdk/v4/scheduler/storage/badger"
	"github.com/stretchr/testify/require"
)

func newBadger(t *testing.T) *bstore.Storage {
	t.Helper()
	dir := t.TempDir()
	s, err := bstore.Open(bstore.Options{Path: dir})
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close(context.Background()) })
	return s
}

func TestBadger_PushPopWorkRoundtrip(t *testing.T) {
	s := newBadger(t)
	ctx := context.Background()
	want := storage.SerializedWorkUnit{TableName: "t1", ClientID: "c1", ParentID: "p1"}
	require.NoError(t, s.PushWork(ctx, want))

	got, err := s.PopWork(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, want, *got)

	got, err = s.PopWork(ctx)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestBadger_WorkLen(t *testing.T) {
	s := newBadger(t)
	ctx := context.Background()
	n, err := s.WorkLen(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, n)

	require.NoError(t, s.PushWork(ctx, storage.SerializedWorkUnit{TableName: "t"}))
	require.NoError(t, s.PushWorkBatch(ctx, []storage.SerializedWorkUnit{{TableName: "t"}, {TableName: "t"}}))

	n, err = s.WorkLen(ctx)
	require.NoError(t, err)
	require.Equal(t, 3, n)
}
