package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type DestinationClient interface {
	Close(ctx context.Context) error
	Read(ctx context.Context, table *schema.Table, res chan<- arrow.RecordBatch) error
	Write(ctx context.Context, res <-chan message.WriteMessage) error
}

// writeOne is currently used mostly for testing, so it's not a public api
func (p *Plugin) writeOne(ctx context.Context, resource message.WriteMessage) error {
	resources := []message.WriteMessage{resource}
	return p.WriteAll(ctx, resources)
}

// WriteAll is currently used mostly for testing, so it's not a public api
func (p *Plugin) WriteAll(ctx context.Context, resources []message.WriteMessage) error {
	ch := make(chan message.WriteMessage, len(resources))
	for _, resource := range resources {
		ch <- resource
	}
	close(ch)
	return p.Write(ctx, ch)
}

func (p *Plugin) Write(ctx context.Context, res <-chan message.WriteMessage) error {
	if p.client == nil {
		return errors.New("plugin is not initialized. call Init first")
	}
	return p.client.Write(ctx, res)
}

// Read is read data from the requested table to the given channel, returned in the same format as the table
func (p *Plugin) Read(ctx context.Context, table *schema.Table, res chan<- arrow.RecordBatch) error {
	if !p.mu.TryLock() {
		return errors.New("plugin already in use")
	}
	defer p.mu.Unlock()
	if p.client == nil {
		return errors.New("plugin not initialized. call Init() first")
	}
	if err := p.client.Read(ctx, table, res); err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}
	return nil
}
