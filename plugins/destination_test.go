package plugins

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

type testDestinationClient struct {

}

func newTestDestinationClient(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
	return &testDestinationClient{}, nil
}

func newErrorTestDestinationClient(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
	return nil, fmt.Errorf("failed to create test destination client")
}

func (c *testDestinationClient) Migrate(ctx context.Context, tables schema.Tables) error {
	return nil
}

func (c *testDestinationClient) Write(ctx context.Context, tables schema.Tables, res <-chan *schema.DestinationResource) error {
	for _ = range res {
	}
	return nil
}

func (c *testDestinationClient) Stats() DestinationStats {
	return DestinationStats{}
}

func (c *testDestinationClient) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	return nil	
}

func (c *testDestinationClient) Close(ctx context.Context) error {
	return nil
}


func TestDestinationPlugin(t *testing.T) {
	ctx := context.Background()
	d := NewDestinationPlugin("testDestinationPlugin", "development", newErrorTestDestinationClient)
	if err := d.Init(ctx, zerolog.New(zerolog.NewTestWriter(t)), specs.Destination{}); err == nil {
		t.Fatal("expected error got nil")
	}

	d = NewDestinationPlugin("test", "development", newTestDestinationClient)
	if d.Name() != "test" {
		t.Errorf("expected name to be test but got %s", d.Name())
	}
	if d.Version() != "development" {
		t.Errorf("expected version to be development but got %s", d.Version())
	}

	if err := d.Init(ctx, zerolog.Nop(), specs.Destination{}); err != nil {
		t.Fatal(err)
	}
	
}