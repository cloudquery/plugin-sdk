package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow/array"
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func TestPluginRoundRobin(t *testing.T) {
	ctx := context.Background()
	p := NewPlugin("test", "v0.0.0", NewMemDBClient, WithUnmanaged())
	testTable := schema.TestTable("test_table", schema.TestSourceOptions{})
	syncTime := time.Now().UTC()
	testRecords := schema.GenTestData(testTable, schema.GenTestDataOptions{
		SourceName: "test",
		SyncTime:   syncTime,
		MaxRows:    1,
	})
	spec := pbPlugin.Spec{
		Name:      "test",
		Path:      "cloudquery/test",
		Version:   "v1.0.0",
		Registry:  pbPlugin.Spec_REGISTRY_GITHUB,
		WriteSpec: &pbPlugin.WriteSpec{},
		SyncSpec:  &pbPlugin.SyncSpec{},
	}
	if err := p.Init(ctx, spec); err != nil {
		t.Fatal(err)
	}

	if err := p.Migrate(ctx, schema.Tables{testTable}); err != nil {
		t.Fatal(err)
	}
	if err := p.writeAll(ctx, spec, syncTime, testRecords); err != nil {
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
	records, err := p.syncAll(ctx, syncTime, *spec.SyncSpec)
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
	records, err = p.syncAll(ctx, syncTime, *spec.SyncSpec)
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
