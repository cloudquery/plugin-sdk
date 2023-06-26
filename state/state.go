package state

import (
	"context"
	"fmt"

	pbDiscovery "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	pbPluginV3 "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	stateV3 "github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v3"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
)

type Client interface {
	SetKey(ctx context.Context, key string, value string) error
	GetKey(ctx context.Context, key string) (string, error)
}

func NewClient(ctx context.Context, conn *grpc.ClientConn, tableName string) (Client, error) {
	discoveryClient := pbDiscovery.NewDiscoveryClient(conn)
	versions, err := discoveryClient.GetVersions(ctx, &pbDiscovery.GetVersions_Request{})
	if err != nil {
		return nil, err
	}
	if slices.Contains(versions.Versions, 3) {
		return stateV3.NewClient(ctx, pbPluginV3.NewPluginClient(conn), tableName)
	}
	return nil, fmt.Errorf("please upgrade your state backend plugin. state supporting version 3 plugin has %v", versions.Versions)
}
