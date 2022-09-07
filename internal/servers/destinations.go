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
	Plugin *plugins.DestinationPlugin
}

func (s *DestinationServer) Configure(ctx context.Context, req *pb.Configure_Request) (*pb.Configure_Response, error) {
	var spec specs.Destination
	if err := json.Unmarshal(req.Config, &spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
	}
	return &pb.Configure_Response{}, s.Plugin.Initialize(ctx, spec)
}

func (s *DestinationServer) GetName(context.Context, *pb.GetName_Request) (*pb.GetName_Response, error) {
	return &pb.GetName_Response{
		Name: s.Plugin.Name(),
	}, nil
}

func (s *DestinationServer) GetVersion(context.Context, *pb.GetVersion_Request) (*pb.GetVersion_Response, error) {
	return &pb.GetVersion_Response{
		Version: s.Plugin.Version(),
	}, nil
}

func (s *DestinationServer) GetExampleConfig(_ context.Context, r *pb.GetDestinationExampleConfig_Request) (*pb.GetDestinationExampleConfig_Response, error) {
	registry, err := specs.RegistryFromString(r.GetRegistry())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid value for registry: %v", err)
	}
	cfg, err := s.Plugin.ExampleConfig(plugins.DestinationExampleConfigOptions{
		Path:     r.GetPath(),
		Registry: registry,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate destination config: %v", err)
	}
	return &pb.GetDestinationExampleConfig_Response{
		Config: cfg,
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
		if err := s.Plugin.Write(context.Background(), resource.TableName, resource.Data); err != nil {
			return fmt.Errorf("write: failed to write resource: %w", err)
		}
	}
}
