package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"google.golang.org/grpc"
)

type DestinationClient struct {
	pbClient pb.DestinationClient
}

func NewDestinationClient(cc grpc.ClientConnInterface) *DestinationClient {
	return &DestinationClient{
		pbClient: pb.NewDestinationClient(cc),
	}
}

func (c *DestinationClient) Name(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to get name: %w", err)
	}
	return res.Name, nil
}

func (c *DestinationClient) Version(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}
	return res.Version, nil
}

func (c *DestinationClient) Initialize(ctx context.Context, spec specs.Destination) error {
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

// Write writes rows as they are received from the channel to the destination plugin.
// resources is marshaled schema.Resource. We are not marshalling this inside the function
// because usually it is alreadun marshalled from the source plugin.
func (c *DestinationClient) Write(ctx context.Context, resources <-chan []byte) (uint64, error) {
	saveClient, err := c.pbClient.Write(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to create save client: %w", err)
	}
	var failedWrites uint64
	for resource := range resources {
		if err := saveClient.Send(&pb.Write_Request{
			Resource: resource,
		}); err != nil {
			failedWrites++
		}
	}

	return failedWrites, nil
}
