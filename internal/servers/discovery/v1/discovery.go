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
	v := make([]int32, len(s.Versions))
	for i := range s.Versions {
		v[i] = int32(s.Versions[i])
	}
	return &pb.GetVersions_Response{Versions: v}, nil
}
