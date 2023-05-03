package serve

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"testing"
	"time"

	pb "github.com/cloudquery/plugin-pb-go/pb/source/v1"
	"github.com/cloudquery/plugin-sdk/v2/plugins/source"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSourceSuccessV1(t *testing.T) {
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

	getNameResponse, err := c.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getNameResponse.Name != "testPlugin" {
		t.Fatalf("expected name to be testPlugin but got %s", getNameResponse.Name)
	}

	getVersionRes, err := c.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getVersionRes.Version != "v1.0.0" {
		t.Fatalf("Expected version to be v1.0.0 but got %s", getVersionRes.Version)
	}

	getTablesRes, err := c.GetTables(ctx, &pb.GetTables_Request{})
	if err != nil {
		t.Fatal(err)
	}
	var tables schema.Tables
	if err := json.Unmarshal(getTablesRes.Tables, &tables); err != nil {
		t.Fatal(err)
	}
	if len(tables) != 2 {
		t.Fatalf("Expected 2 tables but got %d", len(tables))
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
	specBytes, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Init(ctx, &pb.Init_Request{Spec: specBytes}); err != nil {
		t.Fatal(err)
	}

	syncClient, err := c.Sync(ctx, &pb.Sync_Request{})
	if err != nil {
		t.Fatal(err)
	}
	var resources []schema.DestinationResource
	for {
		r, err := syncClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		var resource schema.DestinationResource
		if err := json.Unmarshal(r.Resource, &resource); err != nil {
			t.Fatal(err)
		}
		resources = append(resources, resource)
	}

	totalResources := 0
	for _, resource := range resources {
		if resource.TableName != "test_table" {
			t.Fatalf("Expected resource with table name test_table. got: %s", resource.TableName)
		}
		if len(resource.Data) != 3 {
			t.Fatalf("Expected resource with data length 3 but got %d", len(resource.Data))
		}

		if resource.Data[2] == nil {
			t.Fatalf("Expected resource with data[2] to be not nil")
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
