package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type SourceServer struct {
	pb.UnimplementedSourceServer
	Plugin *source.Plugin
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

func (s *SourceServer) GetTablesForSpec(_ context.Context, req *pb.GetTablesForSpec_Request) (*pb.GetTablesForSpec_Response, error) {
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

func (*SourceServer) GetSyncSummary(context.Context, *pb.GetSyncSummary_Request) (*pb.GetSyncSummary_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSyncSummary is deprecated please upgrade client")
}

func (*SourceServer) Sync(*pb.Sync_Request, pb.Source_SyncServer) error {
	return status.Errorf(codes.Unimplemented, "method Sync is deprecated please upgrade client")
}

func (s *SourceServer) Sync2(req *pb.Sync2_Request, stream pb.Source_Sync2Server) error {
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
			return status.Errorf(codes.InvalidArgument, "failed to marshal resource: %v", err)
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
			return fmt.Errorf("failed to send resource: %v", err)
		}
	}

	return syncErr
}

func (s *SourceServer) GetMetrics(context.Context, *pb.GetSourceMetrics_Request) (*pb.GetSourceMetrics_Response, error) {
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
	if size > MaxMsgSize/2 {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("table", resource.Table.Name)
			scope.SetExtra("bytes", size)
			sentry.CurrentHub().CaptureMessage("Large message detected")
		})
	}
	if size > MaxMsgSize {
		return errors.New("message exceeds max size")
	}
	return nil
}
