// Package inmemory is the default Storage backend — holds all scheduler
// state in process memory. Matches the behavior of the pre-existing
// ConcurrentRandomQueue: random-pop work semantics, atomic refcounts.
package inmemory

import (
	"context"
	"errors"
	"math/rand"
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/scheduler/storage"
)

type Storage struct {
	mu        sync.Mutex
	queue     []storage.SerializedWorkUnit
	resources map[string]*resourceEntry
	random    *rand.Rand
}

type resourceEntry struct {
	data     []byte
	refcount int
}

// New returns a Storage seeded for deterministic random-pop ordering.
func New(seed int64) *Storage {
	return &Storage{
		queue:     make([]storage.SerializedWorkUnit, 0),
		resources: make(map[string]*resourceEntry),
		random:    rand.New(rand.NewSource(seed)),
	}
}

func (s *Storage) PushWork(_ context.Context, w storage.SerializedWorkUnit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = append(s.queue, w)
	return nil
}

func (s *Storage) PushWorkBatch(_ context.Context, ws []storage.SerializedWorkUnit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = append(s.queue, ws...)
	return nil
}

func (s *Storage) PopWork(_ context.Context) (*storage.SerializedWorkUnit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.queue) == 0 {
		return nil, nil
	}
	idx := s.random.Intn(len(s.queue))
	last := len(s.queue) - 1
	s.queue[idx], s.queue[last] = s.queue[last], s.queue[idx]
	item := s.queue[last]
	s.queue = s.queue[:last]
	return &item, nil
}

func (s *Storage) WorkLen(_ context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.queue), nil
}

func (s *Storage) PutResource(_ context.Context, id string, data []byte, refcount int) error {
	if refcount < 1 {
		return errors.New("storage/inmemory: refcount must be >= 1")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]byte, len(data))
	copy(cp, data)
	s.resources[id] = &resourceEntry{data: cp, refcount: refcount}
	return nil
}

func (s *Storage) GetResource(_ context.Context, id string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.resources[id]
	if !ok {
		return nil, storage.ErrResourceNotFound
	}
	out := make([]byte, len(entry.data))
	copy(out, entry.data)
	return out, nil
}

func (s *Storage) DecResourceRefcount(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.resources[id]
	if !ok {
		return storage.ErrResourceNotFound
	}
	entry.refcount--
	if entry.refcount <= 0 {
		delete(s.resources, id)
	}
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = nil
	s.resources = nil
	return nil
}
