package plugins

import (
	"context"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/internal/testdata"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
)

func TestDestinationPlugin(t *testing.T) {
	p := NewDestinationPlugin("test", "development", NewTestDestinationMemDBClient)
	DestinationPluginTestSuiteRunner(t, p, nil,
		DestinationTestSuiteTests{})
}

func TestDestinationOnNewError(t *testing.T) {
	ctx := context.Background()
	p := NewDestinationPlugin("test", "development", newTestDestinationMemDBClientErrOnNew)
	err := p.Init(ctx, getTestLogger(t), specs.Destination{})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDestinationOnWriteError(t *testing.T) {
	ctx := context.Background()
	newClientFunc := getNewTestDestinationMemDBClient(withErrOnWrite())
	p := NewDestinationPlugin("test", "development", newClientFunc)
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
	p.Write(ctx, tables, sourceName, syncTime, ch)
}

func TestDestinationOnWriteCtxCancelled(t *testing.T) {
	ctx := context.Background()
	newClientFunc := getNewTestDestinationMemDBClient(withBlockingWrite())
	p := NewDestinationPlugin("test", "development", newClientFunc)
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
	go func() {
		timer := time.NewTimer(4 * time.Second)
		<-timer.C
		// close the channel after we already write to it so we can test
		// context is cancelled successfully
		close(ch)
	}()
	p.Write(ctx, tables, sourceName, syncTime, ch)
}
