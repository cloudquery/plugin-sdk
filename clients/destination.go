package clients

import (
	"context"
	"fmt"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
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

func (c *DestinationClient) Write(ctx context.Context, resource *schema.Resource) error {
	// var saveClient pb.Destination_SaveClient
	// var err error
	// if c.pbClient != nil {
	// 	saveClient, err = c.pbClient.Write(ctx)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create save client: %w", err)
	// 	}
	// }
	if c.localClient != nil {
		if err := c.localClient.Write(ctx, resource); err != nil {
			return fmt.Errorf("failed to save resources: %w", err)
		}
	}

	return nil
}

// func (c *DestinationClient) CreateTables(ctx context.Context, tables []*schema.Table) error {
// 	if c.localClient != nil {
// 		return c.localClient.CreateTables(ctx, tables)
// 	}
// 	b, err := yaml.Marshal(tables)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal tables: %w", err)
// 	}
// 	if _, err := c.pbClient.CreateTables(ctx, &pb.CreateTables_Request{Tables: b}); err != nil {
// 		return err
// 	}
// 	return nil
// }
