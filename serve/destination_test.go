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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bufDestinationDialer(context.Context, string) (net.Conn, error) {
	return testDestinationListener.Dial()
}

func TestDestination(t *testing.T) {
	plugin := plugins.NewDestinationPlugin("testDestinationPlugin", "development", plugins.NewTestDestinationMemDBClient)
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
	c, err := clients.NewDestinationClient(ctx, specs.RegistryGrpc, "", "", clients.WithDestinationGrpcConn(conn), clients.WithDestinationNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.Terminate(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := c.Initialize(ctx, specs.Destination{
		WriteMode: specs.WriteModeAppend,
	}); err != nil {
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
	if err := resource.Set("test_column", 5); err != nil {
		t.Fatal(err)
	}
	destResource := resource.ToDestinationResource()
	b, err := json.Marshal(destResource)
	if err != nil {
		t.Fatal(err)
	}
	resources := make(chan []byte, 1)
	resources <- b
	close(resources)
	if err := c.Write2(ctx, tables, "test", time.Now(), resources); err != nil {
		t.Fatal(err)
	}

	if err := c.DeleteStale(ctx, nil, "testSource", time.Now()); err != nil {
		t.Fatalf("failed to call DeleteStale: %v", err)
	}

	_, err = c.GetMetrics(ctx)
	if err != nil {
		t.Fatal(err)
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
