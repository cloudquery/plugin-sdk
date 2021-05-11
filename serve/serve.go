package serve

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudquery/cq-provider-sdk/cqproto"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

var Handshake = plugin.HandshakeConfig{
	MagicCookieKey:   "CQ_PLUGIN_COOKIE",
	MagicCookieValue: "6753812e-79c2-4af5-ad01-e6083c374e1f",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"provider": &cqproto.CQPlugin{},
}

type Options struct {
	// Required: Name of provider.
	Name string

	// Required: Provider is the actual provider that will be served.
	Provider cqproto.CQProviderServer

	// Optional: Logger is the logger that go-plugin will use.
	Logger hclog.Logger

	// Optional: Set NoLogOutputOverride to not override the log output with an hclog
	// adapter. This should only be used when running the plugin in
	// acceptance tests.
	NoLogOutputOverride bool

	// TestConfig should only be set when the provider is being tested; it
	// will opt out of go-plugin's lifecycle management and other features,
	// and will use the supplied configuration options to control the
	// plugin's lifecycle and communicate connection information. See the
	// go-plugin GoDoc for more information.
	TestConfig *plugin.ServeTestConfig
}

func Serve(opts *Options) {

	if opts.Name == "" {
		panic("missing provider name")
	}

	if opts.Provider == nil {
		panic("missing provider instance")
	}

	// Check of CQ_PROVIDER_DEBUG is turned on. In case it's true the plugin is executed in debug mode, allowing for
	// the CloudQuery main command to connect to this plugin via the .cq_reattach and the CQ_REATTACH_PROVIDERS env var
	if os.Getenv("CQ_PROVIDER_DEBUG") == "1" {
		if err := Debug(context.Background(), opts.Name, opts); err != nil {
			panic(fmt.Errorf("failed to run debug: %w", err))
		}
		return
	}
	serve(opts)
}

func serve(opts *Options) {

	if !opts.NoLogOutputOverride {
		// In order to allow go-plugin to correctly pass log-levels through to
		// cloudquery, we need to use an hclog.Logger with JSON output. We can
		// inject this into the std `log` package here, so existing providers will
		// make use of it automatically.
		logger := hclog.New(&hclog.LoggerOptions{
			// We send all output to CloudQuery. Go-plugin will take the output and
			// pass it through another hclog.Logger on the client side where it can
			// be filtered.
			Level:      hclog.Trace,
			JSONFormat: true,
		})
		log.SetOutput(logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true}))
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			2: {
				"provider": &cqproto.CQPlugin{Impl: opts.Provider},
			}},
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			return grpc.NewServer(opts...)
		},
		Logger: opts.Logger,
		Test:   opts.TestConfig,
	})
}
