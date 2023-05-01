package discovery

import (
	"context"

	pb "github.com/cloudquery/plugin-sdk/v3/pb/discovery/v0"
)

type Server struct {
	pb.UnimplementedDiscoveryServer
	Versions []string
}

func (s *Server) GetVersions(context.Context, *pb.GetVersions_Request) (*pb.GetVersions_Response, error) {
	return &pb.GetVersions_Response{Versions: s.Versions}, nil
}
