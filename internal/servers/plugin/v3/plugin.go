package plugin

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-pb-go/managedplugin"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const MaxMsgSize = 100 * 1024 * 1024 // 100 MiB

type Server struct {
	pb.UnimplementedPluginServer
	Plugin    *plugin.Plugin
	Logger    zerolog.Logger
	Directory string
	NoSentry  bool
}

func (s *Server) GetTables(ctx context.Context, _ *pb.GetTables_Request) (*pb.GetTables_Response, error) {
	tables, err := s.Plugin.Tables(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tables: %v", err)
	}
	encoded, err := tables.ToArrowSchemas().Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode tables: %w", err)
	}
	return &pb.GetTables_Response{
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
	msgs := make(chan plugin.Message)
	var syncErr error
	ctx := stream.Context()

	syncOptions := plugin.SyncOptions{
		Tables:      req.Tables,
		SkipTables:  req.SkipTables,
		Concurrency: req.Concurrency,
	}

	if req.StateBackend != nil {
		opts := []managedplugin.Option{
			managedplugin.WithLogger(s.Logger),
			managedplugin.WithDirectory(s.Directory),
		}
		if s.NoSentry {
			opts = append(opts, managedplugin.WithNoSentry())
		}
		statePlugin, err := managedplugin.NewClient(ctx, managedplugin.Config{
			Path:     req.StateBackend.Path,
			Registry: managedplugin.Registry(req.StateBackend.Registry),
			Version:  req.StateBackend.Version,
		}, opts...)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to create state plugin: %v", err)
		}
		stateClient, err := newStateClient(ctx, statePlugin.Conn, req.StateBackend)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to create state client: %v", err)
		}
		syncOptions.StateBackend = stateClient
	}

	go func() {
		defer close(msgs)
		err := s.Plugin.Sync(ctx, syncOptions, msgs)
		if err != nil {
			syncErr = fmt.Errorf("failed to sync records: %w", err)
		}
	}()

	pbMsg := &pb.Sync_Response{}
	for msg := range msgs {
		switch m := msg.(type) {
		case *plugin.MessageCreateTable:
			m.Table.ToArrowSchema()
			pbMsg.Message = &pb.Sync_Response_CreateTable{
				CreateTable: &pb.MessageCreateTable{
					Table:        nil,
					MigrateForce: m.MigrateForce,
				},
			}
		case *plugin.MessageInsert:
			recordBytes, err := schema.RecordToBytes(m.Record)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to encode record: %v", err)
			}
			pbMsg.Message = &pb.Sync_Response_Insert{
				Insert: &pb.MessageInsert{
					Record: recordBytes,
					Upsert: m.Upsert,
				},
			}
		case *plugin.MessageDeleteStale:
			tableBytes, err := m.Table.ToArrowSchemaBytes()
			if err != nil {
				return status.Errorf(codes.Internal, "failed to encode record: %v", err)
			}
			pbMsg.Message = &pb.Sync_Response_Delete{
				Delete: &pb.MessageDeleteStale{
					Table:      tableBytes,
					SourceName: m.SourceName,
					SyncTime:   timestamppb.New(m.SyncTime),
				},
			}
		default:
			return status.Errorf(codes.Internal, "unknown message type: %T", msg)
		}

		// err := checkMessageSize(msg, rec)
		// if err != nil {
		// 	sc := rec.Schema()
		// 	tName, _ := sc.Metadata().GetValue(schema.MetadataTableName)
		// 	s.Logger.Warn().Str("table", tName).
		// 		Int("bytes", len(msg.String())).
		// 		Msg("Row exceeding max bytes ignored")
		// 	continue
		// }
		if err := stream.Send(pbMsg); err != nil {
			return status.Errorf(codes.Internal, "failed to send resource: %v", err)
		}
	}

	return syncErr
}

func (s *Server) Write(msg pb.Plugin_WriteServer) error {
	msgs := make(chan plugin.Message)

	eg, ctx := errgroup.WithContext(msg.Context())
	eg.Go(func() error {
		return s.Plugin.Write(ctx, plugin.WriteOptions{}, msgs)
	})

	for {
		r, err := msg.Recv()
		if err == io.EOF {
			close(msgs)
			if err := eg.Wait(); err != nil {
				return status.Errorf(codes.Internal, "write failed: %v", err)
			}
			return msg.SendAndClose(&pb.Write_Response{})
		}
		if err != nil {
			close(msgs)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.Internal, "failed to receive msg: %v and write failed: %v", err, wgErr)
			}
			return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
		}
		var pluginMessage plugin.Message
		var pbMsgConvertErr error
		switch pbMsg := r.Message.(type) {
		case *pb.Write_Request_CreateTable:
			table, err := schema.NewTableFromBytes(pbMsg.CreateTable.Table)
			if err != nil {
				pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create table: %v", err)
				break
			}
			pluginMessage = &plugin.MessageCreateTable{
				Table:        table,
				MigrateForce: pbMsg.CreateTable.MigrateForce,
			}
		case *pb.Write_Request_Insert:
			record, err := schema.NewRecordFromBytes(pbMsg.Insert.Record)
			if err != nil {
				pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create record: %v", err)
				break
			}
			pluginMessage = &plugin.MessageInsert{
				Record: record,
				Upsert: pbMsg.Insert.Upsert,
			}
		case *pb.Write_Request_Delete:
			table, err := schema.NewTableFromBytes(pbMsg.Delete.Table)
			if err != nil {
				pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create record: %v", err)
				break
			}
			pluginMessage = &plugin.MessageDeleteStale{
				Table:      table,
				SourceName: pbMsg.Delete.SourceName,
				SyncTime:   pbMsg.Delete.SyncTime.AsTime(),
			}
		}

		if pbMsgConvertErr != nil {
			close(msgs)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.Internal, "failed to convert message: %v and write failed: %v", pbMsgConvertErr, wgErr)
			}
			return pbMsgConvertErr
		}

		select {
		case msgs <- pluginMessage:
		case <-ctx.Done():
			close(msgs)
			if err := eg.Wait(); err != nil {
				return status.Errorf(codes.Internal, "Context done: %v and failed to wait for plugin: %v", ctx.Err(), err)
			}
			return status.Errorf(codes.Internal, "Context done: %v", ctx.Err())
		}
	}
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

func (s *Server) Close(ctx context.Context, _ *pb.Close_Request) (*pb.Close_Response, error) {
	return &pb.Close_Response{}, s.Plugin.Close(ctx)
}
