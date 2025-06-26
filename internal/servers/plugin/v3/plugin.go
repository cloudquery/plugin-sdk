package plugin

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/apache/arrow-go/v18/arrow"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
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
}

func (s *Server) GetTables(ctx context.Context, req *pb.GetTables_Request) (*pb.GetTables_Response, error) {
	tables, err := s.Plugin.Tables(ctx, plugin.TableOptions{
		Tables:              req.Tables,
		SkipTables:          req.SkipTables,
		SkipDependentTables: req.SkipDependentTables,
	})
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

func (s *Server) GetSpecSchema(context.Context, *pb.GetSpecSchema_Request) (*pb.GetSpecSchema_Response, error) {
	sc := s.Plugin.JSONSchema()
	if len(sc) == 0 {
		return &pb.GetSpecSchema_Response{}, nil
	}

	return &pb.GetSpecSchema_Response{JsonSchema: &sc}, nil
}

func (s *Server) TestConnection(ctx context.Context, req *pb.TestConnection_Request) (*pb.TestConnection_Response, error) {
	err := s.Plugin.TestConnection(ctx, s.Logger, req.Spec)
	if err == nil {
		return &pb.TestConnection_Response{Success: true}, nil
	}

	const unknown = "UNKNOWN"
	var testConnErr *plugin.TestConnError
	if !errors.As(err, &testConnErr) {
		if errors.Is(err, plugin.ErrNotImplemented) {
			return nil, status.Errorf(codes.Unimplemented, "TestConnection feature is not implemented in this plugin")
		}

		return &pb.TestConnection_Response{
			Success:            false,
			FailureCode:        unknown,
			FailureDescription: err.Error(),
		}, nil
	}

	resp := &pb.TestConnection_Response{
		Success:     false,
		FailureCode: testConnErr.Code,
	}
	if resp.FailureCode == "" {
		resp.FailureCode = unknown
	}
	if testConnErr.Message != nil {
		resp.FailureDescription = testConnErr.Message.Error()
	}
	return resp, nil
}

func (s *Server) Init(ctx context.Context, req *pb.Init_Request) (*pb.Init_Response, error) {
	if err := s.Plugin.Init(ctx, req.Spec, plugin.NewClientOptions{NoConnection: req.NoConnection, InvocationID: req.InvocationId}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to init plugin: %v", err)
	}
	return &pb.Init_Response{}, nil
}

func (s *Server) Read(req *pb.Read_Request, stream pb.Plugin_ReadServer) error {
	records := make(chan arrow.Record)
	var readErr error
	ctx := stream.Context()

	sc, err := pb.NewSchemaFromBytes(req.Table)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to create schema from bytes: %v", err)
	}
	table, err := schema.NewTableFromArrowSchema(sc)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to create table from schema: %v", err)
	}
	go func() {
		defer close(records)
		err := s.Plugin.Read(ctx, table, records)
		if err != nil {
			readErr = fmt.Errorf("failed to read records: %w", err)
		}
	}()

	for rec := range records {
		recBytes, err := pb.RecordToBytes(rec)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to convert record to bytes: %v", err)
		}
		resp := &pb.Read_Response{
			Record: recBytes,
		}
		if err := stream.Send(resp); err != nil {
			return status.Errorf(codes.Internal, "failed to send read response: %v", err)
		}
	}

	return readErr
}

