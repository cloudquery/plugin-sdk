package nop

import "context"

func New() *Backend {
	return &Backend{}
}

// Backend can be used in cases where no backend is specified to avoid the need to check for nil
// pointers in all resolvers.
type Backend struct{}

func (*Backend) Set(ctx context.Context, table, clientID, value string) error {
	return nil
}

func (*Backend) Get(ctx context.Context, table, clientID string) (string, error) {
	return "", nil
}

func (*Backend) Close(ctx context.Context) error {
	return nil
}
