package serve

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	clients "github.com/cloudquery/plugin-sdk/v2/clients/destination/v0"
	"github.com/cloudquery/plugin-sdk/v2/internal/deprecated"
	"github.com/cloudquery/plugin-sdk/v2/internal/memdb"
	servers "github.com/cloudquery/plugin-sdk/v2/internal/servers/destination/v0"
	"github.com/cloudquery/plugin-sdk/v2/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bufDestinationDialer(context.Context, string) (net.Conn, error) {
	testDestinationListenerLock.Lock()
	defer testDestinationListenerLock.Unlock()
	return testDestinationListener.Dial()
}

func TestDestination(t *testing.T) {
	plugin := destination.NewPlugin("testDestinationPlugin", "development", memdb.NewClient)
	s := &destinationServe{
		plugin: plugin,
	}
	cmd := newCmdDestinationRoot(s)
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

	// wait for the server to start
	for {
		testDestinationListenerLock.Lock()
		if testDestinationListener != nil {
			testDestinationListenerLock.Unlock()
			break
		}
		testDestinationListenerLock.Unlock()
		t.Log("waiting for grpc server to start")
		time.Sleep(time.Millisecond * 200)
	}

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDestinationDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	c, err := clients.NewClient(ctx, specs.RegistryGrpc, "", "", clients.WithGrpcConn(conn), clients.WithNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.Terminate(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := c.Initialize(ctx, specs.Destination{
		WriteMode: specs.WriteModeAppend,
	}); err != nil {
		t.Fatal(err)
	}

	name, err := c.Name(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if name != "testDestinationPlugin" {
		t.Fatalf("expected name to be testDestinationPlugin but got %s", name)
	}

	version, err := c.Version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if version != "development" {
		t.Fatalf("expected version to be development but got %s", version)
	}

	tableName := "test_destination_serve"
	sourceName := "test_destination_serve_source"
	syncTime := time.Now()
	table := schema.TestTable(tableName)
	tables := schema.Tables{table}
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	if err := c.Migrate(ctx, tables); err != nil {
		t.Fatal(err)
	}

	destResource := schema.DestinationResource{
		TableName: tableName,
		Data:      deprecated.GenTestData(table),
	}
	_ = destResource.Data[0].Set(sourceName)
	_ = destResource.Data[1].Set(syncTime)
	tableV2 := servers.TableV1ToV2(table)
	destRecord := schema.CQTypesOneToRecord(memory.DefaultAllocator, destResource.Data, tableV2.ToArrowSchema())
	b, err := json.Marshal(destResource)
	if err != nil {
		t.Fatal(err)
	}

	resources := make(chan []byte, 1)
	resources <- b
	close(resources)
	if err := c.Write2(ctx, sourceSpec, tables, syncTime, resources); err != nil {
		t.Fatal(err)
	}

	readCh := make(chan arrow.Record, 1)
	if err := plugin.Read(ctx, tableV2, sourceName, readCh); err != nil {
		t.Fatal(err)
	}
	close(readCh)
	totalResources := 0
	for resource := range readCh {
		totalResources++
		if !array.RecordEqual(destRecord, resource) {
			diff := destination.RecordDiff(destRecord, resource)
			t.Fatalf("expected %v but got %v. Diff: %v", destRecord, resource, diff)
		}
	}
	if totalResources != 1 {
		t.Fatalf("expected 1 resource but got %d", totalResources)
	}

	if err := c.DeleteStale(ctx, nil, "testSource", time.Now()); err != nil {
		t.Fatalf("failed to call DeleteStale: %v", err)
	}

	_, err = c.GetMetrics(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Close(ctx); err != nil {
		t.Fatalf("failed to call Close: %v", err)
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
