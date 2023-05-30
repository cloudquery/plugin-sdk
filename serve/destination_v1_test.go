package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/ipc"
	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v1"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDestinationV1(t *testing.T) {
	plugin := plugin.NewPlugin("testDestinationPlugin", "development", memdb.NewClient)
	s := &pluginServe{
		plugin: plugin,
	}
	cmd := newCmdPluginRoot(s)
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
		testPluginListenerLock.Lock()
		if testPluginListener != nil {
			testPluginListenerLock.Unlock()
			break
		}
		testPluginListenerLock.Unlock()
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
	if _, err := c.Configure(ctx, &pb.Configure_Request{Config: specBytes}); err != nil {
		t.Fatal(err)
	}

	getNameRes, err := c.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getNameRes.Name != "testDestinationPlugin" {
		t.Fatalf("expected name to be testDestinationPlugin but got %s", getNameRes.Name)
	}

	getVersionRes, err := c.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getVersionRes.Version != "development" {
		t.Fatalf("expected version to be development but got %s", getVersionRes.Version)
	}

	tableName := "test_destination_serve"
	sourceName := "test_destination_serve_source"
	syncTime := time.Now()
	table := schema.TestTable(tableName, schema.TestSourceOptions{})
	tables := schema.Tables{table}
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	encodedTables, err := tables.ToArrowSchemas().Encode()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.Migrate(ctx, &pb.Migrate_Request{
		Tables: encodedTables,
	}); err != nil {
		t.Fatal(err)
	}

	rec := schema.GenTestData(table, schema.GenTestDataOptions{
		SourceName: sourceName,
		SyncTime:   syncTime,
		MaxRows:    1,
	})[0]

	sourceSpecBytes, err := json.Marshal(sourceSpec)
	if err != nil {
		t.Fatal(err)
	}
	writeClient, err := c.Write(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeClient.Send(&pb.Write_Request{
		SourceSpec: sourceSpecBytes,
		Source:     sourceSpec.Name,
		Timestamp:  timestamppb.New(syncTime.Truncate(time.Microsecond)),
		Tables:     encodedTables,
	}); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	wr := ipc.NewWriter(&buf, ipc.WithSchema(rec.Schema()))
	if err := wr.Write(rec); err != nil {
		t.Fatal(err)
	}
	if err := wr.Close(); err != nil {
		t.Fatal(err)
	}
	if err := writeClient.Send(&pb.Write_Request{
		Resource: buf.Bytes(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := writeClient.CloseAndRecv(); err != nil {
		t.Fatal(err)
	}
	// serversDestination
	readCh := make(chan arrow.Record, 1)
	if err := plugin.Read(ctx, table, sourceName, readCh); err != nil {
		t.Fatal(err)
	}
	close(readCh)
	totalResources := 0
	for resource := range readCh {
		totalResources++
		if !array.RecordEqual(rec, resource) {
			diff := destination.RecordDiff(rec, resource)
			t.Fatalf("expected %v but got %v. Diff: %v", rec, resource, diff)
		}
	}
	if totalResources != 1 {
		t.Fatalf("expected 1 resource but got %d", totalResources)
	}
	if _, err := c.DeleteStale(ctx, &pb.DeleteStale_Request{
		Source:    "testSource",
		Timestamp: timestamppb.New(time.Now().Truncate(time.Microsecond)),
		Tables:    encodedTables,
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
