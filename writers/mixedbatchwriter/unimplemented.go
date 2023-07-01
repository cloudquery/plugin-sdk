package mixedbatchwriter

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

type UnimplementedMigrateTableBatch struct{}

func (UnimplementedMigrateTableBatch) MigrateTableBatch(context.Context, message.WriteMigrateTables) error {
	return nil
}

type UnimplementedDeleteStaleBatch struct{}

func (UnimplementedDeleteStaleBatch) DeleteStaleBatch(context.Context, message.WriteDeleteStales) error {
	return nil
}
