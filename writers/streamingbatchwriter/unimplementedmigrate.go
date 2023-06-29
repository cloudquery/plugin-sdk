package streamingbatchwriter

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

type UnimplementedMigrateTable struct{}

func (UnimplementedMigrateTable) MigrateTable(_ context.Context, ch <-chan *message.WriteMigrateTable) error {
	// nolint:revive
	for range ch {
	}
	return nil
}
