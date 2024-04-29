package state

import (
	"context"
	"fmt"
	"slices"

	pbDiscovery "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	pbPluginV3 "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	stateV3 "github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v3"
	"google.golang.org/grpc"
)

type Client interface {
	SetKey(ctx context.Context, key string, value string) error
	GetKey(ctx context.Context, key string) (string, error)
	Flush(ctx context.Context) error
	Close() error
}

type ClientOptions struct {
	Versioned bool
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
	if slices.Contains(versions.Versions, 3) {
		if opts.Versioned {
			return stateV3.NewClientWithTable(ctx, pbPluginV3.NewPluginClient(conn), stateV3.VersionedTable(tableName), conn)
		}
		return stateV3.NewClient(ctx, pbPluginV3.NewPluginClient(conn), tableName, conn)
	}
	return nil, fmt.Errorf("please upgrade your state backend plugin. state supporting version 3 plugin has %v", versions.Versions)
}

type NoOpClient struct{}

func (c *NoOpClient) SetKey(_ context.Context, _ string, _ string) error {
	return nil
}

func (c *NoOpClient) GetKey(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (c *NoOpClient) Flush(_ context.Context) error {
	return nil
}

func (c *NoOpClient) Close() error {
	return nil
}

// static check
var _ Client = (*NoOpClient)(nil)
