// Package badger implements the Storage interface using a local embedded
// BadgerDB. Primary v1 backend for memory offload.
package badger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	badgerdb "github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
)

// Options configure the Badger backend.
type Options struct {
	// Path is the on-disk directory for the Badger database. Must be writable.
	// The caller is responsible for ensuring isolation (e.g. appending an
	// invocation ID) across concurrent syncs of the same plugin.
	Path string
}

// Storage is a BadgerDB-backed Storage.
type Storage struct {
	db *badgerdb.DB
}

// Open creates or opens a Badger database at opts.Path.
func Open(opts Options) (*Storage, error) {
	if opts.Path == "" {
		return nil, errors.New("badger: Options.Path is required")
	}
	bopts := badgerdb.DefaultOptions(opts.Path).WithLogger(nil)
	db, err := badgerdb.Open(bopts)
	if err != nil {
		return nil, fmt.Errorf("badger: open %q: %w", opts.Path, err)
	}
	return &Storage{db: db}, nil
}

// --- Storage methods below this line will be filled in subsequent tasks. ---

const (
	workPrefix = "w/"
)

// workKey generates a unique key for a queued work unit. UUID gives random
// scan order when iterating the prefix, satisfying the contract's "no
// particular pop order" requirement.
func workKey() []byte {
	return []byte(workPrefix + uuid.NewString())
}

func (s *Storage) PushWork(ctx context.Context, w storage.SerializedWorkUnit) error {
	data, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("badger: marshal work: %w", err)
	}
	return s.db.Update(func(txn *badgerdb.Txn) error {
		return txn.Set(workKey(), data)
	})
}

func (s *Storage) PushWorkBatch(ctx context.Context, ws []storage.SerializedWorkUnit) error {
	if len(ws) == 0 {
		return nil
	}
	wb := s.db.NewWriteBatch()
	defer wb.Cancel()
	for _, w := range ws {
		data, err := json.Marshal(w)
		if err != nil {
			return fmt.Errorf("badger: marshal work: %w", err)
		}
		if err := wb.Set(workKey(), data); err != nil {
			return fmt.Errorf("badger: write batch: %w", err)
		}
	}
	return wb.Flush()
}

func (s *Storage) PopWork(ctx context.Context) (*storage.SerializedWorkUnit, error) {
	var out *storage.SerializedWorkUnit
	err := s.db.Update(func(txn *badgerdb.Txn) error {
		it := txn.NewIterator(badgerdb.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(workPrefix)
		it.Seek(prefix)
		if !it.ValidForPrefix(prefix) {
			return nil // empty queue
		}
		item := it.Item()
		key := item.KeyCopy(nil)
		value, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("badger: read work value: %w", err)
		}
		var w storage.SerializedWorkUnit
		if err := json.Unmarshal(value, &w); err != nil {
			return fmt.Errorf("badger: unmarshal work: %w", err)
		}
		out = &w
		return txn.Delete(key)
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Storage) WorkLen(ctx context.Context) (int, error) {
	count := 0
	err := s.db.View(func(txn *badgerdb.Txn) error {
		opts := badgerdb.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(workPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			count++
		}
		return nil
	})
	return count, err
}

func (s *Storage) PutResource(ctx context.Context, id string, data []byte, refcount int) error {
	return errors.New("not implemented")
}
func (s *Storage) GetResource(ctx context.Context, id string) ([]byte, error) {
	return nil, errors.New("not implemented")
}
func (s *Storage) DecResourceRefcount(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
func (s *Storage) Close(ctx context.Context) error {
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