func flushMetrics() {
	traceProvider, ok := otel.GetTracerProvider().(*trace.TracerProvider)
	if ok && traceProvider != nil {
		traceProvider.ForceFlush(context.Background())
	}
	meterProvider, ok := otel.GetMeterProvider().(*metric.MeterProvider)
	if ok && meterProvider != nil {
		meterProvider.ForceFlush(context.Background())
	}
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
	if req.Backend != nil {
		syncOptions.BackendOptions = &plugin.BackendOptions{
			TableName:  req.Backend.TableName,
			Connection: req.Backend.Connection,
		}
	}
	if req.Shard != nil {
		syncOptions.Shard = &plugin.Shard{
			Num:   req.Shard.Num,
			Total: req.Shard.Total,
		}
	}

	go func() {
		defer flushMetrics()
		defer close(msgs)
		err := s.Plugin.Sync(ctx, syncOptions, msgs)
		if err != nil {
			syncErr = fmt.Errorf("failed to sync records: %w", err)
		}
	}()
	var err error
	for msg := range msgs {
		msg, err = s.Plugin.OnBeforeSend(ctx, msg)
		if err != nil {
			syncErr = fmt.Errorf("failed before sending message: %w", err)
			return syncErr
		}
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
		case *message.SyncDeleteRecord:
			whereClause := make([]*pb.PredicatesGroup, len(m.WhereClause))
			for j, predicateGroup := range m.WhereClause {
				whereClause[j] = &pb.PredicatesGroup{
					GroupingType: pb.PredicatesGroup_GroupingType(pb.PredicatesGroup_GroupingType_value[predicateGroup.GroupingType]),
					Predicates:   make([]*pb.Predicate, len(predicateGroup.Predicates)),
				}
				for i, predicate := range predicateGroup.Predicates {
					record, err := pb.RecordToBytes(predicate.Record)
					if err != nil {
						return status.Errorf(codes.Internal, "failed to encode record: %v", err)
					}

					whereClause[j].Predicates[i] = &pb.Predicate{
						Record:   record,
						Column:   predicate.Column,
						Operator: pb.Predicate_Operator(pb.Predicate_Operator_value[predicate.Operator]),
					}
				}
			}

			tableRelations := make([]*pb.TableRelation, len(m.TableRelations))
			for i, tr := range m.TableRelations {
				tableRelations[i] = &pb.TableRelation{
					TableName:   tr.TableName,
					ParentTable: tr.ParentTable,
				}
			}
			pbMsg.Message = &pb.Sync_Response_DeleteRecord{
				DeleteRecord: &pb.Sync_MessageDeleteRecord{
					TableName:      m.TableName,
					TableRelations: tableRelations,
					WhereClause:    whereClause,
				},
			}
		case *message.SyncError:
			if !req.WithErrorMessages {
				continue
			}
			pbMsg.Message = &pb.Sync_Response_Error{
				Error: &pb.Sync_MessageError{
					TableName: m.TableName,
					Error:     m.Error,
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

	if err := s.Plugin.OnSyncFinish(ctx); err != nil {
		return status.Errorf(codes.Internal, "failed to finish sync: %v", err)
	}

	return syncErr
}

func (s *Server) Write(stream pb.Plugin_WriteServer) error {
	msgs := make(chan message.WriteMessage)
	ctx := stream.Context()
	eg, gctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.Plugin.Write(gctx, msgs)
	})

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			close(msgs)
			if err := eg.Wait(); err != nil {
				return status.Errorf(codes.Internal, "write failed: %v", err)
			}
			return stream.SendAndClose(&pb.Write_Response{})
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
				Table:        table,
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
			pluginMessage = &message.WriteDeleteStale{
				TableName:  pbMsg.Delete.TableName,
				SourceName: pbMsg.Delete.SourceName,
				SyncTime:   pbMsg.Delete.SyncTime.AsTime(),
			}

		case *pb.Write_Request_DeleteRecord:
			whereClause := make(message.PredicateGroups, len(pbMsg.DeleteRecord.WhereClause))

			for j, predicateGroup := range pbMsg.DeleteRecord.WhereClause {
				whereClause[j].GroupingType = predicateGroup.GroupingType.String()
				whereClause[j].Predicates = make(message.Predicates, len(predicateGroup.Predicates))
				for i, predicate := range predicateGroup.Predicates {
					record, err := pb.NewRecordFromBytes(predicate.Record)
					if err != nil {
						pbMsgConvertErr = status.Errorf(codes.InvalidArgument, "failed to create record: %v", err)
						break
					}
					whereClause[j].Predicates[i] = message.Predicate{
						Record:   record,
						Column:   predicate.Column,
						Operator: predicate.Operator.String(),
					}
				}
			}

			tableRelations := make([]message.TableRelation, len(pbMsg.DeleteRecord.TableRelations))
			for i, tr := range pbMsg.DeleteRecord.TableRelations {
				tableRelations[i] = message.TableRelation{
					TableName:   tr.TableName,
					ParentTable: tr.ParentTable,
				}
			}
			pluginMessage = &message.WriteDeleteRecord{
				DeleteRecord: message.DeleteRecord{
					TableName:      pbMsg.DeleteRecord.TableName,
					TableRelations: tableRelations,
					WhereClause:    whereClause,
				},
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
		case <-gctx.Done():
			close(msgs)
			if err := eg.Wait(); err != nil {
				return status.Errorf(codes.Canceled, "plugin returned error: %v", err)
			}
			return status.Errorf(codes.Internal, "write failed for unknown reason")
		case <-ctx.Done():
			close(msgs)
			if err := eg.Wait(); err != nil {
				return status.Errorf(codes.Internal, "context done: %v and failed to wait for plugin: %v", ctx.Err(), err)
			}
			return status.Errorf(codes.Canceled, "context done: %v", ctx.Err())
		}
	}
}

