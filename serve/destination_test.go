package serve

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/clients"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bufDestinationDialer(context.Context, string) (net.Conn, error) {
	return testDestinationListener.Dial()
}

type testDestinationClient struct {
}

func (*testDestinationClient) Name() string {
	return "testDestinationPlugin"
}
func (*testDestinationClient) Version() string {
	return "development"
}
func (*testDestinationClient) Initialize(context.Context, specs.Destination) error {
	return nil
}
func (*testDestinationClient) Migrate(context.Context, schema.Tables) error {
	return nil
}
func (*testDestinationClient) Write(context.Context, string, map[string]interface{}) error {
	return nil
}
func (*testDestinationClient) SetLogger(zerolog.Logger) {
}
func (*testDestinationClient) Close(context.Context) error {
	return nil
}

func TestDestination(t *testing.T) {
	plugin := plugins.DestinationPlugin(&testDestinationClient{})
	s := &destinationServe{
		plugin: plugin,
	}

	cmd := newCmdDestinationRoot(s)
	cmd.SetArgs([]string{"serve", "--network", "test"})

	var serveErr error
	go func() {
		serveErr = cmd.Execute()
	}()

	// wait for the server to start
	for {
		if testDestinationListener != nil {
			break
		}
		t.Log("waiting for grpc server to start")
		time.Sleep(time.Millisecond * 200)
		if serveErr != nil {
			t.Fatal(serveErr)
		}
	}

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDestinationDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	c, err := clients.NewDestinationClient(ctx, specs.RegistryGrpc, "", "", clients.WithDestinationGrpcConn(conn))
	if err != nil {
		t.Fatal(err)
	}
	resources := make(chan []byte)
	wg := errgroup.Group{}
	wg.Go(func() error {
		defer close(resources)
		name, err := c.Name(ctx)
		if err != nil {
			return err
		}
		if name != "testDestinationPlugin" {
			return fmt.Errorf("expected name to be testDestinationPlugin but got %s", name)
		}
		// call all methods as sanity check
		if err := c.Close(); err != nil {
			return fmt.Errorf("failed to close: %w", err)
		}
		return nil
	})
	if err := wg.Wait(); err != nil {
		t.Fatalf("Failed to get name from destination plugin: %v", err)
	}
}
