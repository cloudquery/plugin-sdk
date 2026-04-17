// Package storage defines the contract for backends that hold scheduler
// work state (work units and parent resources) during a sync.
//
// Backends store opaque bytes — the scheduler owns all serialization.
// See docs/superpowers/specs/2026-04-17-external-queue-scheduler-design.md
// for the full design.
package storage

import (
	"context"
	"errors"
)

// ErrResourceNotFound is returned by GetResource / DecResourceRefcount when
// the ID is absent. Callers should treat this as a programming error
// (indicates a leaked reference or double-free) and fail the sync.
var ErrResourceNotFound = errors.New("resource not found")

// SerializedWorkUnit is the on-the-wire representation of a scheduled unit of
// work. It holds references only — the concrete Table/Client/Parent are
// reconstituted in-process by the scheduler.
type SerializedWorkUnit struct {
	TableName string // lookup key into the plugin's registered tables
	ClientID  string // lookup key into the plugin's initialized clients
	ParentID  string // empty if top-level; else ID in resource KV
}

// Storage is the pluggable backend for scheduler work state.
type Storage interface {
	PushWork(ctx context.Context, w SerializedWorkUnit) error
	PushWorkBatch(ctx context.Context, ws []SerializedWorkUnit) error
	// PopWork removes and returns a work unit. Returns (nil, nil) when empty;
	// returns an error only on backend failure. Pop semantics are defined by
	// the backend (random for in-memory, FIFO-ish for badger — callers must
	// not assume an ordering beyond "work eventually drains").
	PopWork(ctx context.Context) (*SerializedWorkUnit, error)
	WorkLen(ctx context.Context) (int, error)

	// PutResource inserts a resource blob with an initial refcount.
	// refcount must be >= 1 (a resource with zero pins should never exist).
	PutResource(ctx context.Context, id string, data []byte, refcount int) error
	GetResource(ctx context.Context, id string) ([]byte, error)
	// DecResourceRefcount decrements refcount by 1 and deletes when it
	// reaches zero, atomically within a single backend operation.
	DecResourceRefcount(ctx context.Context, id string) error

	Close(ctx context.Context) error
}
