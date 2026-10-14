package state

import (
	"context"
	"fmt"
	"slices"

	pbDiscovery "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	"github.com/cloudquery/plugin-sdk/v4/internal/checkpoint"
	stateV3 "github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v3"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultMaxMsgSizeInBytes = 100 * 1024 * 1024 // 100 MiB

// ErrSyncAborted is returned by Flush when the sync was cancelled before the state could be sent as a checkpoint.
// The state stays pending.
var ErrSyncAborted = checkpoint.ErrSyncAborted

type Client interface {
	SetKey(ctx context.Context, key string, value string) error
	GetKey(ctx context.Context, key string) (string, error)
	Flush(ctx context.Context) error
	Close() error
}

type ClientOptions struct {
	// Checkpoints is the res channel of the plugin's Sync, the same one passed to the scheduler.
	// When the CLI supports checkpoints, Flush sends state on it instead of writing to the backend, and the CLI persists
	// the state once the destinations commit the resources synced before it. Changes made within a table resolver (using
	// its ctx) wait only for the resources that resolver produced; other changes must be flushed after the resources
	// they summarize are sent, e.g. once the scheduler's Sync returns.
	// Only NewConnectedClientWithOptions uses it, as only the backend options say whether the CLI supports checkpoints.
	Checkpoints chan<- message.SyncMessage
}

type ConnectionOptions struct {
	MaxMsgSizeInBytes int
}

func NewClient(ctx context.Context, conn *grpc.ClientConn, tableName string) (Client, error) {
	return NewClientWithOptions(ctx, conn, tableName, ClientOptions{})
}

func NewClientWithOptions(ctx context.Context, conn *grpc.ClientConn, tableName string, _ ClientOptions) (Client, error) {
	return newStateClient(ctx, conn, tableName, nil)
}

// newStateClient makes Flush send checkpoints on checkpoints, if set, instead of writing to the backend.
func newStateClient(ctx context.Context, conn *grpc.ClientConn, tableName string, checkpoints chan<- message.SyncMessage) (Client, error) {
	discoveryClient := pbDiscovery.NewDiscoveryClient(conn)
	versions, err := discoveryClient.GetVersions(ctx, &pbDiscovery.GetVersions_Request{})
	if err != nil {
		return nil, err
	}
	if !slices.Contains(versions.Versions, 3) {
		return nil, fmt.Errorf("please upgrade your state backend plugin. state supporting version 3 plugin has %v", versions.Versions)
	}

	var opts []stateV3.Option
	if checkpoints != nil {
		opts = append(opts, stateV3.WithCheckpoints(checkpoints))
	}
	return stateV3.NewClient(ctx, conn, tableName, opts...)
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
	// backend_options missing or configured but empty
	if backendOpts == nil || (backendOpts.Connection == "" && backendOpts.TableName == "") {
		return &NoOpClient{}, nil
	}

	// backend_options configured incorrectly, e.g. missing connection or table_name
	if backendOpts.Connection == "" || backendOpts.TableName == "" {
		return nil, fmt.Errorf("backend_options must contain both connection and table_name: %v", backendOpts)
	}

	var checkpoints chan<- message.SyncMessage
	if backendOpts.SupportsCheckpoints {
		checkpoints = clOpts.Checkpoints
	}
	if connOpts.MaxMsgSizeInBytes <= 0 {
		connOpts.MaxMsgSizeInBytes = defaultMaxMsgSizeInBytes
	}

	backendConn, err := grpc.NewClient(backendOpts.Connection,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(connOpts.MaxMsgSizeInBytes),
			grpc.MaxCallSendMsgSize(connOpts.MaxMsgSizeInBytes),
		),
	)

	if err != nil {
		return &NoOpClient{}, fmt.Errorf("failed to dial grpc source plugin at %s: %w", backendOpts.Connection, err)
	}

	stateClient, err := newStateClient(ctx, backendConn, backendOpts.TableName, checkpoints)
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
