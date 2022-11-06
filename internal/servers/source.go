package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudquery/plugin-sdk/v2/internal/pb"
	"github.com/cloudquery/plugin-sdk/v2/plugins"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SourceServer struct {
	pb.UnimplementedSourceServer
	Plugin *plugins.SourcePlugin
	Logger zerolog.Logger
}

func (*SourceServer) GetProtocolVersion(context.Context, *pb.GetProtocolVersion_Request) (*pb.GetProtocolVersion_Response, error) {
	return &pb.GetProtocolVersion_Response{
		Version: 2,
	}, nil
}

func (s *SourceServer) GetTables(context.Context, *pb.GetTables_Request) (*pb.GetTables_Response, error) {
	b, err := json.Marshal(s.Plugin.Tables())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tables: %w", err)
	}
	return &pb.GetTables_Response{
		Tables: b,
	}, nil
}

func (s *SourceServer) GetName(context.Context, *pb.GetName_Request) (*pb.GetName_Response, error) {
	return &pb.GetName_Response{
		Name: s.Plugin.Name(),
	}, nil
}

func (s *SourceServer) GetVersion(context.Context, *pb.GetVersion_Request) (*pb.GetVersion_Response, error) {
	return &pb.GetVersion_Response{
		Version: s.Plugin.Version(),
	}, nil
}

func (s *SourceServer) Sync(req *pb.Sync_Request, stream pb.Source_SyncServer) error {
	resources := make(chan *schema.Resource)
	var syncErr error

	var spec specs.Source
	dec := json.NewDecoder(bytes.NewReader(req.Spec))
	dec.UseNumber()
	// TODO: warn about unknown fields
	if err := dec.Decode(&spec); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to decode spec: %v", err)
	}

	go func() {
		defer close(resources)
		err := s.Plugin.Sync(stream.Context(), spec, resources)
		if err != nil {
			syncErr = fmt.Errorf("failed to sync resources: %w", err)
		}
	}()

	for resource := range resources {
		destResource := resource.ToDestinationResource()
		b, err := json.Marshal(destResource)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal resource: %v", err)
		}

		if err := stream.Send(&pb.Sync_Response{
			Resource: b,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to send resource: %v", err)
		}
	}

	return syncErr
}

func (s *SourceServer) GetMetrics(context.Context, *pb.GetSourceMetrics_Request) (*pb.GetSourceMetrics_Response, error) {
	b, err := json.Marshal(s.Plugin.Metrics())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal source metrics: %w", err)
	}
	return &pb.GetSourceMetrics_Response{
		Metrics: b,
	}, nil
}
