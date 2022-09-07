package plugins

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

//go:embed templates/destination.go.tpl
var destinationFS embed.FS

type DestinationNewExecutionClientFunc func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error)

type DestinationClient interface {
	Initialize(ctx context.Context, spec specs.Destination) error
	Migrate(ctx context.Context, tables schema.Tables) error
	Write(ctx context.Context, table string, data map[string]interface{}) error
	SetLogger(logger zerolog.Logger)
}

// DestinationPlugin is the base structure required by calls to serve.Serve.
type DestinationPlugin struct {
	// name of plugin i.e aws, gcp, azure, etc
	name string
	// version of the plugin
	version string
	// example config for the plugin
	exampleConfig string
	// called on init to create a new DestinationClient
	newExecutionClient DestinationNewExecutionClientFunc
	// logger that should be used by the plugin
	logger zerolog.Logger
	// client returned by call to newExecutionClient
	client DestinationClient
	// configTemplate will be used to generate example config
	configTemplate *template.Template
}

type DestinationOption func(*DestinationPlugin)

// DestinationExampleConfigOptions can be used to override default example values.
type DestinationExampleConfigOptions struct {
	Path     string
	Registry specs.Registry
}

// WithDestinationExampleConfig sets an example config to user. It should only contain
// the inner "spec" part of a destination config. In other words, only the part specific to
// each destination plugin. The standard destination plugin config should not be included.
func WithDestinationExampleConfig(exampleConfig string) DestinationOption {
	return func(p *DestinationPlugin) {
		p.exampleConfig = exampleConfig
	}
}

func WithDestinationLogger(logger zerolog.Logger) DestinationOption {
	return func(p *DestinationPlugin) {
		p.logger = logger
	}
}

func (p *DestinationPlugin) Name() string {
	return p.name
}

func (p *DestinationPlugin) Version() string {
	return p.version
}

func (p *DestinationPlugin) Initialize(ctx context.Context, spec specs.Destination) error {
	c, err := p.newExecutionClient(ctx, p.logger, spec)
	if err != nil {
		return fmt.Errorf("failed to create execution client for destination plugin %s: %w", p.name, err)
	}
	return c.Initialize(ctx, spec)
}

// ExampleConfig returns a full example yaml config for the plugin.
func (p *DestinationPlugin) ExampleConfig(opts DestinationExampleConfigOptions) (string, error) {
	spec := specs.Destination{
		Name:     p.name,
		Version:  p.version,
		Path:     opts.Path,
		Registry: opts.Registry,
		Spec:     p.exampleConfig,
	}
	spec.SetDefaults()
	w := bytes.NewBufferString("")
	err := p.configTemplate.Execute(w, spec)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return w.String(), nil
}

func (p *DestinationPlugin) Migrate(ctx context.Context, tables schema.Tables) error {
	return p.client.Migrate(ctx, tables)
}

func (p *DestinationPlugin) Write(ctx context.Context, table string, data map[string]interface{}) error {
	return p.client.Write(ctx, table, data)
}

func (p *DestinationPlugin) SetLogger(logger zerolog.Logger) {
	p.logger = logger
	p.client.SetLogger(logger)
}

func NewDestinationPlugin(name string, version string, newExecutionClient DestinationNewExecutionClientFunc, opts ...DestinationOption) *DestinationPlugin {
	cfgTemplate := "destination.go.tpl"
	tpl, err := template.New(cfgTemplate).Funcs(template.FuncMap{
		"indent": indentSpaces,
	}).ParseFS(destinationFS, "templates/"+cfgTemplate)
	if err != nil {
		panic("failed to parse " + cfgTemplate + ":" + err.Error())
	}

	p := DestinationPlugin{
		name:               name,
		version:            version,
		exampleConfig:      "",
		logger:             zerolog.Logger{},
		newExecutionClient: newExecutionClient,
		configTemplate:     tpl,
	}
	for _, opt := range opts {
		opt(&p)
	}
	if err := p.validate(); err != nil {
		panic(err)
	}
	return &p
}

func (p *DestinationPlugin) validate() error {
	if p.newExecutionClient == nil {
		return fmt.Errorf("newExecutionClient function not defined for source plugin: " + p.name)
	}

	if p.name == "" {
		return errors.New("plugin name should not be empty")
	}

	return nil
}

func indentSpaces(text string, spaces int) string {
	s := strings.Repeat(" ", spaces)
	return s + strings.ReplaceAll(text, "\n", "\n"+s)
}
