package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v1"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedDestinationServer
	Plugin      *plugin.Plugin
	Logger      zerolog.Logger
	spec        specs.Destination
	migrateMode plugin.MigrateMode
}

func (s *Server) Configure(ctx context.Context, req *pb.Configure_Request) (*pb.Configure_Response, error) {
	var spec specs.Destination
	if err := json.Unmarshal(req.Config, &spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
	}
	s.spec = spec
	pluginSpec, err := json.Marshal(s.spec.Spec)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to marshal spec: %v", err)
	}
	var pluginSpecMap map[string]any
	if err := json.Unmarshal(pluginSpec, &pluginSpecMap); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal plugin spec: %v", err)
	}
	// this is for backward compatibility
	if s.spec.BatchSize > 0 {
		pluginSpecMap["batch_size"] = s.spec.BatchSize
	}
	if s.spec.BatchSizeBytes > 0 {
		pluginSpecMap["batch_size_bytes"] = s.spec.BatchSizeBytes
	}
	pluginSpec, err = json.Marshal(pluginSpecMap)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to marshal spec: %v", err)
	}
	return &pb.Configure_Response{}, s.Plugin.Init(ctx, pluginSpec, plugin.NewClientOptions{})
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

func (s *Server) Migrate(ctx context.Context, req *pb.Migrate_Request) (*pb.Migrate_Response, error) {
	schemas, err := NewSchemasFromBytes(req.Tables)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create schemas: %v", err)
	}
	tables, err := schema.NewTablesFromArrowSchemas(schemas)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create tables: %v", err)
	}
	s.setPKsForTables(tables)

	writeCh := make(chan message.WriteMessage)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.Plugin.Write(ctx, writeCh)
	})
	for _, table := range tables {
		writeCh <- &message.WriteMigrateTable{
			Table:        table,
			MigrateForce: s.migrateMode == plugin.MigrateModeForce,
		}
	}
	close(writeCh)
	if err := eg.Wait(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write: %v", err)
	}
	return &pb.Migrate_Response{}, nil
}

// Note the order of operations in this method is important!
// Trying to insert into the `resources` channel before starting the reader goroutine will cause a deadlock.
func (s *Server) Write(msg pb.Destination_WriteServer) error {
	msgs := make(chan message.WriteMessage)

	r, err := msg.Recv()
	if err != nil {
		if err == io.EOF {
			return msg.SendAndClose(&pb.Write_Response{})
		}
		return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
	}

	schemas, err := NewSchemasFromBytes(r.Tables)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to create schemas: %v", err)
	}
	tables, err := schema.NewTablesFromArrowSchemas(schemas)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to create tables: %v", err)
	}
	var sourceSpec specs.Source
	if r.SourceSpec == nil {
		// this is for backward compatibility
		sourceSpec = specs.Source{
			Name: r.Source,
		}
	} else {
		if err := json.Unmarshal(r.SourceSpec, &sourceSpec); err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to unmarshal source spec: %v", err)
		}
	}
	s.setPKsForTables(tables)
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
		rdr, err := ipc.NewReader(bytes.NewReader(r.Resource))
		if err != nil {
			close(msgs)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.InvalidArgument, "failed to create reader: %v and write failed: %v", err, wgErr)
			}
			return status.Errorf(codes.InvalidArgument, "failed to create reader: %v", err)
		}
		for rdr.Next() {
			rec := rdr.Record()
			rec.Retain()
			msg := &message.WriteInsert{
				Record: rec,
			}
			select {
			case msgs <- msg:
			case <-ctx.Done():
				close(msgs)
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

func setCQIDAsPrimaryKeysForTables(tables schema.Tables) {
	for _, table := range tables {
		for i, col := range table.Columns {
			table.Columns[i].PrimaryKey = col.Name == schema.CqIDColumn.Name
		}
		setCQIDAsPrimaryKeysForTables(table.Relations)
	}
}

func (*Server) GetMetrics(context.Context, *pb.GetDestinationMetrics_Request) (*pb.GetDestinationMetrics_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetrics is deprecated. please upgrade CLI")
}

func (s *Server) DeleteStale(ctx context.Context, req *pb.DeleteStale_Request) (*pb.DeleteStale_Response, error) {
	schemas, err := NewSchemasFromBytes(req.Tables)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create schemas: %v", err)
	}
	tables, err := schema.NewTablesFromArrowSchemas(schemas)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create tables: %v", err)
	}

	msgs := make(chan message.WriteMessage)
	var writeErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		writeErr = s.Plugin.Write(ctx, msgs)
	}()
	for _, table := range tables {
		bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
		bldr.Field(table.Columns.Index(schema.CqSourceNameColumn.Name)).(*array.StringBuilder).Append(req.Source)
		bldr.Field(table.Columns.Index(schema.CqSyncTimeColumn.Name)).(*array.TimestampBuilder).AppendTime(req.Timestamp.AsTime())
		msgs <- &message.WriteDeleteStale{
			TableName:  table.Name,
			SourceName: req.Source,
			SyncTime:   req.Timestamp.AsTime(),
		}
	}
	close(msgs)
	wg.Wait()
	return &pb.DeleteStale_Response{}, writeErr
}

func (s *Server) setPKsForTables(tables schema.Tables) {
	if s.spec.PKMode == specs.PKModeCQID {
		setCQIDAsPrimaryKeysForTables(tables)
	}
}

func (s *Server) Close(ctx context.Context, _ *pb.Close_Request) (*pb.Close_Response, error) {
	return &pb.Close_Response{}, s.Plugin.Close(ctx)
}
