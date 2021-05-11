package cqproto

import (
	"context"

	"github.com/cloudquery/cq-provider-sdk/cqproto/internal"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// CQPlugin This is the implementation of plugin.GRPCServer so we can serve/consume this.
type CQPlugin struct {
	// GRPCPlugin must still implement the Stub interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl CQProviderServer
}

func (p *CQPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	internal.RegisterProviderServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *CQPlugin) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{broker: broker, client: internal.NewProviderClient(c)}, nil
}
