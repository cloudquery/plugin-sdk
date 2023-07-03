package mixedbatchwriter

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

type IgnoreMigrateTableBatch struct{}

func (IgnoreMigrateTableBatch) MigrateTableBatch(context.Context, message.WriteMigrateTables) error {
	return nil
}

type UnimplementedDeleteStaleBatch struct{}

func (UnimplementedDeleteStaleBatch) DeleteStaleBatch(context.Context, message.WriteDeleteStales) error {
	return fmt.Errorf("DeleteStaleBatch: %w", plugin.ErrNotImplemented)
}
