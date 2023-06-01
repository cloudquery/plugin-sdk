package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func TestPluginUnmanagedSync(t *testing.T) {
	ctx := context.Background()
	p := NewPlugin("test", "v0.0.0", NewMemDBClient, WithUnmanagedSync())
	testTable := schema.TestTable("test_table", schema.TestSourceOptions{})
	syncTime := time.Now().UTC()
	sourceName := "test"
	testRecords := schema.GenTestData(testTable, schema.GenTestDataOptions{
		SourceName: sourceName,
		SyncTime:   syncTime,
		MaxRows:    1,
	})
	if err := p.Init(ctx, nil); err != nil {
		t.Fatal(err)
	}

	if err := p.Migrate(ctx, schema.Tables{testTable}, MigrateModeSafe); err != nil {
		t.Fatal(err)
	}
	if err := p.writeAll(ctx, sourceName, syncTime, WriteModeOverwrite, testRecords); err != nil {
		t.Fatal(err)
	}
	gotRecords, err := p.readAll(ctx, testTable, "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(gotRecords) != len(testRecords) {
		t.Fatalf("got %d records, want %d", len(gotRecords), len(testRecords))
	}
	if !array.RecordEqual(testRecords[0], gotRecords[0]) {
		t.Fatal("records are not equal")
	}
	records, err := p.syncAll(ctx, sourceName, syncTime, SyncOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("got %d resources, want 1", len(records))
	}

	if !array.RecordEqual(testRecords[0], records[0]) {
		t.Fatal("records are not equal")
	}

	newSyncTime := time.Now().UTC()
	if err := p.DeleteStale(ctx, schema.Tables{testTable}, "test", newSyncTime); err != nil {
		t.Fatal(err)
	}
	records, err = p.syncAll(ctx, sourceName, syncTime, SyncOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 0 {
		t.Fatalf("got %d resources, want 0", len(records))
	}

	if err := p.Close(ctx); err != nil {
		t.Fatal(err)
	}
}

// func TestPluginInit(t *testing.T) {
// 	const (
// 		batchSize      = uint64(100)
// 		batchSizeBytes = uint64(1000)
// 	)

// 	var (
// 		batchSizeObserved      uint64
// 		batchSizeBytesObserved uint64
// 	)
// 	p := NewPlugin(
// 		"test",
// 		"development",
// 		func(ctx context.Context, logger zerolog.Logger, s any) (Client, error) {
// 			batchSizeObserved = s.WriteSpec.BatchSize
// 			batchSizeBytesObserved = s.WriteSpec.BatchSizeBytes
// 			return NewMemDBClient(ctx, logger, s)
// 		},
// 		WithDefaultBatchSize(int(batchSize)),
// 		WithDefaultBatchSizeBytes(int(batchSizeBytes)),
// 	)
// 	require.NoError(t, p.Init(context.TODO(), nil))

// 	require.Equal(t, batchSize, batchSizeObserved)
// 	require.Equal(t, batchSizeBytes, batchSizeBytesObserved)
// }
