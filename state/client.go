package state

import (
	"context"
	"fmt"
	"slices"

	pbDiscovery "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	stateV3 "github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v3"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultMaxMsgSizeInBytes = 100 * 1024 * 1024 // 100 MiB

type Client interface {
	SetKey(ctx context.Context, key string, value string) error
	GetKey(ctx context.Context, key string) (string, error)
	Flush(ctx context.Context) error
	Close() error
}

type ClientOptions struct {
	Versioned bool
}

type ConnectionOptions struct {
	MaxMsgSizeInBytes int
}

func NewClient(ctx context.Context, conn *grpc.ClientConn, tableName string) (Client, error) {
	return NewClientWithOptions(ctx, conn, tableName, ClientOptions{})
}

func NewClientWithOptions(ctx context.Context, conn *grpc.ClientConn, tableName string, opts ClientOptions) (Client, error) {
	discoveryClient := pbDiscovery.NewDiscoveryClient(conn)
	versions, err := discoveryClient.GetVersions(ctx, &pbDiscovery.GetVersions_Request{})
	if err != nil {
		return nil, err
	}
	if !slices.Contains(versions.Versions, 3) {
		return nil, fmt.Errorf("please upgrade your state backend plugin. state supporting version 3 plugin has %v", versions.Versions)
	}

	if opts.Versioned {
		return stateV3.NewClientWithTable(ctx, conn, stateV3.VersionedTable(tableName))
	}
	return stateV3.NewClient(ctx, conn, tableName)
}

// NewConnectedClient returns a state client and initialises the gRPC connection to the state backend with a 100MiB max message size.
// The state client is guaranteed to be non-nil (it defaults to the NoOpClient).
// You must call Close() on the returned Client object.
func NewConnectedClient(ctx context.Context, backendOpts *plugin.BackendOptions) (Client, error) {
	return NewConnectedClientWithOptions(ctx, backendOpts, ConnectionOptions{MaxMsgSizeInBytes: defaultMaxMsgSizeInBytes}, ClientOptions{})
}

// NewConnectedClientWithOptions returns a state client and initialises the gRPC connection to the state backend.
// The state client is guaranteed to be non-nil (it defaults to the NoOpClient).
// You must call Close() on the returned Client object.
func NewConnectedClientWithOptions(ctx context.Context, backendOpts *plugin.BackendOptions, connOpts ConnectionOptions, clOpts ClientOptions) (Client, error) {
	if backendOpts == nil {
		return &NoOpClient{}, nil
	}

	// TODO: Remove once there's a documented migration path per https://github.com/grpc/grpc-go/issues/7244
	// nolint:staticcheck
	backendConn, err := grpc.DialContext(ctx, backendOpts.Connection,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(connOpts.MaxMsgSizeInBytes),
			grpc.MaxCallSendMsgSize(connOpts.MaxMsgSizeInBytes),
		),
	)

	if err != nil {
		return &NoOpClient{}, fmt.Errorf("failed to dial grpc source plugin at %s: %w", backendOpts.Connection, err)
	}

	stateClient, err := NewClientWithOptions(ctx, backendConn, backendOpts.TableName, clOpts)
	if err != nil {
		return &NoOpClient{}, fmt.Errorf("failed to create state client: %w", err)
	}

	return stateClient, nil
}

type NoOpClient struct{}

func (*NoOpClient) SetKey(_ context.Context, _ string, _ string) error {
	return nil
}

func (*NoOpClient) GetKey(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (*NoOpClient) Flush(_ context.Context) error {
	return nil
}

func (*NoOpClient) Close() error {
	return nil
}

// static check
var _ Client = (*NoOpClient)(nil)
