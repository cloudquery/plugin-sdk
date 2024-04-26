package state

import (
	"context"
	"fmt"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultMaxMsgSizeInBytes = 100 * 1024 * 1024 // 100 MiB

type ConnectionOptions struct {
	MaxMsgSizeInBytes int
}

// NewGrpcConnectedClient returns a state client and a gRPC connection to the state backend with a 100MiB max message size.
// The state client is guaranteed to be non-nil (it defaults to the NoOpClient).
func NewGrpcConnectedClient(ctx context.Context, backendOpts *plugin.BackendOptions) (Client, *grpc.ClientConn, error) {
	return NewGrpcConnectedClientWithOptions(ctx, backendOpts, ConnectionOptions{MaxMsgSizeInBytes: defaultMaxMsgSizeInBytes})
}

// NewGrpcConnectedClientWithOptions returns a state client and a gRPC connection to the state backend.
// The state client is guaranteed to be non-nil (it defaults to the NoOpClient).
func NewGrpcConnectedClientWithOptions(ctx context.Context, backendOpts *plugin.BackendOptions, opts ConnectionOptions) (Client, *grpc.ClientConn, error) {
	if backendOpts == nil {
		return &NoOpClient{}, nil, nil
	}

	backendConn, err := grpc.DialContext(ctx, backendOpts.Connection,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(opts.MaxMsgSizeInBytes),
			grpc.MaxCallSendMsgSize(opts.MaxMsgSizeInBytes),
		),
	)

	if err != nil {
		return &NoOpClient{}, nil, fmt.Errorf("failed to dial grpc source plugin at %s: %w", backendOpts.Connection, err)
	}

	stateClient, err := NewClient(ctx, backendConn, backendOpts.TableName)
	if err != nil {
		return &NoOpClient{}, nil, fmt.Errorf("failed to create state client: %w", err)
	}

	return stateClient, backendConn, nil
}
