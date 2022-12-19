package memdb

import (
	"context"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/internal/testdata"
	"github.com/cloudquery/plugin-sdk/plugins/destination"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

func TestPluginUnmanagedClient(t *testing.T) {
	p := destination.NewPlugin("test", "development", NewClient)
	destination.PluginTestSuiteRunner(t, p, nil,
		destination.PluginTestSuiteTests{})
}

func TestPluginManagedClient(t *testing.T) {
	p := destination.NewPlugin("test", "development", NewClient, destination.WithManagerWriter())
	destination.PluginTestSuiteRunner(t, p, nil,
		destination.PluginTestSuiteTests{})
}

func TestPluginOnNewError(t *testing.T) {
	ctx := context.Background()
	p := destination.NewPlugin("test", "development", NewClientErrOnNew)
	err := p.Init(ctx, getTestLogger(t), specs.Destination{})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOnWriteError(t *testing.T) {
	ctx := context.Background()
	newClientFunc := GetNewClient(WithErrOnWrite())
	p := destination.NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, getTestLogger(t), specs.Destination{}); err != nil {
		t.Fatal(err)
	}
	table := testdata.TestTable("test")
	tables := []*schema.Table{
		table,
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	ch := make(chan schema.DestinationResource, 1)
	ch <- schema.DestinationResource{
		TableName: "test",
		Data:      testdata.GenTestData(table),
	}
	close(ch)
	err := p.Write(ctx, tables, sourceName, syncTime, ch)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "errOnWrite" {
		t.Fatalf("expected errOnWrite, got %s", err.Error())
	}
}

func TestOnWriteCtxCancelled(t *testing.T) {
	ctx := context.Background()
	newClientFunc := GetNewClient(WithBlockingWrite())
	p := destination.NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, getTestLogger(t), specs.Destination{}); err != nil {
		t.Fatal(err)
	}
	table := testdata.TestTable("test")
	tables := []*schema.Table{
		testdata.TestTable("test"),
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	ch := make(chan schema.DestinationResource, 1)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	ch <- schema.DestinationResource{
		TableName: "test",
		Data:      testdata.GenTestData(table),
	}
	defer cancel()
	err := p.Write(ctx, tables, sourceName, syncTime, ch)
	if err != nil {
		t.Fatal(err)
	}
}
