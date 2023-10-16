package streamingbatchwriter

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

// IgnoreMigrateTable is a dummy handler to consume WriteMigrateTable messages
type IgnoreMigrateTable struct{}

func (IgnoreMigrateTable) MigrateTable(_ context.Context, ch <-chan *message.WriteMigrateTable) error {
	// nolint:revive
	for range ch {
	}
	return nil
}

// UnimplementedDeleteStale is a dummy handler to consume and error on DeleteStale messages
type UnimplementedDeleteStale struct{}

func (UnimplementedDeleteStale) DeleteStale(_ context.Context, ch <-chan *message.WriteDeleteStale) error {
	// nolint:revive
	for range ch {
	}
	return fmt.Errorf("DeleteStale: %w", plugin.ErrNotImplemented)
}

type UnimplementedDeleteRecordsBatch struct{}

func (UnimplementedDeleteRecordsBatch) DeleteRecords(_ context.Context, ch <-chan *message.WriteDeleteRecord) error {
	return fmt.Errorf("DeleteRecordsBatch: %w", plugin.ErrNotImplemented)
}
