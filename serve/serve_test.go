package serve

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/clients"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestSourcePluginSpec struct {
	Accounts []string `json:"accounts,omitempty" yaml:"accounts,omitempty"`
}

type testExecutionClient struct {
	logger zerolog.Logger
}

var _ schema.ClientMeta = &testExecutionClient{}

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

func (c *testExecutionClient) Logger() *zerolog.Logger {
	return &c.logger
}

func newTestExecutionClient(context.Context, zerolog.Logger, specs.Source) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

// https://stackoverflow.com/questions/32840687/timeout-for-waitgroup-wait
// func waitTimeout(wg *errgroup.Group, timeout time.Duration) (bool, error) {
// 	c := make(chan struct{})
// 	var err error
// 	go func() {
// 		defer close(c)
// 		err = wg.Wait()
// 	}()
// 	select {
// 	case <-c:
// 		return false, err // completed normally
// 	case <-time.After(timeout):
// 		return true, err // timed out
// 	}
// }

func bufDialer(context.Context, string) (net.Conn, error) {
	return testListener.Dial()
}

func TestServe(t *testing.T) {
	plugin := plugins.NewSourcePlugin(
		"testSourcePlugin",
		"v1.0.0",
		[]*schema.Table{testTable()},
		newTestExecutionClient,
		plugins.WithSourceExampleConfig(`--- 
kind: source
spec: 
  name: testSourcePlugin
  path: cloudquery/testSourcePlugin
  spec: 
    accounts: ["all"]
  tables: 
    - "*"
  version: v1.0.0
`),
		plugins.WithSourceLogger(zerolog.New(zerolog.NewTestWriter(t))),
	)

	cmd := newCmdRoot(Options{
		SourcePlugin: plugin,
	})
	cmd.SetArgs([]string{"serve", "--network", "test"})

	var serveErr error
	go func() {
		serveErr = cmd.Execute()
	}()

	// wait for the server to start
	for {
		if testListener != nil {
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
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	c := clients.NewSourceClient(conn)
	resources := make(chan *schema.Resource)
	wg := errgroup.Group{}
	wg.Go(func() error {
		defer close(resources)
		return c.Sync(ctx,
			specs.Source{
				Name:     "testSourcePlugin",
				Version:  "v1.0.0",
				Registry: specs.RegistryGithub,
				Spec:     TestSourcePluginSpec{Accounts: []string{"cloudquery/plugin-sdk"}},
			},
			resources)
	})
	for resource := range resources {
		if resource.TableName != "test_table" {
			t.Fatalf("Expected resource with table name test: %s", resource.TableName)
		}
		if int(resource.Data["test_column"].(float64)) != 3 {
			t.Fatalf("Expected resource {'test_column':3} got: %v", resource.Data)
		}
	}
	if err := wg.Wait(); err != nil {
		t.Fatalf("Failed to sync resources: %v", err)
	}

	exampleConfig, err := c.ExampleConfig(ctx)
	if err != nil {
		t.Fatalf("Failed to get example config: %v", err)
	}
	var exampleSpec specs.Spec
	if err := specs.SpecUnmarshalYamlStrict([]byte(exampleConfig), &exampleSpec); err != nil {
		t.Fatalf("Failed to unmarshal example config: %v", err)
	}
	// skip internal validation for now

	if diff := cmp.Diff(expectedExampleSpecConfig, exampleSpec); diff != "" {
		t.Fatalf("Spec mismatch (-want +got):\n%s", diff)
	}
}
