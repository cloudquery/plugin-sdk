package proto

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/proto/internal"
)

type GRPCClient struct{ client internal.ProviderClient }

func (m *GRPCClient) Init(driver string, dsn string, verbose bool) error {
	_, err := m.client.Init(context.Background(), &internal.InitRequest{
		Driver:  driver,
		Dsn:     dsn,
		Verbose: verbose,
	})
	return err
}

func (m *GRPCClient) GenConfig() (string, error) {
	res, err := m.client.GenConfig(context.Background(), &internal.GenConfigRequest{})
	if err != nil {
		return "", err
	}
	return res.Yaml, nil
}

func (m *GRPCClient) Fetch(data []byte) error {
	_, err := m.client.Fetch(context.Background(), &internal.FetchRequest{
		Data: data,
	})
	return err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl CQProvider
	internal.UnimplementedProviderServer
}

func (m *GRPCServer) Init(ctx context.Context, req *internal.InitRequest) (*internal.InitResponse, error) {
	return &internal.InitResponse{}, m.Impl.Init(req.Driver, req.Dsn, req.Verbose)
}

func (m *GRPCServer) GenConfig(ctx context.Context, req *internal.GenConfigRequest) (*internal.GenConfigResponse, error) {
	r, err := m.Impl.GenConfig()
	if err != nil {
		return nil, err
	}
	return &internal.GenConfigResponse{Yaml: r}, nil
}

func (m *GRPCServer) Fetch(ctx context.Context, req *internal.FetchRequest) (*internal.FetchResponse, error) {
	err := m.Impl.Fetch(req.Data)
	return &internal.FetchResponse{}, err
}
