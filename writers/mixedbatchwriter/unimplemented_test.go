package mixedbatchwriter_test

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/writers/mixedbatchwriter"
)

type testDummyClient struct {
	mixedbatchwriter.IgnoreMigrateTableBatch
	mixedbatchwriter.UnimplementedDeleteStaleBatch
}

func (testDummyClient) InsertBatch(context.Context, message.WriteInserts) error {
	return nil
}

var _ mixedbatchwriter.Client = (*testDummyClient)(nil)
