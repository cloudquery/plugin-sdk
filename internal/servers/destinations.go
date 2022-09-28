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

func (s *DestinationServer) Migrate(ctx context.Context, req *pb.Migrate_Request) (*pb.Migrate_Response, error) {
	var tables []*schema.Table
	if err := json.Unmarshal(req.Tables, &tables); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	return &pb.Migrate_Response{}, s.Plugin.Migrate(ctx, tables)
}

func (s *DestinationServer) Write(msg pb.Destination_WriteServer) error {
	failedWrites := uint64(0)
	for {
		r, err := msg.Recv()
		if err != nil {
			if err == io.EOF {
				return msg.SendAndClose(&pb.Write_Response{
					FailedWrites: failedWrites,
				})
			}
			return fmt.Errorf("write: failed to receive msg: %w", err)
		}
		var resource *schema.Resource
		if err := json.Unmarshal(r.Resource, &resource); err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
		}
		if err := s.Plugin.Write(msg.Context(), resource.TableName, resource.Data); err != nil {
			failedWrites++
		}
	}
}
