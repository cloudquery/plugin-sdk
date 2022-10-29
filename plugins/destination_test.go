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

func (*testDestinationClient) Migrate(context.Context, schema.Tables) error {
	return nil
}

func (*testDestinationClient) Read(context.Context, schema.Tables, chan<- *schema.DestinationResource) error {
	//nolint:revive
	return nil
}

func (*testDestinationClient) Write(_ context.Context, _ schema.Tables, res <-chan *schema.DestinationResource) error {
	//nolint:revive
	for range res {
	}
	return nil
}

func (*testDestinationClient) Metrics() DestinationMetrics {
	return DestinationMetrics{}
}

func (*testDestinationClient) DeleteStale(context.Context, schema.Tables, string, time.Time) error {
	return nil
}

func (*testDestinationClient) Close(context.Context) error {
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
