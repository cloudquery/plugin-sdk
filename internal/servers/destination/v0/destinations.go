package destination

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/memory"
	pbBase "github.com/cloudquery/plugin-pb-go/pb/base/v0"
	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v0"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v2/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedDestinationServer
	Plugin *destination.Plugin
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
	return &pbBase.Configure_Response{}, s.Plugin.Init(ctx, s.Logger, spec)
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
	var tables schema.Tables
	if err := json.Unmarshal(req.Tables, &tables); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	SetDestinationManagedCqColumns(tables)
	s.setPKsForTables(tables)

	return &pb.Migrate_Response{}, s.Plugin.Migrate(ctx, tables.ToArrowSchemas())
}

func (*Server) Write(pb.Destination_WriteServer) error {
	return status.Errorf(codes.Unimplemented, "method Write is deprecated please upgrade client")
}

// Note the order of operations in this method is important!
// Trying to insert into the `resources` channel before starting the reader goroutine will cause a deadlock.
func (s *Server) Write2(msg pb.Destination_Write2Server) error {
	resources := make(chan arrow.Record)

	r, err := msg.Recv()
	if err != nil {
		if err == io.EOF {
			return msg.SendAndClose(&pb.Write2_Response{})
		}
		return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
	}
	var tables schema.Tables
	if err := json.Unmarshal(r.Tables, &tables); err != nil {
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

	syncTime := r.Timestamp.AsTime()
	SetDestinationManagedCqColumns(tables)
	s.setPKsForTables(tables)
	schemas := tables.ToArrowSchemas()
	eg, ctx := errgroup.WithContext(msg.Context())
	eg.Go(func() error {
		return s.Plugin.Write(ctx, sourceSpec, schemas, syncTime, resources)
	})
	sourceColumn := &schema.Text{}
	_ = sourceColumn.Set(sourceSpec.Name)
	syncTimeColumn := &schema.Timestamptz{}
	_ = syncTimeColumn.Set(syncTime)

	for {
		r, err := msg.Recv()
		if err == io.EOF {
			close(resources)
			if err := eg.Wait(); err != nil {
				return status.Errorf(codes.Internal, "write failed: %v", err)
			}
			return msg.SendAndClose(&pb.Write2_Response{})
		}
		if err != nil {
			close(resources)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.Internal, "failed to receive msg: %v and write failed: %v", err, wgErr)
			}
			return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
		}
		var origResource schema.DestinationResource
		if err := json.Unmarshal(r.Resource, &origResource); err != nil {
			close(resources)
			if wgErr := eg.Wait(); wgErr != nil {
				return status.Errorf(codes.InvalidArgument, "failed to unmarshal resource: %v and write failed: %v", err, wgErr)
			}
			return status.Errorf(codes.InvalidArgument, "failed to unmarshal resource: %v", err)
		}
		// this is a check to keep backward compatible for sources that are not adding
		// source and sync time
		if len(origResource.Data) < len(tables.Get(origResource.TableName).Columns) {
			origResource.Data = append([]schema.CQType{sourceColumn, syncTimeColumn}, origResource.Data...)
		}
		convertedResource := schema.CQTypesToRecord(memory.DefaultAllocator, []schema.CQTypes{origResource.Data}, schema.CQSchemaToArrow(tables.Get(origResource.TableName)))
		select {
		case resources <- convertedResource:
		case <-ctx.Done():
			convertedResource.Release()
			close(resources)
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
			table.Columns[i].CreationOptions.PrimaryKey = col.Name == schema.CqIDColumn.Name
		}
		setCQIDAsPrimaryKeysForTables(table.Relations)
	}
}

// Overwrites or adds the CQ columns that are managed by the destination plugins (_cq_sync_time, _cq_source_name).
func SetDestinationManagedCqColumns(tables []*schema.Table) {
	for _, table := range tables {
		for i := range table.Columns {
			if table.Columns[i].Name == schema.CqIDColumn.Name {
				table.Columns[i].CreationOptions.Unique = true
				table.Columns[i].CreationOptions.NotNull = true
			}
		}
		table.OverwriteOrAddColumn(&schema.CqSyncTimeColumn)
		table.OverwriteOrAddColumn(&schema.CqSourceNameColumn)
		SetDestinationManagedCqColumns(table.Relations)
	}
}

func (s *Server) GetMetrics(context.Context, *pb.GetDestinationMetrics_Request) (*pb.GetDestinationMetrics_Response, error) {
	stats := s.Plugin.Metrics()
	b, err := json.Marshal(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stats: %w", err)
	}
	return &pb.GetDestinationMetrics_Response{
		Metrics: b,
	}, nil
}

func (s *Server) DeleteStale(ctx context.Context, req *pb.DeleteStale_Request) (*pb.DeleteStale_Response, error) {
	var tables schema.Tables
	if err := json.Unmarshal(req.Tables, &tables); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	SetDestinationManagedCqColumns(tables)
	schemas := make(schema.Schemas, len(tables.FlattenTables()))
	for i, table := range tables.FlattenTables() {
		schemas[i] = table.ToArrowSchema()
	}
	if err := s.Plugin.DeleteStale(ctx, schemas, req.Source, req.Timestamp.AsTime()); err != nil {
		return nil, err
	}

	return &pb.DeleteStale_Response{}, nil
}

func (s *Server) setPKsForTables(tables schema.Tables) {
	if s.spec.PKMode == specs.PKModeCQID {
		setCQIDAsPrimaryKeysForTables(tables)
	}
}

func (s *Server) Close(ctx context.Context, _ *pb.Close_Request) (*pb.Close_Response, error) {
	return &pb.Close_Response{}, s.Plugin.Close(ctx)
}
