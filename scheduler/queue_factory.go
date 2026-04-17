package scheduler

import (
	"path/filepath"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
	badgerstore "github.com/cloudquery/plugin-sdk/v4/scheduler/storage/badger"
	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage/inmemory"
)

// NewStorageFromConfig constructs a Storage backend from the user-facing
// QueueConfig. The seed is used by the in-memory backend for deterministic
// random-pop ordering. invocationID is appended to disk-backed paths so
// concurrent syncs of the same plugin don't collide on the same directory.
//
// A nil cfg is treated as in-memory (the default).
func NewStorageFromConfig(cfg *QueueConfig, seed int64, invocationID string) (storage.Storage, error) {
	if cfg == nil {
		return inmemory.New(seed), nil
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	switch cfg.Type {
	case "", QueueTypeInMemory:
		return inmemory.New(seed), nil
	case QueueTypeBadger:
		path := filepath.Join(cfg.Path, invocationID)
		return badgerstore.Open(badgerstore.Options{Path: path})
	default:
		// unreachable given Validate above, but keeps the switch exhaustive
		return inmemory.New(seed), nil
	}
}
