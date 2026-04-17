package scheduler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStorageFromConfig_InMemoryDefault(t *testing.T) {
	s, err := NewStorageFromConfig(nil, 42, "inv-1")
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close(context.Background())
}

func TestNewStorageFromConfig_InMemoryExplicit(t *testing.T) {
	s, err := NewStorageFromConfig(&QueueConfig{Type: QueueTypeInMemory}, 42, "inv-1")
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close(context.Background())
}

func TestNewStorageFromConfig_BadgerOpens(t *testing.T) {
	dir := t.TempDir()
	s, err := NewStorageFromConfig(&QueueConfig{Type: QueueTypeBadger, Path: dir}, 42, "inv-1")
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close(context.Background())
}

func TestNewStorageFromConfig_Invalid(t *testing.T) {
	_, err := NewStorageFromConfig(&QueueConfig{Type: "nope"}, 42, "inv-1")
	require.Error(t, err)
}
