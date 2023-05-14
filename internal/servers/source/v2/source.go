package source

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/ipc"
	"github.com/apache/arrow/go/v13/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/source/v2"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v3/plugins/source"
	"github.com/cloudquery/plugin-sdk/v3/scalar"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const MaxMsgSize = 100 * 1024 * 1024 // 100 MiB

type Server struct {
	pb.UnimplementedSourceServer
	Plugin *source.Plugin
	Logger zerolog.Logger
}

func (s *Server) GetTables(context.Context, *pb.GetTables_Request) (*pb.GetTables_Response, error) {
	tables := s.Plugin.Tables().ToArrowSchemas()
	encoded, err := tables.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode tables: %w", err)
	}
	return &pb.GetTables_Response{
		Tables: encoded,
	}, nil
}

func (s *Server) GetDynamicTables(context.Context, *pb.GetDynamicTables_Request) (*pb.GetDynamicTables_Response, error) {
	tables := s.Plugin.GetDynamicTables().ToArrowSchemas()
	encoded, err := tables.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode tables: %w", err)
	}
	return &pb.GetDynamicTables_Response{
		Tables: encoded,
	}, nil
}

func (s *Server) GetName(context.Context, *pb.GetName_Request) (*pb.GetName_Response, error) {
	return &pb.GetName_Response{
		Name: s.Plugin.Name(),
	}, nil
}

func (s *Server) GetVersion(context.Context, *pb.GetVersion_Request) (*pb.GetVersion_Response, error) {
	return &pb.GetVersion_Response{
		Version: s.Plugin.Version(),
	}, nil
}

func (s *Server) Init(ctx context.Context, req *pb.Init_Request) (*pb.Init_Response, error) {
	var spec specs.Source
	dec := json.NewDecoder(bytes.NewReader(req.Spec))
	dec.UseNumber()
	// TODO: warn about unknown fields
	if err := dec.Decode(&spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode spec: %v", err)
	}

	if err := s.Plugin.Init(ctx, spec); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to init plugin: %v", err)
	}
	return &pb.Init_Response{}, nil
}

func (s *Server) Sync(_ *pb.Sync_Request, stream pb.Source_SyncServer) error {
	resources := make(chan *schema.Resource)
	var syncErr error
	ctx := stream.Context()

	go func() {
		defer close(resources)
		err := s.Plugin.Sync(ctx, resources)
		if err != nil {
			syncErr = fmt.Errorf("failed to sync resources: %w", err)
		}
	}()

	for resource := range resources {
		vector := resource.GetValues()
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, resource.Table.ToArrowSchema())
		scalar.AppendToRecordBuilder(bldr, vector)
		rec := bldr.NewRecord()

		var buf bytes.Buffer
		w := ipc.NewWriter(&buf, ipc.WithSchema(rec.Schema()))
		if err := w.Write(rec); err != nil {
			return status.Errorf(codes.Internal, "failed to write record: %v", err)
		}
		if err := w.Close(); err != nil {
			return status.Errorf(codes.Internal, "failed to close writer: %v", err)
		}

		msg := &pb.Sync_Response{
			Resource: buf.Bytes(),
		}
		err := checkMessageSize(msg, resource)
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

func (s *Server) GetMetrics(context.Context, *pb.GetMetrics_Request) (*pb.GetMetrics_Response, error) {
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
	return &pb.GetMetrics_Response{
		Metrics: b,
	}, nil
}

func (s *Server) GenDocs(_ context.Context, req *pb.GenDocs_Request) (*pb.GenDocs_Response, error) {
	err := s.Plugin.GeneratePluginDocs(req.Path, req.Format.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate docs: %w", err)
	}
	return &pb.GenDocs_Response{}, nil
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
