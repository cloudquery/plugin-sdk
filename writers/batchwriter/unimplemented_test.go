package batchwriter_test

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/writers/batchwriter"
)

type testDummyClient struct {
	batchwriter.IgnoreMigrateTables
	batchwriter.UnimplementedDeleteStale
	batchwriter.UnimplementedDeleteRecord
}

func (testDummyClient) WriteTableBatch(context.Context, string, message.WriteInserts) error {
	return nil
}

var _ batchwriter.Client = (*testDummyClient)(nil)
