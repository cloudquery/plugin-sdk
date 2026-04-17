package scheduler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueueConfig_ValidateInMemoryNoPath(t *testing.T) {
	cfg := &QueueConfig{Type: QueueTypeInMemory}
	require.NoError(t, cfg.Validate())
}

func TestQueueConfig_ValidateBadgerRequiresPath(t *testing.T) {
	cfg := &QueueConfig{Type: QueueTypeBadger}
	err := cfg.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "path")
}

func TestQueueConfig_ValidateUnknownType(t *testing.T) {
	cfg := &QueueConfig{Type: "redis"}
	err := cfg.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "redis")
	require.Contains(t, err.Error(), "in-memory")
	require.Contains(t, err.Error(), "badger")
}

func TestQueueConfig_ValidateBadgerOK(t *testing.T) {
	cfg := &QueueConfig{Type: QueueTypeBadger, Path: "/tmp/q"}
	require.NoError(t, cfg.Validate())
}

func TestQueueConfig_RequiresShuffleQueueStrategy(t *testing.T) {
	cfg := &QueueConfig{Type: QueueTypeBadger, Path: "/tmp/q"}
	require.NoError(t, cfg.ValidateWithStrategy(StrategyShuffleQueue))

	err := cfg.ValidateWithStrategy(StrategyDFS)
	require.Error(t, err)
	require.Contains(t, err.Error(), "shuffle-queue")
}

func TestQueueConfig_InMemoryAllowsAnyStrategy(t *testing.T) {
	cfg := &QueueConfig{Type: QueueTypeInMemory}
	require.NoError(t, cfg.ValidateWithStrategy(StrategyDFS))
}
