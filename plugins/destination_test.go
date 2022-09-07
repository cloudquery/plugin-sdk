package plugins

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
)

type testDestinationClient struct {
	logger *zerolog.Logger
}

const wantDestinationConfig = `
kind: destination
spec:
  # Name of the plugin.
  name: "test"

  # Version of the plugin to use.
  version: "v1.2.3"

  # Registry to use (one of "github", "local" or "grpc").
  registry: "github"

  # Path to plugin. Required format depends on the registry.
  path: "cloudquery/test"

  # Write mode (either "overwrite" or "append").
  write_mode: "append"

  # Plugin-specific configuration.
  spec:
    test: hello
`

func (*testDestinationClient) Initialize(ctx context.Context, spec specs.Destination) error {
	return nil
}
func (*testDestinationClient) Migrate(ctx context.Context, tables schema.Tables) error {
	return nil
}
func (*testDestinationClient) Write(ctx context.Context, table string, data map[string]interface{}) error {
	return nil
}
func (t *testDestinationClient) SetLogger(logger zerolog.Logger) {
	t.logger = &logger
}

func TestDestinationPlugin(t *testing.T) {
	f := func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error) {
		return &testDestinationClient{}, nil
	}
	logger := zerolog.New(zerolog.NewTestWriter(t))
	giveExample := "test: hello"
	p := NewDestinationPlugin("test", "v1.2.3", f, WithDestinationLogger(logger), WithDestinationExampleConfig(giveExample))
	name := p.Name()
	if name != "test" {
		t.Errorf("plugin.Name() = %q, want %q", name, "test")
	}

	version := p.Version()
	if version != "v1.2.3" {
		t.Errorf("plugin.Version() = %q, want %q", version, "v1.2.3")
	}

	ctx := context.Background()
	spec := specs.Destination{
		Name:      "test",
		Version:   "v0.0.1",
		Path:      "cloudquery/test",
		Registry:  0,
		WriteMode: 0,
		Spec:      nil,
	}
	err := p.Initialize(ctx, spec)
	if err != nil {
		t.Fatalf("unexpected error calling Initialize: %v", err)
	}

	// check the generated config
	opts := DestinationExampleConfigOptions{
		Registry: specs.RegistryGithub,
	}
	example, err := p.ExampleConfig(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := strings.TrimSpace(wantDestinationConfig)
	if diff := cmp.Diff(example, want); diff != "" {
		t.Errorf("generated destination config not as expected (-got, +want): %v", diff)
	}
}
