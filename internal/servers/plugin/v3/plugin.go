package plugin

import (
	"context"
	"fmt"
	"io"

	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
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
	schemas := tables.ToArrowSchemas()
	encoded := make([][]byte, len(schemas))
	for i, sc := range schemas {
		encoded[i], err = pb.SchemaToBytes(sc)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to encode tables: %v", err)
		}
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
	msgs := make(chan message.SyncMessage)
	var syncErr error
	ctx := stream.Context()

	syncOptions := plugin.SyncOptions{
		Tables:              req.Tables,
		SkipTables:          req.SkipTables,
		SkipDependentTables: req.SkipDependentTables,
		DeterministicCQID:   req.DeterministicCqId,
	}

	go func() {
		defer close(msgs)
		err := s.Plugin.Sync(ctx, syncOptions, msgs)
		if err != nil {
			syncErr = fmt.Errorf("failed to sync records: %w", err)
		}
	}()

	for msg := range msgs {
		pbMsg := &pb.Sync_Response{}
		switch m := msg.(type) {
		case *message.SyncMigrateTable:
			tableSchema := m.Table.ToArrowSchema()
			schemaBytes, err := pb.SchemaToBytes(tableSchema)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to encode table schema: %v", err)
			}
			pbMsg.Message = &pb.Sync_Response_MigrateTable{
				MigrateTable: &pb.Sync_MessageMigrateTable{
					Table: schemaBytes,
				},
			}

		case *message.SyncInsert:
			recordBytes, err := pb.RecordToBytes(m.Record)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to encode record: %v", err)
			}
			pbMsg.Message = &pb.Sync_Response_Insert{
				Insert: &pb.Sync_MessageInsert{
					Record: recordBytes,
				},
			}
		default:
			return status.Errorf(codes.Internal, "unknown message type: %T", msg)
		}

		size := proto.Size(pbMsg)
		if size > MaxMsgSize {
			s.Logger.Error().Int("bytes", size).Msg("Message exceeds max size")
			continue
		}
		if err := stream.Send(pbMsg); err != nil {
			return status.Errorf(codes.Internal, "failed to send message: %v", err)
		}
	}

	return syncErr
}

func (s *Server) Write(msg pb.Plugin_WriteServer) error {
	msgs := make(chan message.WriteMessage)
	eg, ctx := errgroup.WithContext(msg.Context())
	eg.Go(func() error {
		return s.Plugin.Write(ctx, msgs)
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
		var pluginMessage message.WriteMessage
		var pbMsgConvertErr error
		switch pbMsg := r.Message.(type) {
		case *pb.Write_Request_MigrateTable:
			sc, err := pb.NewSchemaFromBytes(pbMsg.MigrateTable.Table)
			if err != nil {
				pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create schema from bytes: %v", err)
				break
			}
			table, err := schema.NewTableFromArrowSchema(sc)
			if err != nil {
				pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create table from schema: %v", err)
				break
			}
			pluginMessage = &message.WriteMigrateTable{
				Table: table,
				MigrateForce: pbMsg.MigrateTable.MigrateForce,
			}
		case *pb.Write_Request_Insert:
			record, err := pb.NewRecordFromBytes(pbMsg.Insert.Record)
			if err != nil {
				pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create record: %v", err)
				break
			}
			pluginMessage = &message.WriteInsert{
				Record: record,
			}
		case *pb.Write_Request_Delete:
			sc, err := pb.NewSchemaFromBytes(pbMsg.Delete.Table)
			if err != nil {
				pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create schema from bytes: %v", err)
				break
			}
			table, err := schema.NewTableFromArrowSchema(sc)
			if err != nil {
				pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create table from schema: %v", err)
				break
			}
			pluginMessage = &message.WriteDeleteStale{
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

func (s *Server) Close(ctx context.Context, _ *pb.Close_Request) (*pb.Close_Response, error) {
	return &pb.Close_Response{}, s.Plugin.Close(ctx)
}
