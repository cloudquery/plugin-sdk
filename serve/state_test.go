package serve

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/premium"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestState(t *testing.T) {
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
	// TODO: Remove once there's a documented migration path per https://github.com/grpc/grpc-go/issues/7244
	// nolint:staticcheck
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(srv.bufPluginDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	c := pb.NewPluginClient(conn)
	if _, err := c.Init(ctx, &pb.Init_Request{}); err != nil {
		t.Fatal(err)
	}
	stateClient, err := state.NewClient(ctx, conn, "test")
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
	stateClient, err = state.NewClient(ctx, conn, "test")
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

func TestStateOverwriteGetLatest(t *testing.T) {
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
	// TODO: Remove once there's a documented migration path per https://github.com/grpc/grpc-go/issues/7244
	// nolint:staticcheck
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

	stateClient, err := state.NewClient(ctx, conn, table.Name)
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

	stateClient, err = state.NewClient(ctx, conn, table.Name)
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

func TestStateCloseSavesUnflushedChanges(t *testing.T) {
	ctx := context.Background()
	dial := serveStateBackend(t)

	stateClient, err := state.NewClient(ctx, dial(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if err := stateClient.SetKey(ctx, "key", "value"); err != nil {
		t.Fatal(err)
	}
	if err := stateClient.Close(); err != nil {
		t.Fatal(err)
	}

	assertStateValue(ctx, t, dial(), "value")
}

func TestStateCloseDropsUnflushedChangesWhenSyncIsCancelled(t *testing.T) {
	ctx := context.Background()
	dial := serveStateBackend(t)

	syncCtx, cancel := context.WithCancel(ctx)
	stateClient, err := state.NewClient(syncCtx, dial(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if err := stateClient.SetKey(ctx, "key", "saved"); err != nil {
		t.Fatal(err)
	}
	if err := stateClient.Flush(ctx); err != nil {
		t.Fatal(err)
	}
	if err := stateClient.SetKey(ctx, "key", "dropped"); err != nil {
		t.Fatal(err)
	}
	cancel()
	if err := stateClient.Close(); err != nil {
		t.Fatal(err)
	}

	assertStateValue(ctx, t, dial(), "saved")
}

func TestStateCloseSavesUnflushedChangesWhenQuotaRunsOut(t *testing.T) {
	ctx := context.Background()
	dial := serveStateBackend(t)

	quota := &quotaMonitor{}
	syncCtx, err := premium.WithCancelOnQuotaExceeded(ctx, quota, premium.WithQuotaCheckPeriod(time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	stateClient, err := state.NewClient(syncCtx, dial(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if err := stateClient.SetKey(ctx, "key", "value"); err != nil {
		t.Fatal(err)
	}
	quota.exceeded.Store(true)
	<-syncCtx.Done()
	if err := stateClient.Close(); err != nil {
		t.Fatal(err)
	}

	assertStateValue(ctx, t, dial(), "value")
}

type quotaMonitor struct {
	exceeded atomic.Bool
}

func (*quotaMonitor) TeamName() string {
	return "test-team"
}

func (q *quotaMonitor) CheckQuota(context.Context) (premium.CheckQuotaResult, error) {
	return premium.CheckQuotaResult{HasQuota: !q.exceeded.Load()}, nil
}

func assertStateValue(ctx context.Context, t *testing.T, conn *grpc.ClientConn, want string) {
	t.Helper()
	stateClient, err := state.NewClient(ctx, conn, "test")
	if err != nil {
		t.Fatal(err)
	}
	val, err := stateClient.GetKey(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}
	if val != want {
		t.Fatalf("expected value to be %q but got %q", want, val)
	}
}

// serveStateBackend serves a memdb destination and returns a func that opens a new connection to it.
// Each state client needs its own connection, because Close closes it.
func serveStateBackend(t *testing.T) func() *grpc.ClientConn {
	t.Helper()
	p := plugin.NewPlugin("testPluginV3", "v1.0.0", memdb.NewMemDBClient)
	srv := Plugin(p, WithArgs("serve"), WithTestListener())
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Serve(ctx); err != nil {
			t.Error(err)
		}
	}()
	t.Cleanup(func() {
		cancel()
		wg.Wait()
	})

	dial := func() *grpc.ClientConn {
		// nolint:staticcheck
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(srv.bufPluginDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		t.Cleanup(func() { _ = conn.Close() })
		return conn
	}
	if _, err := pb.NewPluginClient(dial()).Init(ctx, &pb.Init_Request{}); err != nil {
		t.Fatal(err)
	}
	return dial
}
