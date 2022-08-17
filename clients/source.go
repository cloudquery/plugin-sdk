// package clients is a wrapper around grpc clients so clients can work
// with non protobuf structs and handle unmarshaling
package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"google.golang.org/grpc"
)

type SourceClient struct {
	pbClient pb.SourceClient
}

type FetchResultMessage struct {
	Resource []byte
}

func NewSourceClient(cc grpc.ClientConnInterface) *SourceClient {
	return &SourceClient{
		pbClient: pb.NewSourceClient(cc),
	}
}

func (c *SourceClient) GetTables(ctx context.Context) ([]*schema.Table, error) {
	res, err := c.pbClient.GetTables(ctx, &pb.GetTables_Request{})
	if err != nil {
		return nil, err
	}
	var tables []*schema.Table
	if err := json.Unmarshal(res.Tables, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

func (c *SourceClient) ExampleConfig(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetExampleConfig(ctx, &pb.GetExampleConfig_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to get example config: %w", err)
	}
	return res.Config, nil
}

func (c *SourceClient) Sync(ctx context.Context, spec specs.Source, res chan<- *schema.Resource) error {
	b, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal source spec: %w", err)
	}
	stream, err := c.pbClient.Sync(ctx, &pb.Sync_Request{
		Spec: b,
	})
	if err != nil {
		return fmt.Errorf("failed to sync resources: %w", err)
	}
	for {
		r, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to fetch resources from stream: %w", err)
		}
		var resource schema.Resource
		err = json.Unmarshal(r.Resource, &resource)
		if err != nil {
			return fmt.Errorf("failed to unmarshal resource: %w", err)
		}

		res <- &resource
	}
}
