package serve

import (
	"context"
	"sync"
	"testing"

	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v4"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestStateV4(t *testing.T) {
	p := plugin.NewPlugin(
		"testPluginV3",
		"v1.0.0",
		memdb.NewMemDBClient)
	srv := Plugin(p, WithArgs("serve"), WithTestListener())
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	var serverErr error
	go func() {
		defer wg.Done()
		serverErr = srv.Serve(ctx)
	}()
	defer func() {
		cancel()
		wg.Wait()
	}()

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(srv.bufPluginDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	c := pb.NewPluginClient(conn)
	if _, err := c.Init(ctx, &pb.Init_Request{}); err != nil {
		t.Fatal(err)
	}
	stateClient, err := state.NewClient(ctx, c, "test")
	if err != nil {
		t.Fatal(err)
	}

	if err := stateClient.SetKey(ctx, "key", "value"); err != nil {
		t.Fatal(err)
	}

	val, err := stateClient.GetKey(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}
	if val != "value" {
		t.Fatalf("expected value to be value but got %s", val)
	}

	if err := stateClient.Flush(ctx); err != nil {
		t.Fatal(err)
	}
	stateClient, err = state.NewClient(ctx, c, "test")
	if err != nil {
		t.Fatal(err)
	}
	val, err = stateClient.GetKey(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}
	if val != "value" {
		t.Fatalf("expected value to be value but got %s", val)
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}

func TestStateV4OverwriteGetLatest(t *testing.T) {
	p := plugin.NewPlugin(
		"testPluginV3",
		"v1.0.0",
		memdb.NewMemDBClient)
	srv := Plugin(p, WithArgs("serve"), WithTestListener())
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	var serverErr error
	go func() {
		defer wg.Done()
		serverErr = srv.Serve(ctx)
	}()
	defer func() {
		cancel()
		wg.Wait()
	}()

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(srv.bufPluginDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	c := pb.NewPluginClient(conn)
	if _, err := c.Init(ctx, &pb.Init_Request{}); err != nil {
		t.Fatal(err)
	}

	table := state.Table("test_no_pk")
	// Remove PKs
	for i := range table.Columns {
		table.Columns[i].PrimaryKey = false
	}

	stateClient, err := state.NewClientWithTable(ctx, c, table)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range []string{"valua", "value1", "value3", "value2"} {
		if err := stateClient.SetKey(ctx, "key", v); err != nil {
			t.Fatal(err)
		}
		// Without Flush(), value will only be updated in memory, and we won't get duplicate entries in memdb
		if err := stateClient.Flush(ctx); err != nil {
			t.Fatal(err)
		}
	}

	stateClient, err = state.NewClientWithTable(ctx, c, table)
	if err != nil {
		t.Fatal(err)
	}
	val, err := stateClient.GetKey(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}
	if val != `value2` {
		t.Fatalf("expected value to be %q but got %q", `value2`, val)
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
