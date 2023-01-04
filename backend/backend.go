package backend

import "context"

type Backend interface {
	// Set sets the value for the given key.
	Set(ctx context.Context, table, key, value string) error
	// Get returns the value for the given key.
	Get(ctx context.Context, table, key string) (string, error)
	// Close closes the backend.
	Close(ctx context.Context) error
}
