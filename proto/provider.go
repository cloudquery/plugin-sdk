package proto

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/proto/internal"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type CQProvider interface {
	Init(driver string, dsn string, verbose bool) error
	Fetch(data []byte) error
	GenConfig() (string, error)
}

// CQPlugin this is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type CQPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl CQProvider
}

func (p *CQPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	internal.RegisterProviderServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *CQPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: internal.NewProviderClient(c)}, nil
}
