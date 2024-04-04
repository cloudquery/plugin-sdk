package serve

import (
	"context"
	"sync"
	"testing"

	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
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

func TestStateOverwrite(t *testing.T) {
	cases := []struct {
		name   string
		values []string
		expect string
	}{
		{
			"Overwrite",
			[]string{"valua1", "value1", "value3", "value2"}, // All same length, expect largest value
			"value3", // Largest value lexicographically
		},
		{
			"Overwrite with integers",
			[]string{"1", "32", "4"},
			"1", // First value written, last value read from memdb?
		},
		{
			"Overwrite with timestamps",
			[]string{"2024-04-03T16:02:55.20412Z", "2024-04-03T16:03:37.06Z", "2024-04-03T16:03:22.440487Z", "2024-04-03T16:03:37.058413Z"},
			"2024-04-03T16:03:37.06Z", // Latest timestamp despite rounding zeroes
		},
		{
			"Overwrite with float unix timestamps",
			[]string{"1712226133.860000", "1712226134.759", "1712226133.859000"},
			"1712226134.759",
		},
		{
			"Overwrite with float unix timestamps with int in between",
			[]string{"1712226133.860000", "1712226134", "1712226133.859000"},
			"1712226134",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
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

			for _, v := range tc.values {
				if err := stateClient.SetKey(ctx, "key", v); err != nil {
					t.Fatal(err)
				}
				// Without Flush(), value will only be updated in memory, and we won't get duplicate entries in memdb
				if err := stateClient.Flush(ctx); err != nil {
					t.Fatal(err)
				}
			}

			val, err := stateClient.GetKey(ctx, "key")
			if err != nil {
				t.Fatal(err)
			}
			if finalValueWritten := tc.values[len(tc.values)-1]; val != finalValueWritten {
				t.Fatalf("expected value to be %q but got %q", finalValueWritten, val)
			}

			stateClient, err = state.NewClientWithTable(ctx, c, table)
			if err != nil {
				t.Fatal(err)
			}
			val, err = stateClient.GetKey(ctx, "key")
			if err != nil {
				t.Fatal(err)
			}
			if val != tc.expect {
				t.Fatalf("expected value to be %q but got %q", tc.expect, val)
			}

			cancel()
			wg.Wait()
			if serverErr != nil {
				t.Fatal(serverErr)
			}
		})
	}
}
