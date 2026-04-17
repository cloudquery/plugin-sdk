// Package badger implements the Storage interface using a local embedded
// BadgerDB. Primary v1 backend for memory offload.
package badger

import (
	"context"
	"errors"
	"fmt"

	badgerdb "github.com/dgraph-io/badger/v4"

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

func (s *Storage) PushWork(ctx context.Context, w storage.SerializedWorkUnit) error {
	return errors.New("not implemented")
}
func (s *Storage) PushWorkBatch(ctx context.Context, ws []storage.SerializedWorkUnit) error {
	return errors.New("not implemented")
}
func (s *Storage) PopWork(ctx context.Context) (*storage.SerializedWorkUnit, error) {
	return nil, errors.New("not implemented")
}
func (s *Storage) WorkLen(ctx context.Context) (int, error) {
	return 0, errors.New("not implemented")
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
