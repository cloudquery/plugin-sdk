package backend

import "context"

type Backend interface {
	// Set sets the value for the given table and client id.
	Set(ctx context.Context, table, clientID, value string) error
	// Get returns the value for the given table and client id.
	Get(ctx context.Context, table, clientID string) (string, error)
	// Close closes the backend.
	Close(ctx context.Context) error
}
