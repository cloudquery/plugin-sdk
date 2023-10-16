package batchwriter

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

type IgnoreMigrateTables struct{}

func (IgnoreMigrateTables) MigrateTables(context.Context, message.WriteMigrateTables) error {
	return nil
}

type UnimplementedDeleteStale struct{}

func (UnimplementedDeleteStale) DeleteStale(context.Context, message.WriteDeleteStales) error {
	return fmt.Errorf("DeleteStale: %w", plugin.ErrNotImplemented)
}

type UnimplementedDeleteRecord struct{}

func (UnimplementedDeleteRecord) DeleteRecord(context.Context, message.WriteDeleteRecords) error {
	return fmt.Errorf("DeleteRecord: %w", plugin.ErrNotImplemented)
}
