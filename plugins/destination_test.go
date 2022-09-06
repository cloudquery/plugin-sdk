package plugins

import (
	"context"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"testing"
)

func TestDestinationPlugin(t *testing.T) {
	f := func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
		return nil, nil
	}
	logger := zerolog.New(zerolog.NewTestWriter(t))
	giveExample := "hello"
	p := NewDestinationPlugin("test", "v1.2.3", f, WithDestinationLogger(logger), WithDestinationExampleConfig(giveExample))
	name := p.Name()
	if name != "test" {
		t.Errorf("plugin.Name() = %q, want %q", name, "test")
	}

	version := p.Version()
	if version != "v1.2.3" {
		t.Errorf("plugin.Version() = %q, want %q", version, "v1.2.3")
	}

	example := p.ExampleConfig()
	wantExample := "hello"
	if example != wantExample {
		t.Errorf("plugin.ExampleConfig() = %q, want %q", example, wantExample)
	}
}
