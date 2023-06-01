package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/ipc"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const MaxMsgSize = 100 * 1024 * 1024 // 100 MiB

type Server struct {
	pb.UnimplementedPluginServer
	Plugin *plugin.Plugin
	Logger zerolog.Logger
}

func (s *Server) GetStaticTables(context.Context, *pb.GetStaticTables_Request) (*pb.GetStaticTables_Response, error) {
	tables := s.Plugin.StaticTables().ToArrowSchemas()
	encoded, err := tables.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode tables: %w", err)
	}
	return &pb.GetStaticTables_Response{
		Tables: encoded,
	}, nil
}

func (s *Server) GetDynamicTables(context.Context, *pb.GetDynamicTables_Request) (*pb.GetDynamicTables_Response, error) {
	tables := s.Plugin.DynamicTables()
	if tables == nil {
		return &pb.GetDynamicTables_Response{}, nil
	}
	encoded, err := tables.ToArrowSchemas().Encode()
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
	if err := s.Plugin.Init(ctx, req.Spec); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to init plugin: %v", err)
	}
	return &pb.Init_Response{}, nil
}

func (s *Server) Sync(req *pb.Sync_Request, stream pb.Plugin_SyncServer) error {
	records := make(chan arrow.Record)
	var syncErr error
	ctx := stream.Context()

	syncOptions := plugin.SyncOptions{
		Tables:      req.Tables,
		SkipTables:  req.SkipTables,
		Concurrency: req.Concurrency,
		Scheduler:   plugin.SchedulerDFS,
	}
	if req.Scheduler == pb.SCHEDULER_SCHEDULER_ROUND_ROBIN {
		syncOptions.Scheduler = plugin.SchedulerRoundRobin
	}

	sourceName := req.SourceName

	go func() {
		defer close(records)
		err := s.Plugin.Sync(ctx, sourceName, req.SyncTime.AsTime(), syncOptions, records)
		if err != nil {
			syncErr = fmt.Errorf("failed to sync records: %w", err)
		}
	}()

	for rec := range records {
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
		err := checkMessageSize(msg, rec)
		if err != nil {
			sc := rec.Schema()
			tName, _ := sc.Metadata().GetValue(schema.MetadataTableName)
			s.Logger.Warn().Str("table", tName).
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
	agg := &plugin.TableClientMetrics{}
	for _, table := range m.TableClient {
		for _, tableClient := range table {
			agg.Resources += tableClient.Resources
			agg.Errors += tableClient.Errors
			agg.Panics += tableClient.Panics
		}
	}
	b, err := json.Marshal(&plugin.Metrics{
		TableClient: map[string]map[string]*plugin.TableClientMetrics{"": {"": agg}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal source metrics: %w", err)
	}
	return &pb.GetMetrics_Response{
		Metrics: b,
	}, nil
}

func (s *Server) Migrate(ctx context.Context, req *pb.Migrate_Request) (*pb.Migrate_Response, error) {
	schemas, err := schema.NewSchemasFromBytes(req.Tables)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create schemas: %v", err)
	}
	tables, err := schema.NewTablesFromArrowSchemas(schemas)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create tables: %v", err)
	}
	if req.PkMode == pb.PK_MODE_CQ_ID_ONLY {
		setCQIDAsPrimaryKeysForTables(tables)
	}
	migrateMode := plugin.MigrateModeSafe
	switch req.MigrateMode {
	case pb.MIGRATE_MODE_SAFE:
		migrateMode = plugin.MigrateModeSafe
	case pb.MIGRATE_MODE_FORCE:
		migrateMode = plugin.MigrateModeForced
	}
	// switch req.
	return &pb.Migrate_Response{}, s.Plugin.Migrate(ctx, tables, migrateMode)
}

func (s *Server) Write(msg pb.Plugin_WriteServer) error {
	resources := make(chan arrow.Record)

	r, err := msg.Recv()
	if err != nil {
		if err == io.EOF {
			return msg.SendAndClose(&pb.Write_Response{})
		}
		return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
	}

	schemas, err := schema.NewSchemasFromBytes(r.Tables)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to create schemas: %v", err)
	}
	tables, err := schema.NewTablesFromArrowSchemas(schemas)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to create tables: %v", err)
	}
	if r.PkMode == pb.PK_MODE_CQ_ID_ONLY {
		setCQIDAsPrimaryKeysForTables(tables)
	}
	sourceName := r.SourceName
	syncTime := r.SyncTime.AsTime()
	writeMode := plugin.WriteModeOverwrite
	switch r.WriteMode {
	case pb.WRITE_MODE_WRITE_MODE_APPEND:
		writeMode = plugin.WriteModeAppend
	case pb.WRITE_MODE_WRITE_MODE_OVERWRITE:
		writeMode = plugin.WriteModeOverwrite
	case pb.WRITE_MODE_WRITE_MODE_OVERWRITE_DELETE_STALE:
		writeMode = plugin.WriteModeOverwriteDeleteStale
	}
	eg, ctx := errgroup.WithContext(msg.Context())
	eg.Go(func() error {
		return s.Plugin.Write(ctx, sourceName, tables, syncTime, writeMode, resources)
	})

	for {
		r, err := msg.Recv()
		if err == io.EOF {
			close(resources)
			if err := eg.Wait(); err != nil {
				return status.Errorf(codes.Internal, "write failed: %v", err)
			}
			return msg.SendAndClose(&pb.Write_Response{})
		}
		if err != nil {
			close(resources)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.Internal, "failed to receive msg: %v and write failed: %v", err, wgErr)
			}
			return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
		}
		rdr, err := ipc.NewReader(bytes.NewReader(r.Resource))
		if err != nil {
			close(resources)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.InvalidArgument, "failed to create reader: %v and write failed: %v", err, wgErr)
			}
			return status.Errorf(codes.InvalidArgument, "failed to create reader: %v", err)
		}
		for rdr.Next() {
			rec := rdr.Record()
			rec.Retain()
			select {
			case resources <- rec:
			case <-ctx.Done():
				close(resources)
				if err := eg.Wait(); err != nil {
					return status.Errorf(codes.Internal, "Context done: %v and failed to wait for plugin: %v", ctx.Err(), err)
				}
				return status.Errorf(codes.Internal, "Context done: %v", ctx.Err())
			}
		}
		if err := rdr.Err(); err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to read resource: %v", err)
		}
	}
}