func (s *Server) Transform(stream pb.Plugin_TransformServer) error {
	var (
		recvRecords = make(chan arrow.Record)
		sendRecords = make(chan arrow.Record)
		ctx         = stream.Context()
		eg, gctx    = errgroup.WithContext(ctx)
	)

	// Run the plugin's transform with both channels.
	//
	// When the plugin is done, it must return with either an error or nil.
	// The plugin must not close either channel.
	eg.Go(func() error {
		if err := s.Plugin.Transform(gctx, recvRecords, sendRecords); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		return nil
	})

	// Write transformed records from transformer to destination.
	//
	// Currently the `sendRecords` channel is never closed. Instead, the plugin finishes this goroutine
	// when it returns, either with an error or null.
	//
	// The reading never closes the writer, because it's up to the Plugin to decide when to finish
	// writing, regardless of if the reading finished.
	eg.Go(func() error {
		var sendErr error
		for record := range sendRecords {
			// We cannot terminate the stream here, because the plugin may still be sending records. So if error was returned channel has to be drained
			if sendErr != nil {
				continue
			}
			recordBytes, err := pb.RecordToBytes(record)
			if err != nil {
				sendErr = status.Errorf(codes.Internal, "failed to convert record to bytes: %v", err)
				continue
			}
			if err := stream.Send(&pb.Transform_Response{Record: recordBytes}); err != nil {
				sendErr = status.Errorf(codes.Internal, "error sending response: %v", err)
				continue
			}
		}
		return sendErr
	})

	// Read records from source to transformer
	//
	// If there's an error receiving or deserialising records, or if there are no more records,
	// the `recvRecords` channel will be closed. This will tell the plugin's transformer that
	// no more transforming can be done.
	//
	// The writer cannot stop the reader even on error, but the plugin will when it returns,
	// by setting `doneReading` to true.
	eg.Go(func() error {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				close(recvRecords)
				return nil
			}
			if err != nil {
				close(recvRecords)
				if status.Code(err) == codes.Canceled {
					// Ignore context cancellation errors
					return nil
				}
				return status.Errorf(codes.Internal, "Error receiving request: %v", err)
			}
			record, err := pb.NewRecordFromBytes(req.Record)
			if err != nil {
				close(recvRecords)
				return status.Errorf(codes.InvalidArgument, "failed to create record: %v", err)
			}

			select {
			case recvRecords <- record:
			case <-gctx.Done():
				close(recvRecords)
				return status.Errorf(codes.Canceled, "context done: %v", gctx.Err())
			case <-ctx.Done():
				close(recvRecords)
				return status.Errorf(codes.Canceled, "context done: %v", ctx.Err())
			}
		}
	})

	return eg.Wait()
}

func (s *Server) TransformSchema(ctx context.Context, req *pb.TransformSchema_Request) (*pb.TransformSchema_Response, error) {
	sc, err := pb.NewSchemaFromBytes(req.Schema)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create schema from bytes: %v", err)
	}
	newSchema, err := s.Plugin.TransformSchema(ctx, sc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transform schema: %v", err)
	}
	encoded, err := pb.SchemaToBytes(newSchema)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to encode schema: %v", err)
	}
	return &pb.TransformSchema_Response{Schema: encoded}, nil
}

func (s *Server) Close(ctx context.Context, _ *pb.Close_Request) (*pb.Close_Response, error) {
	return &pb.Close_Response{}, s.Plugin.Close(ctx)
}
