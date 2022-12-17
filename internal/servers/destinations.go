package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DestinationServer struct {
	pb.UnimplementedDestinationServer
	Plugin *plugins.DestinationPlugin
	Logger zerolog.Logger
}

func (*DestinationServer) GetProtocolVersion(context.Context, *pb.GetProtocolVersion_Request) (*pb.GetProtocolVersion_Response, error) {
	return &pb.GetProtocolVersion_Response{
		Version: 2,
	}, nil
}

func (s *DestinationServer) Configure(ctx context.Context, req *pb.Configure_Request) (*pb.Configure_Response, error) {
	var spec specs.Destination
	if err := json.Unmarshal(req.Config, &spec); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal spec: %v", err)
	}
	return &pb.Configure_Response{}, s.Plugin.Init(ctx, s.Logger, spec)
}

func (s *DestinationServer) GetName(context.Context, *pb.GetName_Request) (*pb.GetName_Response, error) {
	return &pb.GetName_Response{
		Name: s.Plugin.Name(),
	}, nil
}

func (s *DestinationServer) GetVersion(context.Context, *pb.GetVersion_Request) (*pb.GetVersion_Response, error) {
	return &pb.GetVersion_Response{
		Version: s.Plugin.Version(),
	}, nil
}

func (s *DestinationServer) Migrate(ctx context.Context, req *pb.Migrate_Request) (*pb.Migrate_Response, error) {
	var tables []*schema.Table
	if err := json.Unmarshal(req.Tables, &tables); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	return &pb.Migrate_Response{}, s.Plugin.Migrate(ctx, tables)
}

func (*DestinationServer) Write(pb.Destination_WriteServer) error {
	return status.Errorf(codes.Unimplemented, "method Write is deprecated please upgrade client")
}

// Note the order of operations in this method is important!
// Trying to insert into the `resources` channel before starting the reader goroutine will cause a deadlock.
func (s *DestinationServer) Write2(msg pb.Destination_Write2Server) error {
	resources := make(chan schema.DestinationResource)

	r, err := msg.Recv()
	if err != nil {
		if err == io.EOF {
			return msg.SendAndClose(&pb.Write2_Response{})
		}
		return fmt.Errorf("write: failed to receive msg: %w", err)
	}
	var tables schema.Tables
	if err := json.Unmarshal(r.Tables, &tables); err != nil {
		return fmt.Errorf("write: failed to unmarshal tables: %w", err)
	}
	sourceName := r.Source
	syncTime := r.Timestamp.AsTime()

	eg, ctx := errgroup.WithContext(msg.Context())
	eg.Go(func() error {
		return s.Plugin.Write(ctx, tables, sourceName, syncTime, resources)
	})

	for {
		r, err := msg.Recv()
		if err != nil {
			close(resources)
			if err == io.EOF {
				if err := eg.Wait(); err != nil {
					return fmt.Errorf("got EOF. plugin returned: %w", err)
				}
				return msg.SendAndClose(&pb.Write2_Response{})
			}
			if pluginErr := eg.Wait(); pluginErr != nil {
				return fmt.Errorf("failed to receive msg: %v. plugin returned %w", err, pluginErr)
			}
			return fmt.Errorf("failed to receive msg: %w", err)
		}
		var resource schema.DestinationResource
		if err := json.Unmarshal(r.Resource, &resource); err != nil {
			close(resources)
			if err := eg.Wait(); err != nil {
				s.Logger.Error().Err(err).Msg("failed to unmarshal resource. failed to wait for plugin")
			}
			return status.Errorf(codes.InvalidArgument, "failed to unmarshal resource: %v", err)
		}
		select {
		case resources <- resource:
		case <-ctx.Done():
			close(resources)
			if err := eg.Wait(); err != nil {
				s.Logger.Error().Err(err).Msg("failed to wait")
			}
			return ctx.Err()
		}
	}
}

func (s *DestinationServer) GetMetrics(context.Context, *pb.GetDestinationMetrics_Request) (*pb.GetDestinationMetrics_Response, error) {
	stats := s.Plugin.Metrics()
	b, err := json.Marshal(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stats: %w", err)
	}
	return &pb.GetDestinationMetrics_Response{
		Metrics: b,
	}, nil
}

func (s *DestinationServer) DeleteStale(ctx context.Context, req *pb.DeleteStale_Request) (*pb.DeleteStale_Response, error) {
	var tables schema.Tables
	if err := json.Unmarshal(req.Tables, &tables); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	if err := s.Plugin.DeleteStale(ctx, tables, req.Source, req.Timestamp.AsTime()); err != nil {
		return nil, err
	}

	return &pb.DeleteStale_Response{}, nil
}

func (s *DestinationServer) Close(ctx context.Context, _ *pb.Close_Request) (*pb.Close_Response, error) {
	return &pb.Close_Response{}, s.Plugin.Close(ctx)
}
