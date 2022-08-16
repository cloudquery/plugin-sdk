package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DestinationServer struct {
	pb.UnimplementedDestinationServer
	Plugin plugins.DestinationPlugin
}

func (s *DestinationServer) Configure(ctx context.Context, req *pb.Configure_Request) (*pb.Configure_Response, error) {
	var spec specs.Destination
	if err := json.Unmarshal(req.Config, &spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
	}
	return &pb.Configure_Response{}, s.Plugin.Initialize(ctx, spec)
}

func (s *DestinationServer) GetExampleConfig(ctx context.Context, req *pb.GetExampleConfig_Request) (*pb.GetExampleConfig_Response, error) {
	return &pb.GetExampleConfig_Response{
		Config: s.Plugin.ExampleConfig(),
	}, nil
}

func (s *DestinationServer) Write(msg pb.Destination_WriteServer) error {
	for {
		r, err := msg.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("write: failed to receive msg: %w", err)
		}
		var resource *schema.Resource
		if err := json.Unmarshal(r.Resource, &resource); err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
		}
		if err := s.Plugin.Write(context.Background(), resource); err != nil {
			return fmt.Errorf("write: failed to write resource: %w", err)
		}
	}
}
