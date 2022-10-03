package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SourceServer struct {
	pb.UnimplementedSourceServer
	Plugin  *plugins.SourcePlugin
	Logger  zerolog.Logger
	summary *schema.SyncSummary
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

func (s *SourceServer) GetSyncSummary(context.Context, *pb.GetSyncSummary_Request) (*pb.GetSyncSummary_Response, error) {
	b, err := json.Marshal(s.summary)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal summary: %w", err)
	}
	return &pb.GetSyncSummary_Response{
		Summary: b,
	}, nil
}

func (s *SourceServer) Sync(req *pb.Sync_Request, stream pb.Source_SyncServer) error {
	resources := make(chan *schema.Resource)
	var syncErr error

	var spec specs.Source
	dec := json.NewDecoder(bytes.NewReader(req.Spec))
	dec.UseNumber()
	dec.DisallowUnknownFields()
	if err := dec.Decode(&spec); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to decode spec: %v", err)
	}

	go func() {
		defer close(resources)
		var err error
		s.summary, err = s.Plugin.Sync(stream.Context(), s.Logger, spec, resources)
		if err != nil {
			syncErr = fmt.Errorf("failed to sync resources: %w", err)
		}
	}()

	for resource := range resources {
		b, err := json.Marshal(resource)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal resource: %v", err)
		}
		if err := stream.Send(&pb.Sync_Response{
			Resource: b,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to send resource: %v", err)
		}
	}
	if syncErr != nil {
		return syncErr
	}

	return nil
}
