package destination

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	pbBase "github.com/cloudquery/plugin-pb-go/pb/base/v0"
	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v0"
	"github.com/cloudquery/plugin-pb-go/specs"
	schemav2 "github.com/cloudquery/plugin-sdk/v2/schema"
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
	Plugin *plugin.Plugin
	Logger zerolog.Logger
	spec   specs.Destination
}

func (*Server) GetProtocolVersion(context.Context, *pbBase.GetProtocolVersion_Request) (*pbBase.GetProtocolVersion_Response, error) {
	return &pbBase.GetProtocolVersion_Response{
		Version: 2,
	}, nil
}

func (s *Server) Configure(ctx context.Context, req *pbBase.Configure_Request) (*pbBase.Configure_Response, error) {
	var spec specs.Destination
	if err := json.Unmarshal(req.Config, &spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
	}
	s.spec = spec
	pluginSpec, err := json.Marshal(s.spec.Spec)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to marshal spec: %v", err)
	}
	return &pbBase.Configure_Response{}, s.Plugin.Init(ctx, pluginSpec, plugin.NewClientOptions{})
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

func (s *Server) Migrate(ctx context.Context, req *pb.Migrate_Request) (*pb.Migrate_Response, error) {
	var tablesV2 schemav2.Tables
	if err := json.Unmarshal(req.Tables, &tablesV2); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	tables := TablesV2ToV3(tablesV2).FlattenTables()
	SetDestinationManagedCqColumns(tables)
	s.setPKsForTables(tables)
	writeCh := make(chan message.WriteMessage)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.Plugin.Write(ctx, writeCh)
	})
	for _, table := range tables {
		writeCh <- &message.WriteMigrateTable{
			Table:        table,
			MigrateForce: s.spec.MigrateMode == specs.MigrateModeForced,
		}
	}
	close(writeCh)
	if err := eg.Wait(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write: %v", err)
	}
	return &pb.Migrate_Response{}, nil
}

func (*Server) Write(pb.Destination_WriteServer) error {
	return status.Errorf(codes.Unimplemented, "method Write is deprecated please upgrade client")
}

// Note the order of operations in this method is important!
// Trying to insert into the `resources` channel before starting the reader goroutine will cause a deadlock.
func (s *Server) Write2(msg pb.Destination_Write2Server) error {
	msgs := make(chan message.WriteMessage)

	r, err := msg.Recv()
	if err != nil {
		if err == io.EOF {
			return msg.SendAndClose(&pb.Write2_Response{})
		}
		return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
	}
	var tablesV2 schemav2.Tables
	if err := json.Unmarshal(r.Tables, &tablesV2); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
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
	tables := TablesV2ToV3(tablesV2).FlattenTables()
	syncTime := r.Timestamp.AsTime()
	SetDestinationManagedCqColumns(tables)
	s.setPKsForTables(tables)
	eg, ctx := errgroup.WithContext(msg.Context())
	// sourceName := r.Source
	eg.Go(func() error {
		return s.Plugin.Write(ctx, msgs)
	})

	for _, table := range tables {
		msgs <- &message.WriteMigrateTable{
			Table:        table,
			MigrateForce: s.spec.MigrateMode == specs.MigrateModeForced,
		}
	}

	sourceColumn := &schemav2.Text{}
	_ = sourceColumn.Set(sourceSpec.Name)
	syncTimeColumn := &schemav2.Timestamptz{}
	_ = syncTimeColumn.Set(syncTime)

	for {
		r, err := msg.Recv()
		if err == io.EOF {
			close(msgs)
			if err := eg.Wait(); err != nil {
				return status.Errorf(codes.Internal, "write failed: %v", err)
			}
			return msg.SendAndClose(&pb.Write2_Response{})
		}
		if err != nil {
			close(msgs)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.Internal, "failed to receive msg: %v and write failed: %v", err, wgErr)
			}
			return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
		}

		var origResource schemav2.DestinationResource
		if err := json.Unmarshal(r.Resource, &origResource); err != nil {
			close(msgs)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.InvalidArgument, "failed to unmarshal resource: %v and write failed: %v", err, wgErr)
			}
			return status.Errorf(codes.InvalidArgument, "failed to unmarshal resource: %v", err)
		}

		table := tables.Get(origResource.TableName)
		if table == nil {
			close(msgs)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.InvalidArgument, "failed to get table: %s and write failed: %v", origResource.TableName, wgErr)
			}
			return status.Errorf(codes.InvalidArgument, "failed to get table: %s", origResource.TableName)
		}
		// this is a check to keep backward compatible for sources that are not adding
		// source and sync time
		if len(origResource.Data) < len(table.Columns) {
			origResource.Data = append([]schemav2.CQType{sourceColumn, syncTimeColumn}, origResource.Data...)
		}
		convertedResource := CQTypesToRecord(memory.DefaultAllocator, []schemav2.CQTypes{origResource.Data}, table.ToArrowSchema())
		msg := &message.WriteInsert{
			Record: convertedResource,
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
}

func setCQIDAsPrimaryKeysForTables(tables schema.Tables) {
	for _, table := range tables {
		for i, col := range table.Columns {
			table.Columns[i].PrimaryKey = col.Name == schema.CqIDColumn.Name
		}
		setCQIDAsPrimaryKeysForTables(table.Relations)
	}
}

// Overwrites or adds the CQ columns that are managed by the destination plugins (_cq_sync_time, _cq_source_name).
func SetDestinationManagedCqColumns(tables []*schema.Table) {
	for _, table := range tables {
		for i := range table.Columns {
			if table.Columns[i].Name == schema.CqIDColumn.Name {
				table.Columns[i].Unique = true
				table.Columns[i].NotNull = true
			}
		}
		table.OverwriteOrAddColumn(&schema.CqSyncTimeColumn)
		table.OverwriteOrAddColumn(&schema.CqSourceNameColumn)
		SetDestinationManagedCqColumns(table.Relations)
	}
}

func (*Server) GetMetrics(context.Context, *pb.GetDestinationMetrics_Request) (*pb.GetDestinationMetrics_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetrics is deprecated. Please update CLI")
}

func (s *Server) DeleteStale(ctx context.Context, req *pb.DeleteStale_Request) (*pb.DeleteStale_Response, error) {
	var tablesV2 schemav2.Tables
	if err := json.Unmarshal(req.Tables, &tablesV2); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	tables := TablesV2ToV3(tablesV2).FlattenTables()
	SetDestinationManagedCqColumns(tables)

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
