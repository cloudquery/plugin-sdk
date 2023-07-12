package discovery

import (
	"context"

	pb "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
)

type Server struct {
	pb.UnimplementedDiscoveryServer
	Versions []int32
}

func (s *Server) GetVersions(context.Context, *pb.GetVersions_Request) (*pb.GetVersions_Response, error) {
	return &pb.GetVersions_Response{Versions: s.Versions}, nil
}
