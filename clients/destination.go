package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"google.golang.org/grpc"
)

type DestinationClient struct {
	pbClient pb.DestinationClient
	// this can be used if we have a plugin which is compiled in, so we don't need to do any grpc requests
	localClient *plugins.DestinationPlugin
}

func NewDestinationClient(cc grpc.ClientConnInterface) *DestinationClient {
	return &DestinationClient{
		pbClient: pb.NewDestinationClient(cc),
	}
}

func NewLocalDestinationClient(p *plugins.DestinationPlugin) *DestinationClient {
	return &DestinationClient{
		localClient: p,
	}
}

func (c *DestinationClient) Name(ctx context.Context) (string, error) {
	if c.localClient != nil {
		return c.localClient.Name(), nil
	}
	res, err := c.pbClient.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to get name: %w", err)
	}
	return res.Name, nil
}

func (c *DestinationClient) Version(ctx context.Context) (string, error) {
	if c.localClient != nil {
		return c.localClient.Version(), nil
	}
	res, err := c.pbClient.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}
	return res.Version, nil
}

func (c *DestinationClient) GetExampleConfig(ctx context.Context) (string, error) {
	if c.localClient != nil {
		return c.localClient.ExampleConfig(), nil
	}
	res, err := c.pbClient.GetExampleConfig(ctx, &pb.GetExampleConfig_Request{})
	if err != nil {
		return "", err
	}
	return res.Config, nil
}

func (c *DestinationClient) Migrate(ctx context.Context, spec specs.Destination, tables []*schema.Table) error {
	if c.localClient != nil {
		return c.localClient.Migrate(ctx, tables)
	}
	b, err := json.Marshal(tables)
	if err != nil {
		return fmt.Errorf("destination migrate: failed to marshal plugin: %w", err)
	}
	_, err = c.pbClient.Migrate(ctx, &pb.Migrate_Request{Tables: b})
	if err != nil {
		return fmt.Errorf("destination migrate: failed to migrate: %w", err)
	}
	return nil
}

func (c *DestinationClient) Write(ctx context.Context, spec specs.Destination, table string, data map[string]interface{}) error {
	// var saveClient pb.Destination_SaveClient
	// var err error
	// if c.pbClient != nil {
	// 	saveClient, err = c.pbClient.Write(ctx)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create save client: %w", err)
	// 	}
	// }
	if c.localClient != nil {
		if err := c.localClient.Write(ctx, table, data); err != nil {
			return fmt.Errorf("failed to save resources: %w", err)
		}
	}

	return nil
}
