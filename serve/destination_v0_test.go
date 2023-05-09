package serve

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	pbBase "github.com/cloudquery/plugin-pb-go/pb/base/v0"
	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v0"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v2/internal/deprecated"
	"github.com/cloudquery/plugin-sdk/v2/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v2/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/testdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	c := pb.NewDestinationClient(conn)
	spec := specs.Destination{
		WriteMode: specs.WriteModeAppend,
	}
	specBytes, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Configure(ctx, &pbBase.Configure_Request{Config: specBytes}); err != nil {
		t.Fatal(err)
	}

	getNameRes, err := c.GetName(ctx, &pbBase.GetName_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getNameRes.Name != "testDestinationPlugin" {
		t.Fatalf("expected name to be testDestinationPlugin but got %s", getNameRes.Name)
	}

	getVersionRes, err := c.GetVersion(ctx, &pbBase.GetVersion_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getVersionRes.Version != "development" {
		t.Fatalf("expected version to be development but got %s", getVersionRes.Version)
	}

	tableName := "test_destination_serve"
	sourceName := "test_destination_serve_source"
	syncTime := time.Now()
	table := testdata.TestTable(tableName)
	tables := schema.Tables{table}
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	tablesBytes, err := json.Marshal(tables)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Migrate(ctx, &pb.Migrate_Request{
		Tables: tablesBytes,
	}); err != nil {
		t.Fatal(err)
	}

	destResource := schema.DestinationResource{
		TableName: tableName,
		Data:      deprecated.GenTestData(table),
	}
	_ = destResource.Data[0].Set(sourceName)
	_ = destResource.Data[1].Set(syncTime)
	destRecord := schema.CQTypesOneToRecord(memory.DefaultAllocator, destResource.Data, table.ToArrowSchema())
	destResourceBytes, err := json.Marshal(destResource)
	if err != nil {
		t.Fatal(err)
	}
	sourceSpecBytes, err := json.Marshal(sourceSpec)
	if err != nil {
		t.Fatal(err)
	}
	writeClient, err := c.Write2(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeClient.Send(&pb.Write2_Request{
		SourceSpec: sourceSpecBytes,
		Source:     sourceSpec.Name,
		Timestamp:  timestamppb.New(syncTime.Truncate(time.Microsecond)),
		Tables:     tablesBytes,
	}); err != nil {
		t.Fatal(err)
	}

	if err := writeClient.Send(&pb.Write2_Request{
		Resource: destResourceBytes,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := writeClient.CloseAndRecv(); err != nil {
		t.Fatal(err)
	}

	readCh := make(chan arrow.Record, 1)
	if err := plugin.Read(ctx, table.ToArrowSchema(), sourceName, readCh); err != nil {
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
	if _, err := c.DeleteStale(ctx, &pb.DeleteStale_Request{
		Source:    "testSource",
		Timestamp: timestamppb.New(time.Now().Truncate(time.Microsecond)),
		Tables:    tablesBytes,
	}); err != nil {
		t.Fatal(err)
	}

	_, err = c.GetMetrics(ctx, &pb.GetDestinationMetrics_Request{})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.Close(ctx, &pb.Close_Request{}); err != nil {
		t.Fatalf("failed to call Close: %v", err)
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
