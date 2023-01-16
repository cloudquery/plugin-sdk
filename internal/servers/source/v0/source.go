package source

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	pbBase "github.com/cloudquery/plugin-sdk/internal/pb/base/v0"
	pb "github.com/cloudquery/plugin-sdk/internal/pb/source/v0"
	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	pb.UnimplementedSourceServer
	Plugin *source.Plugin
	Logger zerolog.Logger
}

// Deprecated: use GetSupportedProtocolVersions instead
func (*Server) GetProtocolVersion(context.Context, *pbBase.GetProtocolVersion_Request) (*pbBase.GetProtocolVersion_Response, error) {
	return &pbBase.GetProtocolVersion_Response{
		Version: 2,
	}, nil
}

func (s *Server) GetTables(context.Context, *pb.GetTables_Request) (*pb.GetTables_Response, error) {
	b, err := json.Marshal(s.Plugin.Tables())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tables: %w", err)
	}
	return &pb.GetTables_Response{
		Tables: b,
	}, nil
}

func (s *Server) GetTablesForSpec(_ context.Context, req *pb.GetTablesForSpec_Request) (*pb.GetTablesForSpec_Response, error) {
	if s.Plugin.HasDynamicTables() {
		return nil, errors.New("plugin has dynamic tables, please upgrade CLI version")
	}
	if len(req.Spec) == 0 {
		b, err := json.Marshal(s.Plugin.Tables())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tables: %w", err)
		}
		return &pb.GetTablesForSpec_Response{
			Tables: b,
		}, nil
	}

	var spec specs.Source
	dec := json.NewDecoder(bytes.NewReader(req.Spec))
	dec.UseNumber()
	if err := dec.Decode(&spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode spec: %v", err)
	}
	tables, err := s.Plugin.TablesForSpec(spec)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}
	b, err := json.Marshal(tables)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tables: %w", err)
	}
	return &pb.GetTablesForSpec_Response{
		Tables: b,
	}, nil
}

func (s *Server) GetName(context.Context, *pbBase.GetName_Request) (*pbBase.GetName_Response, error) {
	return &pbBase.GetName_Response{
		Name: s.Plugin.Name(),
	}, nil
}

func (s *Server) GetVersion(context.Context, *pbBase.GetVersion_Request) (*pbBase.GetVersion_Response, error) {
	return &pbBase.GetVersion_Response{
		Version: s.Plugin.Version(),
	}, nil
}

func (*Server) GetSyncSummary(context.Context, *pb.GetSyncSummary_Request) (*pb.GetSyncSummary_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSyncSummary is deprecated please upgrade client")
}

func (*Server) Sync(*pb.Sync_Request, pb.Source_SyncServer) error {
	return status.Errorf(codes.Unimplemented, "method Sync is deprecated please upgrade client")
}

func (s *Server) Sync2(req *pb.Sync2_Request, stream pb.Source_Sync2Server) error {
	if s.Plugin.HasDynamicTables() {
		return errors.New("plugin has dynamic tables, please upgrade CLI version")
	}
	resources := make(chan *schema.Resource)
	var syncErr error

	var spec specs.Source
	dec := json.NewDecoder(bytes.NewReader(req.Spec))
	dec.UseNumber()
	// TODO: warn about unknown fields
	if err := dec.Decode(&spec); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to decode spec: %v", err)
	}
	ctx := stream.Context()
	if err := s.Plugin.Init(ctx, spec); err != nil {
		return status.Errorf(codes.Internal, "failed to init plugin: %v", err)
	}
	defer func() {
		if err := s.Plugin.Close(ctx); err != nil {
			s.Logger.Error().Err(err).Msg("failed to close plugin")
		}
	}()

	go func() {
		defer close(resources)
		err := s.Plugin.Sync(ctx, resources)
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

		msg := &pb.Sync2_Response{
			Resource: b,
		}
		err = checkMessageSize(msg, resource)
		if err != nil {
			s.Logger.Warn().Str("table", resource.Table.Name).
				Int("bytes", len(msg.String())).
				Msg("Row exceeding max bytes ignored")
			continue
		}
		if err := stream.Send(msg); err != nil {
			return status.Errorf(codes.Internal, "failed to send resource: %v", err)
		}
	}

	return syncErr
}

func (s *Server) GetMetrics(context.Context, *pb.GetSourceMetrics_Request) (*pb.GetSourceMetrics_Response, error) {
	// Aggregate metrics before sending to keep response size small.
	// Temporary fix for https://github.com/cloudquery/cloudquery/issues/3962
	m := s.Plugin.Metrics()
	agg := &source.TableClientMetrics{}
	for _, table := range m.TableClient {
		for _, tableClient := range table {
			agg.Resources += tableClient.Resources
			agg.Errors += tableClient.Errors
			agg.Panics += tableClient.Panics
		}
	}
	b, err := json.Marshal(&source.Metrics{
		TableClient: map[string]map[string]*source.TableClientMetrics{"": {"": agg}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal source metrics: %w", err)
	}
	return &pb.GetSourceMetrics_Response{
		Metrics: b,
	}, nil
}

func checkMessageSize(msg proto.Message, resource *schema.Resource) error {
	size := proto.Size(msg)
	// log error to Sentry if row exceeds half of the max size
	if size > pb.MaxMsgSize/2 {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("table", resource.Table.Name)
			scope.SetExtra("bytes", size)
			sentry.CurrentHub().CaptureMessage("Large message detected")
		})
	}
	if size > pb.MaxMsgSize {
		return errors.New("message exceeds max size")
	}
	return nil
}
