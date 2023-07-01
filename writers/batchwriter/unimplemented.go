package batchwriter

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

type UnimplementedMigrateTables struct{}

func (UnimplementedMigrateTables) MigrateTables(context.Context, message.WriteMigrateTables) error {
	return nil
}

type UnimplementedDeleteStale struct{}

func (UnimplementedDeleteStale) DeleteStale(context.Context, message.WriteDeleteStales) error {
	return nil
}
