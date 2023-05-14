package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/ipc"
	pb "github.com/cloudquery/plugin-pb-go/pb/source/v2"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/plugins/source"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestSourcePluginSpec struct {
	Accounts []string `json:"accounts,omitempty" yaml:"accounts,omitempty"`
}

type testExecutionClient struct{}

var _ schema.ClientMeta = &testExecutionClient{}

// var errTestExecutionClientErr = fmt.Errorf("error in newTestExecutionClientErr")

func testTable(name string) *schema.Table {
	return &schema.Table{
		Name: name,
		Resolver: func(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
			res <- map[string]any{
				"TestColumn": 3,
			}
			return nil
		},
		Columns: []schema.Column{
			{
				Name: "test_column",
				Type: arrow.PrimitiveTypes.Int64,
			},
		},
	}
}

func (*testExecutionClient) ID() string {
	return "testExecutionClient"
}

func newTestExecutionClient(context.Context, zerolog.Logger, specs.Source, source.Options) (schema.ClientMeta, error) {
	return &testExecutionClient{}, nil
}

func bufSourceDialer(context.Context, string) (net.Conn, error) {
	testSourceListenerLock.Lock()
	defer testSourceListenerLock.Unlock()
	return testSourceListener.Dial()
}

func TestSourceSuccess(t *testing.T) {
	plugin := source.NewPlugin(
		"testPlugin",
		"v1.0.0",
		[]*schema.Table{testTable("test_table"), testTable("test_table2")},
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
	c := pb.NewSourceClient(conn)

	getNameRes, err := c.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getNameRes.Name != "testPlugin" {
		t.Fatalf("expected name to be testPlugin but got %s", getNameRes.Name)
	}

	getVersionResponse, err := c.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getVersionResponse.Version != "v1.0.0" {
		t.Fatalf("Expected version to be v1.0.0 but got %s", getVersionResponse.Version)
	}

	spec := specs.Source{
		Name:         "testSourcePlugin",
		Version:      "v1.0.0",
		Path:         "cloudquery/testSourcePlugin",
		Registry:     specs.RegistryGithub,
		Tables:       []string{"test_table"},
		Spec:         TestSourcePluginSpec{Accounts: []string{"cloudquery/plugin-sdk"}},
		Destinations: []string{"test"},
	}
	specMarshaled, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("Failed to marshal spec: %v", err)
	}

	getTablesRes, err := c.GetTables(ctx, &pb.GetTables_Request{})
	if err != nil {
		t.Fatal(err)
	}

	tables, err := schema.NewTablesFromBytes(getTablesRes.Tables)
	if err != nil {
		t.Fatal(err)
	}

	if len(tables) != 2 {
		t.Fatalf("Expected 2 tables but got %d", len(tables))
	}
	if _, err := c.Init(ctx, &pb.Init_Request{Spec: specMarshaled}); err != nil {
		t.Fatal(err)
	}

	getTablesForSpecRes, err := c.GetDynamicTables(ctx, &pb.GetDynamicTables_Request{})
	if err != nil {
		t.Fatal(err)
	}
	tables, err = schema.NewTablesFromBytes(getTablesForSpecRes.Tables)
	if err != nil {
		t.Fatal(err)
	}

	if len(tables) != 1 {
		t.Fatalf("Expected 1 table but got %d", len(tables))
	}

	syncClient, err := c.Sync(ctx, &pb.Sync_Request{})
	if err != nil {
		t.Fatal(err)
	}
	var resources []arrow.Record
	for {
		r, err := syncClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		rdr, err := ipc.NewReader(bytes.NewReader(r.Resource))
		if err != nil {
			t.Fatal(err)
		}
		for rdr.Next() {
			rec := rdr.Record()
			rec.Retain()
			resources = append(resources, rec)
		}
	}

	totalResources := 0
	for _, resource := range resources {
		sc := resource.Schema()
		tableName, ok := sc.Metadata().GetValue(schema.MetadataTableName)
		if !ok {
			t.Fatal("Expected table name metadata to be set")
		}
		if tableName != "test_table" {
			t.Fatalf("Expected resource with table name test_table. got: %s", tableName)
		}
		if len(resource.Columns()) != 3 {
			t.Fatalf("Expected resource with data length 3 but got %d", len(resource.Columns()))
		}
		totalResources++
	}
	if totalResources != 1 {
		t.Fatalf("Expected 1 resource on channel but got %d", totalResources)
	}

	getMetricsRes, err := c.GetMetrics(ctx, &pb.GetMetrics_Request{})
	if err != nil {
		t.Fatal(err)
	}
	var stats source.Metrics
	if err := json.Unmarshal(getMetricsRes.Metrics, &stats); err != nil {
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
