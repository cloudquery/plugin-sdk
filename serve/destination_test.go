package serve

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/clients"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bufDestinationDialer(context.Context, string) (net.Conn, error) {
	return testDestinationListener.Dial()
}

type testDestinationClient struct {
}

func newDestinationClient(context.Context, zerolog.Logger, specs.Destination) (plugins.DestinationClient, error) {
	return &testDestinationClient{}, nil
}

func (*testDestinationClient) Initialize(context.Context, specs.Destination) error {
	return nil
}
func (*testDestinationClient) Migrate(context.Context, schema.Tables) error {
	return nil
}
func (*testDestinationClient) Write(_ context.Context, tables schema.Tables, resources <-chan *schema.DestinationResource) error {
	for _ = range resources {
	}
	return nil
}

func (*testDestinationClient) Stats() plugins.DestinationStats {
	return plugins.DestinationStats{}
}

func (*testDestinationClient) Close(context.Context) error {
	return nil
}
func (*testDestinationClient) DeleteStale(context.Context, schema.Tables, string, time.Time) error {
	return nil
}

func TestDestination(t *testing.T) {
	plugin := plugins.NewDestinationPlugin("testDestinationPlugin", "development", newDestinationClient)
	s := &destinationServe{
		plugin: plugin,
	}
	cmd := newCmdDestinationRoot(s)
	cmd.SetArgs([]string{"serve", "--network", "test"})
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	var serverErr error
	go func() {
		defer wg.Done()
		serverErr = cmd.ExecuteContext(ctx)
	}()
	defer func() {
		cancel()
		wg.Wait()
	}()

	// wait for the server to start
	for {
		if testDestinationListener != nil {
			break
		}
		t.Log("waiting for grpc server to start")
		time.Sleep(time.Millisecond * 200)
	}

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDestinationDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	c, err := clients.NewDestinationClient(ctx, specs.RegistryGrpc, "", "", clients.WithDestinationGrpcConn(conn))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.Terminate(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := c.Initialize(ctx, specs.Destination{}); err != nil {
		t.Fatal(err)
	}

	name, err := c.Name(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if name != "testDestinationPlugin" {
		t.Fatalf("expected name to be testDestinationPlugin but got %s", name)
	}

	version, err := c.Version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if version != "development" {
		t.Fatalf("expected version to be development but got %s", version)
	}

	tables := schema.Tables{testTable()}
	if err := c.Migrate(ctx, tables); err != nil {
		t.Fatal(err)
	}

	resource := schema.NewResourceData(testTable(), nil, nil)
	resource.Set("test_column", 5)
	destResource := resource.ToDestinationResource()
	b, err := json.Marshal(destResource)
	if err != nil {
		t.Fatal(err)
	}
	resources := make(chan []byte, 1)
	resources <- b
	close(resources)
	if err := c.Write(ctx, tables, "test", time.Now(), resources); err != nil {
		t.Fatal(err)
	}

	if err := c.DeleteStale(ctx, nil, "testSource", time.Now()); err != nil {
		t.Fatalf("failed to call DeleteStale: %v", err)
	}

	if err := c.Close(ctx); err != nil {
		t.Fatalf("failed to call Close: %v", err)
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