func (s *Server) GenDocs(req *pb.GenDocs_Request, srv pb.Plugin_GenDocsServer) error {
	tmpDir, err := os.MkdirTemp("", "cloudquery-docs")
	if err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	err = s.Plugin.GeneratePluginDocs(tmpDir, req.Format)
	if err != nil {
		return fmt.Errorf("failed to generate docs: %w", err)
	}

	// list files in tmpDir
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to read tmp dir: %w", err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(tmpDir, f.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		if err := srv.Send(&pb.GenDocs_Response{
			Filename: f.Name(),
			Content:  content,
		}); err != nil {
			return fmt.Errorf("failed to send file: %w", err)
		}
	}
	return nil
}

func checkMessageSize(msg proto.Message, record arrow.Record) error {
	size := proto.Size(msg)
	// log error to Sentry if row exceeds half of the max size
	if size > MaxMsgSize/2 {
		sc := record.Schema()
		tName, _ := sc.Metadata().GetValue(schema.MetadataTableName)
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("table", tName)
			scope.SetExtra("bytes", size)
			sentry.CurrentHub().CaptureMessage("Large message detected")
		})
	}
	if size > MaxMsgSize {
		return errors.New("message exceeds max size")
	}
	return nil
}

func setCQIDAsPrimaryKeysForTables(tables schema.Tables) {
	for _, table := range tables {
		for i, col := range table.Columns {
			table.Columns[i].PrimaryKey = col.Name == schema.CqIDColumn.Name
		}
		setCQIDAsPrimaryKeysForTables(table.Relations)
	}
}
