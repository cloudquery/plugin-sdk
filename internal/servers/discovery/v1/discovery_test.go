package discovery

import (
	"context"
	"testing"

	pb "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
)

func TestDiscovery(t *testing.T) {
	ctx := context.Background()
	s := &Server{
		Versions: []int32{1, 2},
	}
	resp, err := s.GetVersions(ctx, &pb.GetVersions_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Versions) != 2 {
		t.Fatal("expected 2 versions")
	}
	if resp.Versions[0] != 1 {
		t.Fatal("expected version 1")
	}
	if resp.Versions[1] != 2 {
		t.Fatal("expected version 2")
	}
}
