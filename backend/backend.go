package backend

type Backend interface {
	// Set sets the value for the given key.
	Set(table, key, value string) error
	// Get returns the value for the given key.
	Get(table, key string) (string, error)
	// Close closes the backend.
	Close() error
}
