package plugin

import (
	"context"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type WriteOptions struct {
	MigrateForce bool
}

type DestinationClient interface {
	GetSpec() any
	Close(ctx context.Context) error
	Read(ctx context.Context, table *schema.Table, res chan<- arrow.Record) error
	Write(ctx context.Context, options WriteOptions, res <-chan message.Message) error
}

// writeOne is currently used mostly for testing, so it's not a public api
func (p *Plugin) writeOne(ctx context.Context, options WriteOptions, resource message.Message) error {
	resources := []message.Message{resource}
	return p.WriteAll(ctx, options, resources)
}

// WriteAll is currently used mostly for testing, so it's not a public api
func (p *Plugin) WriteAll(ctx context.Context, options WriteOptions, resources []message.Message) error {
	ch := make(chan message.Message, len(resources))
	for _, resource := range resources {
		ch <- resource
	}
	close(ch)
	return p.Write(ctx, options, ch)
}

func (p *Plugin) Write(ctx context.Context, options WriteOptions, res <-chan message.Message) error {
	if err := p.client.Write(ctx, options, res); err != nil {
		return err
	}
	return nil
}
