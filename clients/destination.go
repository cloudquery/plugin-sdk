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
	// this can be used if we have a plugin which is compiled in so we dont need to do any grpc requests
	localClient plugins.DestinationPlugin
}

func NewDestinationClient(cc grpc.ClientConnInterface) *DestinationClient {
	return &DestinationClient{
		pbClient: pb.NewDestinationClient(cc),
	}
}

func NewLocalDestinationClient(p plugins.DestinationPlugin) *DestinationClient {
	return &DestinationClient{
		localClient: p,
	}
}

func (c *DestinationClient) Name(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetExampleConfig(ctx, &pb.GetExampleConfig_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to get example config: %w", err)
	}
	return res.Name, nil
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

func (c *DestinationClient) Initialize(ctx context.Context, spec specs.Destination) error {
	if c.localClient != nil {
		return c.localClient.Initialize(ctx, spec)
	}
	b, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("destination configure: failed to marshal spec: %w", err)
	}
	_, err = c.pbClient.Configure(ctx, &pb.Configure_Request{
		Config: b,
	})
	if err != nil {
		return fmt.Errorf("destination configure: failed to configure: %w", err)
	}
	return nil
}

func (c *DestinationClient) Migrate(ctx context.Context, tables []*schema.Table) error {
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

func (c *DestinationClient) Write(ctx context.Context, table string, data map[string]interface{}) error {
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
