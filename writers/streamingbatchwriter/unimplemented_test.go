package streamingbatchwriter_test

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/writers/streamingbatchwriter"
)

type testDummyClient struct {
	streamingbatchwriter.IgnoreMigrateTable
	streamingbatchwriter.UnimplementedDeleteStale
	streamingbatchwriter.UnimplementedDeleteRecordsBatch
}

func (testDummyClient) WriteTable(context.Context, <-chan *message.WriteInsert) error {
	return nil
}

var _ streamingbatchwriter.Client = (*testDummyClient)(nil)
