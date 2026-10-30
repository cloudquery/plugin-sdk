package serve

import (
	"context"
	"slices"
	"sync"
	"testing"

	pbdiscoveryv0 "github.com/cloudquery/plugin-pb-go/pb/discovery/v0"
	pbdiscoveryv1 "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDestinationDiscoveryOnlyAdvertisesV3(t *testing.T) {
	p := plugin.NewPlugin("testDestinationPlugin", "development", memdb.NewMemDBClient, plugin.WithKind("destination"))
	srv := Plugin(p, WithArgs("serve"), WithDestinationV0V1Server(), WithTestListener())
	ctx, cancel := context.WithCancel(context.Background())
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

	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(srv.bufPluginDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	v1Res, err := pbdiscoveryv1.NewDiscoveryClient(conn).GetVersions(ctx, &pbdiscoveryv1.GetVersions_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(v1Res.Versions, []int32{3}) {
		t.Fatalf("expected discovery v1 versions [3] but got %v", v1Res.Versions)
	}

	v0Res, err := pbdiscoveryv0.NewDiscoveryClient(conn).GetVersions(ctx, &pbdiscoveryv0.GetVersions_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(v0Res.Versions, []string{"v3"}) {
		t.Fatalf("expected discovery v0 versions [v3] but got %v", v0Res.Versions)
	}

	for _, method := range []string{"/proto.Destination/GetName", "/cloudquery.destination.v1.Destination/GetName"} {
		err := conn.Invoke(ctx, method, &emptypb.Empty{}, &emptypb.Empty{})
		if status.Code(err) != codes.Unimplemented {
			t.Fatalf("expected %s to return %s but got %v", method, codes.Unimplemented, err)
		}
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
