package nop

import "context"

func New() *Backend {
	return &Backend{}
}

// Backend can be used in cases where no backend is specified to avoid the need to check for nil
// pointers in all resolvers.
type Backend struct{}

func (*Backend) Set(_ context.Context, _, _, _ string) error {
	return nil
}

func (*Backend) Get(_ context.Context, _, _ string) (string, error) {
	return "", nil
}

func (*Backend) Close(_ context.Context) error {
	return nil
}
