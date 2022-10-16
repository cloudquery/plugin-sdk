package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/internal/versions"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
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
		Version: versions.DestinationProtocolVersion,
	}, nil
}

func (s *DestinationServer) GetStats(context.Context, *pb.GetDestinationStats_Request) (*pb.GetDestinationStats_Response, error) {
	stats := s.Plugin.Stats()
	b, err := json.Marshal(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stats: %w", err)
	}
	return &pb.GetDestinationStats_Response{
		Stats: b,
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

// Note the order of operations in this method is important!
// Trying to insert into the `resources` channel before starting the reader goroutine will cause a deadlock.
func (s *DestinationServer) Write(msg pb.Destination_WriteServer) error {
	resources := make(chan *schema.DestinationResource)

	r, err := msg.Recv()
	if err != nil {
		if err == io.EOF {
			return msg.SendAndClose(&pb.Write_Response{})
		}
		return fmt.Errorf("write: failed to receive msg: %w", err)
	}
	// tables := r.Tables
	source := r.Source
	timestamp := r.Timestamp.AsTime()
	var writeErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if writeErr = s.Plugin.Write(msg.Context(), nil, source, timestamp, resources); writeErr != nil {
			s.Logger.Error().Err(writeErr).Msg("write: failed to write resources")
		}
		// drain the channel
		for range resources {
			s.Logger.Error().Msg("skipping resource")
		}
	}()


	for {
		r, err := msg.Recv()
		if err != nil {
			if err == io.EOF {
				close(resources)
				wg.Wait()
				return msg.SendAndClose(&pb.Write_Response{})
			}
			close(resources)
			wg.Wait()
			return fmt.Errorf("write: failed to receive msg: %w", err)
		}
		// non recoverable write error
		if writeErr != nil {
			close(resources)
			wg.Wait()
			return writeErr
		}
		var resource *schema.DestinationResource
		if err := json.Unmarshal(r.Resource, &resource); err != nil {
			close(resources)
			wg.Wait()
			return status.Errorf(codes.InvalidArgument, "failed to unmarshal resource: %v", err)
		}
		select {
		case resources <- resource:
		case <-msg.Context().Done():
			close(resources)
			wg.Wait()
			return msg.Context().Err()
		}
	}
}

func (s *DestinationServer) DeleteStale(ctx context.Context, req *pb.DeleteStale_Request) (*pb.DeleteStale_Response, error) {
	var tables schema.Tables
	if err := json.Unmarshal(req.Tables, &tables); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	s.Plugin.DeleteStale(ctx, tables.TableNames(), req.Source, req.Timestamp.AsTime())

	return &pb.DeleteStale_Response{}, nil
}

func (s *DestinationServer) Close(ctx context.Context, _ *pb.Close_Request) (*pb.Close_Response, error) {
	return &pb.Close_Response{}, s.Plugin.Close(ctx)
}
