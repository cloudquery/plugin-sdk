package discovery

import (
	"context"

	pb "github.com/cloudquery/plugin-sdk/internal/pb/discovery/v0"
)

type DiscoveryServer struct {
	pb.UnimplementedDiscoveryServer
	Versions []string
}

func (s *DiscoveryServer) GetVersions(context.Context, *pb.GetVersions_Request) (*pb.GetVersions_Response, error) {
	return &pb.GetVersions_Response{Versions: s.Versions}, nil
}