package badger_test

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
	"github.com/stretchr/testify/require"
)

func TestBadger_PutGetResource(t *testing.T) {
	s := newBadger(t)
	ctx := context.Background()
	require.NoError(t, s.PutResource(ctx, "id-1", []byte("hello"), 1))

	got, err := s.GetResource(ctx, "id-1")
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), got)
}

func TestBadger_RefcountDeleteOnZero(t *testing.T) {
	s := newBadger(t)
	ctx := context.Background()
	require.NoError(t, s.PutResource(ctx, "id-1", []byte("x"), 2))

	require.NoError(t, s.DecResourceRefcount(ctx, "id-1"))
	_, err := s.GetResource(ctx, "id-1")
	require.NoError(t, err)

	require.NoError(t, s.DecResourceRefcount(ctx, "id-1"))
	_, err = s.GetResource(ctx, "id-1")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
}

func TestBadger_GetMissing(t *testing.T) {
	s := newBadger(t)
	ctx := context.Background()
	_, err := s.GetResource(ctx, "missing")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
}

func TestBadger_DecMissing(t *testing.T) {
	s := newBadger(t)
	ctx := context.Background()
	err := s.DecResourceRefcount(ctx, "missing")
	require.ErrorIs(t, err, storage.ErrResourceNotFound)
}
