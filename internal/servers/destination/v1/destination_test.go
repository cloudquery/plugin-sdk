package destination

import (
	"context"
	"testing"

	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v1"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

func TestGetName(t *testing.T) {
	ctx := context.Background()
	s := Server{
		Plugin: plugin.NewPlugin("test", "development", memdb.NewMemDBClient),
	}
	res, err := s.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Name != "test" {
		t.Fatalf("expected test, got %s", res.GetName())
	}
}

func TestGetVersion(t *testing.T) {
	ctx := context.Background()
	s := Server{
		Plugin: plugin.NewPlugin("test", "development", memdb.NewMemDBClient),
	}
	resVersion, err := s.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if resVersion.Version != "development" {
		t.Fatalf("expected development, got %s", resVersion.GetVersion())
	}
}
