package destination

import (
	"context"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/internal/testdata"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

func TestPlugin(t *testing.T) {
	p := NewPlugin("test", "development", NewTestDestinationMemDBClient)
	PluginTestSuiteRunner(t, p, nil,
		TestSuiteTests{})
}

func TestOnNewError(t *testing.T) {
	ctx := context.Background()
	p := NewPlugin("test", "development", newTestDestinationMemDBClientErrOnNew)
	err := p.Init(ctx, getTestLogger(t), specs.Destination{})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOnWriteError(t *testing.T) {
	ctx := context.Background()
	newClientFunc := getNewTestDestinationMemDBClient(withErrOnWrite())
	p := NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, getTestLogger(t), specs.Destination{}); err != nil {
		t.Fatal(err)
	}
	tables := []*schema.Table{
		testdata.TestTable("test"),
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	ch := make(chan schema.DestinationResource, 1)
	ch <- schema.DestinationResource{
		TableName: "test",
		Data:      testdata.TestData(),
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
	newClientFunc := getNewTestDestinationMemDBClient(withBlockingWrite())
	p := NewPlugin("test", "development", newClientFunc)
	if err := p.Init(ctx, getTestLogger(t), specs.Destination{}); err != nil {
		t.Fatal(err)
	}
	tables := []*schema.Table{
		testdata.TestTable("test"),
	}
	sourceName := "TestDestinationOnWriteError"
	syncTime := time.Now()
	ch := make(chan schema.DestinationResource, 1)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	ch <- schema.DestinationResource{
		TableName: "test",
		Data:      testdata.TestData(),
	}
	defer cancel()
	err := p.Write(ctx, tables, sourceName, syncTime, ch)
	if err != nil {
		t.Fatal(err)
	}
}
