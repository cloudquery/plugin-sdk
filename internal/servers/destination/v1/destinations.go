package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/ipc"
	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v1"
	"github.com/cloudquery/plugin-pb-go/specs/v0"
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
	writeMode   plugin.WriteMode
	migrateMode plugin.MigrateMode
}

func (s *Server) Configure(ctx context.Context, req *pb.Configure_Request) (*pb.Configure_Response, error) {
	var spec specs.Destination
	if err := json.Unmarshal(req.Config, &spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
	}
	s.spec = spec
	switch s.spec.WriteMode {
	case specs.WriteModeAppend:
		s.writeMode = plugin.WriteModeAppend
	case specs.WriteModeOverwrite:
		s.writeMode = plugin.WriteModeOverwrite
	case specs.WriteModeOverwriteDeleteStale:
		s.writeMode = plugin.WriteModeOverwriteDeleteStale
	}
	switch s.spec.MigrateMode {
	case specs.MigrateModeSafe:
		s.migrateMode = plugin.MigrateModeSafe
	case specs.MigrateModeForced:
		s.migrateMode = plugin.MigrateModeForced
	}
	return &pb.Configure_Response{}, s.Plugin.Init(ctx, s.spec.Spec)
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
	schemas, err := schema.NewSchemasFromBytes(req.Tables)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create schemas: %v", err)
	}
	tables, err := schema.NewTablesFromArrowSchemas(schemas)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create tables: %v", err)
	}
	s.setPKsForTables(tables)

	return &pb.Migrate_Response{}, s.Plugin.Migrate(ctx, tables, s.migrateMode)
}

// Note the order of operations in this method is important!
// Trying to insert into the `resources` channel before starting the reader goroutine will cause a deadlock.
func (s *Server) Write(msg pb.Destination_WriteServer) error {
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
	s.setPKsForTables(tables)
	eg, ctx := errgroup.WithContext(msg.Context())
	sourceName := r.Source

	eg.Go(func() error {
		return s.Plugin.Write(ctx, sourceName, tables, syncTime, s.writeMode, resources)
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

func setCQIDAsPrimaryKeysForTables(tables schema.Tables) {
	for _, table := range tables {
		for i, col := range table.Columns {
			table.Columns[i].PrimaryKey = col.Name == schema.CqIDColumn.Name
		}
		setCQIDAsPrimaryKeysForTables(table.Relations)
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
	schemas, err := schema.NewSchemasFromBytes(req.Tables)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create schemas: %v", err)
	}
	tables, err := schema.NewTablesFromArrowSchemas(schemas)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create tables: %v", err)
	}

	if err := s.Plugin.DeleteStale(ctx, tables, req.Source, req.Timestamp.AsTime()); err != nil {
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
