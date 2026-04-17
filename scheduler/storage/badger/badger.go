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

// maxUpdateRetries bounds the retry loop for Badger Update calls that hit
// SSI transaction conflicts. Under contention from many concurrent workers,
// a handful of retries is typical; exceeding this limit indicates pathology
// (e.g., every worker racing on the same key at once).
const maxUpdateRetries = 100

// updateWithRetry runs fn inside s.db.Update, retrying on ErrConflict up to
// maxUpdateRetries times. Necessary because Badger's SSI Update does not
// auto-retry; callers must when they intend read-modify-write semantics.
func (s *Storage) updateWithRetry(fn func(txn *badgerdb.Txn) error) error {
	for i := 0; i < maxUpdateRetries; i++ {
		err := s.db.Update(fn)
		if !errors.Is(err, badgerdb.ErrConflict) {
			return err
		}
	}
	return fmt.Errorf("badger: transaction conflict after %d retries", maxUpdateRetries)
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
	workPrefix     = "w/"
	resourcePrefix = "r/"
)

// workKey generates a unique key for a queued work unit. UUID gives random
// scan order when iterating the prefix, satisfying the contract's "no
// particular pop order" requirement.
func workKey() []byte {
	return []byte(workPrefix + uuid.NewString())
}

// resourceEntry is the on-disk representation of a parent resource. Stored
// as a single JSON blob; refcount is a field so DecResourceRefcount can
// update it atomically within one transaction.
type resourceEntry struct {
	Data     []byte `json:"data"`
	Refcount int    `json:"refcount"`
	ParentID string `json:"parent_id,omitempty"`
}

func resourceKey(id string) []byte {
	return []byte(resourcePrefix + id)
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
	err := s.updateWithRetry(func(txn *badgerdb.Txn) error {
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

func (s *Storage) PutResource(ctx context.Context, id string, data []byte, refcount int, parentID string) error {
	if refcount < 1 {
		return errors.New("badger: refcount must be >= 1")
	}
	return s.updateWithRetry(func(txn *badgerdb.Txn) error {
		if parentID != "" {
			parentKey := resourceKey(parentID)
			parentItem, err := txn.Get(parentKey)
			if errors.Is(err, badgerdb.ErrKeyNotFound) {
				return storage.ErrResourceNotFound
			}
			if err != nil {
				return err
			}
			var parentEntry resourceEntry
			if err := parentItem.Value(func(val []byte) error {
				return json.Unmarshal(val, &parentEntry)
			}); err != nil {
				return fmt.Errorf("badger: unmarshal parent: %w", err)
			}
			parentEntry.Refcount++
			parentBlob, err := json.Marshal(parentEntry)
			if err != nil {
				return fmt.Errorf("badger: marshal parent: %w", err)
			}
			if err := txn.Set(parentKey, parentBlob); err != nil {
				return err
			}
		}
		entry := resourceEntry{Data: data, Refcount: refcount, ParentID: parentID}
		blob, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("badger: marshal resource: %w", err)
		}
		return txn.Set(resourceKey(id), blob)
	})
}

func (s *Storage) GetResource(ctx context.Context, id string) ([]byte, error) {
	var out []byte
	err := s.db.View(func(txn *badgerdb.Txn) error {
		item, err := txn.Get(resourceKey(id))
		if errors.Is(err, badgerdb.ErrKeyNotFound) {
			return storage.ErrResourceNotFound
		}
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			var entry resourceEntry
			if err := json.Unmarshal(val, &entry); err != nil {
				return fmt.Errorf("badger: unmarshal resource: %w", err)
			}
			out = append(out[:0], entry.Data...)
			return nil
		})
	})
	return out, err
}

func (s *Storage) DecResourceRefcount(ctx context.Context, id string) error {
	return s.updateWithRetry(func(txn *badgerdb.Txn) error {
		return decInTxn(txn, id)
	})
}

func decInTxn(txn *badgerdb.Txn, id string) error {
	key := resourceKey(id)
	item, err := txn.Get(key)
	if errors.Is(err, badgerdb.ErrKeyNotFound) {
		return storage.ErrResourceNotFound
	}
	if err != nil {
		return err
	}
	var entry resourceEntry
	if err := item.Value(func(val []byte) error {
		return json.Unmarshal(val, &entry)
	}); err != nil {
		return fmt.Errorf("badger: unmarshal resource: %w", err)
	}
	entry.Refcount--
	if entry.Refcount <= 0 {
		if err := txn.Delete(key); err != nil {
			return err
		}
		if entry.ParentID != "" {
			if err := decInTxn(txn, entry.ParentID); err != nil && !errors.Is(err, storage.ErrResourceNotFound) {
				return err
			}
		}
		return nil
	}
	blob, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("badger: marshal resource: %w", err)
	}
	return txn.Set(key, blob)
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
