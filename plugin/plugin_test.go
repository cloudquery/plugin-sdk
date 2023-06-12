package plugin

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
)

type testPluginSpec struct {
}

type testPluginClient struct {
	messages []Message
}

func newTestPluginClient(context.Context, zerolog.Logger, any) (Client, error) {
	return &testPluginClient{}, nil
}

func (c *testPluginClient) Tables(ctx context.Context) (schema.Tables, error) {
	return schema.Tables{}, nil
}
func (c *testPluginClient) Sync(ctx context.Context, options SyncOptions, res chan<- Message) error {
	for _, msg := range c.messages {
		res <- msg
	}
	return nil
}
func (c *testPluginClient) Write(ctx context.Context, options WriteOptions, res <-chan Message) error {
	for msg := range res {
		c.messages = append(c.messages, msg)
	}
	return nil
}
func (c *testPluginClient) Close(context.Context) error {
	return nil
}

func TestPluginSuccess(t *testing.T) {
	ctx := context.Background()
	p := NewPlugin("test", "v1.0.0", newTestPluginClient)
	if err := p.Init(ctx, &testPluginSpec{}); err != nil {
		t.Fatal(err)
	}
	tables, err := p.Tables(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(tables) != 0 {
		t.Fatal("expected 0 tables")
	}
	if err := p.WriteAll(ctx, WriteOptions{}, nil); err != nil {
		t.Fatal(err)
	}
	if err := p.WriteAll(ctx, WriteOptions{}, []Message{
		MessageCreateTable{},
	}); err != nil {
		t.Fatal(err)
	}
	if len(p.client.(*testPluginClient).messages) != 1 {
		t.Fatal("expected 1 message")
	}

	messages, err := p.SyncAll(ctx, SyncOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 1 {
		t.Fatal("expected 1 message")
	}

	if err := p.Close(ctx); err != nil {
		t.Fatal(err)
	}
}
