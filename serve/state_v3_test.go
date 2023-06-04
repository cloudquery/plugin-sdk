package serve

import (
	"context"
	"sync"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/internal/state"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestStateV3(t *testing.T) {
	p := plugin.NewPlugin("memdb", "v1.0.0", plugin.NewMemDBClient)
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

	stateClient, err := state.NewClient(ctx, "test", conn)
	if err != nil {
		t.Fatalf("Failed to create state client: %v", err)
	}
	if err := stateClient.SetKey(ctx, "testKey", "testValue"); err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}
	key, err := stateClient.GetKey(ctx, "testKey")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if key != "testValue" {
		t.Fatalf("Unexpected key value: %v", key)
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}