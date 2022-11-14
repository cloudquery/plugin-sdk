package serve

import (
	"context"
	"encoding/json"
	"fmt"
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

type TestSourcePluginSpec struct {
	Accounts []string `json:"accounts,omitempty" yaml:"accounts,omitempty"`
}

type testExecutionClient struct{}

var _ schema.ClientMeta = &testExecutionClient{}

var errTestExecutionClientErr = fmt.Errorf("error in newTestExecutionClientErr")

var expectedExampleSpecConfig = specs.Spec{
	Kind: specs.KindSource,
	Spec: &specs.Source{
		Name:    "testSourcePlugin",
		Path:    "cloudquery/testSourcePlugin",
		Version: "v1.0.0",
		Tables:  []string{"*"},
		Spec:    map[string]interface{}{"accounts": []interface{}{"all"}},
	},
}

func testTable() *schema.Table {
	return &schema.Table{
		Name: "test_table",
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

func (*testExecutionClient) ID() string {
	return "testExecutionClient"
}

func newTestExecutionClient(context.Context, zerolog.Logger, specs.Source) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

func newTestExecutionClientErr(context.Context, zerolog.Logger, specs.Source) (schema.ClientMeta, error) {
	return nil, errTestExecutionClientErr
}

func bufSourceDialer(context.Context, string) (net.Conn, error) {
	testSourceListenerLock.Lock()
	defer testSourceListenerLock.Unlock()
	return testSourceListener.Dial()
}

func TestSourceSuccess(t *testing.T) {
	plugin := plugins.NewSourcePlugin(
		"testSourcePlugin",
		"v1.0.0",
		[]*schema.Table{testTable()},
		newTestExecutionClient)

	cmd := newCmdSourceRoot(&sourceServe{
		plugin: plugin,
	})
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
	for {
		testSourceListenerLock.Lock()
		if testSourceListener != nil {
			testSourceListenerLock.Unlock()
			break
		}
		testSourceListenerLock.Unlock()
		t.Log("waiting for grpc server to start")
		time.Sleep(time.Millisecond * 200)
	}

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufSourceDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	c, err := clients.NewSourceClient(ctx, specs.RegistryGrpc, "", "", clients.WithSourceGRPCConnection(conn), clients.WithSourceNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.Terminate(); err != nil {
			t.Fatal(err)
		}
	}()

	name, err := c.Name(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if name != "testSourcePlugin" {
		t.Fatalf("expected name to be testSourcePlugin but got %s", name)
	}

	version, err := c.Version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if version != "v1.0.0" {
		t.Fatalf("Expected version to be v1.0.0 but got %s", version)
	}

	tables, err := c.GetTables(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(tables) != 1 {
		t.Fatalf("Expected 1 table but got %d", len(tables))
	}

	resources := make(chan []byte, 1)
	if err := c.Sync2(ctx,
		specs.Source{
			Name:         "testSourcePlugin",
			Version:      "v1.0.0",
			Path:         "cloudquery/testSourcePlugin",
			Registry:     specs.RegistryGithub,
			Tables:       []string{"*"},
			Spec:         TestSourcePluginSpec{Accounts: []string{"cloudquery/plugin-sdk"}},
			Destinations: []string{"test"},
		},
		resources); err != nil {
		t.Fatal(err)
	}
	close(resources)

	totalResources := 0
	for resourceB := range resources {
		var resource schema.DestinationResource
		if err := json.Unmarshal(resourceB, &resource); err != nil {
			t.Fatalf("failed to unmarshal resource: %v", err)
		}
		if resource.TableName != "test_table" {
			t.Fatalf("Expected resource with table name test_table. got: %s", resource.TableName)
		}
		if len(resource.Data) != 3 {
			t.Fatalf("Expected resource with data length 3 but got %d", len(resource.Data))
		}
		fmt.Println(resource.Data)
		if resource.Data[2] == nil {
			t.Fatalf("Expected resource with data[2] to be not nil")
		}
		// if resource.Data[2].Type() != schema.TypeInt {
		// 	t.Fatalf("Expected resource with data type int but got %s", resource.Data[2].Type())
		// }
		totalResources++
	}
	if totalResources != 1 {
		t.Fatalf("Expected 1 resource on channel but got %d", totalResources)
	}

	stats, err := c.GetMetrics(ctx)
	if err != nil {
		t.Fatal(err)
	}
	clientStats := stats.TableClient[""][""]
	if clientStats.Resources != 1 {
		t.Fatalf("Expected 1 resource but got %d", clientStats.Resources)
	}

	if clientStats.Errors != 0 {
		t.Fatalf("Expected 0 errors but got %d", clientStats.Errors)
	}

	if clientStats.Panics != 0 {
		t.Fatalf("Expected 0 panics but got %d", clientStats.Panics)
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}

const testSourceFailExpectedErr = "failed to fetch resources from stream: rpc error: code = Unknown desc = failed to sync resources: failed to create execution client for source plugin testSourcePlugin: error in newTestExecutionClientErr"

func TestSourceFail(t *testing.T) {
	plugin := plugins.NewSourcePlugin(
		"testSourcePlugin",
		"v1.0.0",
		[]*schema.Table{testTable()},
		newTestExecutionClientErr)

	cmd := newCmdSourceRoot(&sourceServe{
		plugin: plugin,
	})
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
	for {
		testSourceListenerLock.Lock()
		if testSourceListener != nil {
			testSourceListenerLock.Unlock()
			break
		}
		testSourceListenerLock.Unlock()
		t.Log("waiting for grpc server to start")
		time.Sleep(time.Millisecond * 200)
	}

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufSourceDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	c, err := clients.NewSourceClient(ctx, specs.RegistryGrpc, "", "", clients.WithSourceGRPCConnection(conn), clients.WithSourceNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.Terminate(); err != nil {
			t.Fatal(err)
		}
	}()

	resources := make(chan []byte, 1)
	err = c.Sync2(ctx,
		specs.Source{
			Name:         "testSourcePlugin",
			Version:      "v1.0.0",
			Path:         "cloudquery/testSourcePlugin",
			Registry:     specs.RegistryGithub,
			Tables:       []string{"*"},
			Spec:         TestSourcePluginSpec{Accounts: []string{"cloudquery/plugin-sdk"}},
			Destinations: []string{"test"},
		},
		resources)
	close(resources)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err.Error() != testSourceFailExpectedErr {
		t.Fatalf("expected error %s but got %v", testSourceFailExpectedErr, err)
	}
	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
