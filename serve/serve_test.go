package serve

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ schema.ClientMeta = &testExecutionClient{}

type testExecutionClient struct {
	logger zerolog.Logger
}

func testTable() *schema.Table {
	return &schema.Table{
		Name: "testTable",
		Resolver: func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
			res <- map[string]interface{}{
				"TestColumn": 3,
			}
			return nil
		},
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: schema.TypeInt,
			},
		},
	}
}

func (c *testExecutionClient) Logger() *zerolog.Logger {
	return &c.logger
}

func newTestExecutionClient(context.Context, *plugins.SourcePlugin, specs.SourceSpec) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

// https://stackoverflow.com/questions/32840687/timeout-for-waitgroup-wait
func waitTimeout(wg *errgroup.Group, timeout time.Duration) (bool, error) {
	c := make(chan struct{})
	var err error
	go func() {
		defer close(c)
		err = wg.Wait()
	}()
	select {
	case <-c:
		return false, err // completed normally
	case <-time.After(timeout):
		return true, err // timed out
	}
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return testListener.Dial()
}

func TestServe(t *testing.T) {
	plugin := plugins.NewSourcePlugin(
		"testSourcePlugin",
		"1.0.0",
		[]*schema.Table{testTable()},
		newTestExecutionClient,
		plugins.WithSourceLogger(zerolog.New(zerolog.NewTestWriter(t))))

	cmd := newCmdRoot(Options{
		SourcePlugin: plugin,
	})
	cmd.SetArgs([]string{"serve", "--network", "test"})

	go func() {
		cmd.Execute()
	}()

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	ctx := context.Background()
	_, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	// c := clients.NewSourceClient(conn)
	// c.

	// g := errgroup.Group{}
	// g.Go(func() error {
	// 	return cmd.Execute()
	// })

	// // there is no programmatic way to shutdown server so we just check if returned an
	// if waitTimeout(&g, time.Second*3) {
	// 	t.Fatal("timed out")
	// }
	// if err := g.Wait(); err != nil {
	// 	t.Fatal(err)
	// }
}
