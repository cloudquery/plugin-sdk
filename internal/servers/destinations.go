package servers

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudquery/cq-provider-sdk/internal/pb"
	"github.com/cloudquery/cq-provider-sdk/plugins"
	"github.com/cloudquery/cq-provider-sdk/schema"
	"github.com/cloudquery/cq-provider-sdk/spec"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"
)

type DestinationServer struct {
	pb.UnimplementedDestinationServer
	Plugin plugins.DestinationPlugin
}

func (s *DestinationServer) Configure(ctx context.Context, req *pb.Configure_Request) (*pb.Configure_Response, error) {
	var spec spec.DestinationSpec
	if err := yaml.Unmarshal(req.Config, &spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
	}
	return &pb.Configure_Response{}, s.Plugin.Configure(ctx, spec)
}

func (s *DestinationServer) GetExampleConfig(ctx context.Context, req *pb.GetExampleConfig_Request) (*pb.GetExampleConfig_Response, error) {
	return &pb.GetExampleConfig_Response{
		Config: s.Plugin.GetExampleConfig(ctx),
	}, nil
}

func (s *DestinationServer) Save(msg pb.Destination_SaveServer) error {
	for {
		r, err := msg.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("Save: failed to receive msg: %w", err)
		}
		var resources []*schema.Resource
		if err := yaml.Unmarshal(r.Resources, &resources); err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
		}
		if err := s.Plugin.Save(context.Background(), resources); err != nil {
			return fmt.Errorf("Save: failed to save resources: %w", err)
		}
	}
}
