package servers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins/destination"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DestinationServer struct {
	pb.UnimplementedDestinationServer
	Plugin *destination.Plugin
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

func addErrorDetails(err error, reason string) error {
	st := status.Convert(err)
	st, _ = st.WithDetails(&errdetails.ErrorInfo{Reason: reason})
	return st.Err()
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
		return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
	}
	var tables schema.Tables
	if err := json.Unmarshal(r.Tables, &tables); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to unmarshal tables: %v", err)
	}
	sourceName := r.Source
	syncTime := r.Timestamp.AsTime()

	eg, ctx := errgroup.WithContext(msg.Context())
	eg.Go(func() error {
		return s.Plugin.Write(ctx, tables, sourceName, syncTime, resources)
	})

	for {
		r, err := msg.Recv()
		if err == io.EOF {
			close(resources)
			if err := eg.Wait(); err != nil {
				return addErrorDetails(errors.New("plugin write failed"), err.Error())
			}
			return msg.SendAndClose(&pb.Write2_Response{})
		}
		if err != nil {
			close(resources)
			if err := eg.Wait(); err != nil {
				return addErrorDetails(errors.New("plugin write failed"), err.Error())
			}
			return status.Errorf(codes.Internal, "failed to receive msg: %v", err)
		}
		var resource schema.DestinationResource
		if unmarshalError := json.Unmarshal(r.Resource, &resource); unmarshalError != nil {
			close(resources)
			errWithDetails := addErrorDetails(errors.New("failed to unmarshal resource"), unmarshalError.Error())
			if writeError := eg.Wait(); writeError != nil {
				errWithDetails = addErrorDetails(errors.New("plugin write failed and failed to unmarshal resource"), writeError.Error())
				errWithDetails = addErrorDetails(errWithDetails, unmarshalError.Error())
			}
			return errWithDetails
		}
		select {
		case resources <- resource:
		case <-ctx.Done():
			close(resources)
			if err := eg.Wait(); err != nil {
				return addErrorDetails(errors.New("plugin write failed"), err.Error())
			}
			ctxError := ctx.Err()
			if ctxError == nil {
				return nil
			}
			return addErrorDetails(errors.New("context error"), ctxError.Error())
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
