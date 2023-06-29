package streamingbatchwriter

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

// UnimplementedMigrateTable is a dummy handler to consume WriteMigrateTable messages
type UnimplementedMigrateTable struct{}

func (UnimplementedMigrateTable) MigrateTable(_ context.Context, ch <-chan *message.WriteMigrateTable) error {
	// nolint:revive
	for range ch {
	}
	return nil
}

// UnimplementedDeleteStale is a dummy handler to consume DeleteStale messages
type UnimplementedDeleteStale struct{}

func (UnimplementedDeleteStale) DeleteStale(_ context.Context, ch <-chan *message.WriteDeleteStale) error {
	// nolint:revive
	for range ch {
	}
	return nil
}
