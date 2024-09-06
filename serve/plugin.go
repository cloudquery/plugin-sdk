package serve

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cloudquery/plugin-sdk/v4/helpers/grpczerolog"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/premium"
	"github.com/cloudquery/plugin-sdk/v4/types"

	pbDestinationV0 "github.com/cloudquery/plugin-pb-go/pb/destination/v0"
	pbDestinationV1 "github.com/cloudquery/plugin-pb-go/pb/destination/v1"
	pbdiscoveryv0 "github.com/cloudquery/plugin-pb-go/pb/discovery/v0"
	pbdiscoveryv1 "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	pbv3 "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	discoveryServerV0 "github.com/cloudquery/plugin-sdk/v4/internal/servers/discovery/v0"
	discoveryServerV1 "github.com/cloudquery/plugin-sdk/v4/internal/servers/discovery/v1"

	serverDestinationV0 "github.com/cloudquery/plugin-sdk/v4/internal/servers/destination/v0"
	serverDestinationV1 "github.com/cloudquery/plugin-sdk/v4/internal/servers/destination/v1"
	serversv3 "github.com/cloudquery/plugin-sdk/v4/internal/servers/plugin/v3"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type PluginServe struct {
	plugin                *plugin.Plugin
	args                  []string
	destinationV0V1Server bool
	testListener          bool
	testListenerConn      *bufconn.Listener
	versions              []int
}

type PluginOption func(*PluginServe)

// WithDestinationV0V1Server is used to include destination v0 and v1 server to work
// with older sources
func WithDestinationV0V1Server() PluginOption {
	return func(s *PluginServe) {
		s.destinationV0V1Server = true
	}
}

// WithArgs used to serve the plugin with predefined args instead of os.Args
func WithArgs(args ...string) PluginOption {
	return func(s *PluginServe) {
		s.args = args
	}
}

// WithTestListener means that the plugin will be served with an in-memory listener
// available via testListener() method instead of a network listener.
func WithTestListener() PluginOption {
	return func(s *PluginServe) {
		s.testListener = true
		s.testListenerConn = bufconn.Listen(testBufSize)
	}
}

const servePluginShort = `Start plugin server`

func Plugin(p *plugin.Plugin, opts ...PluginOption) *PluginServe {
	s := &PluginServe{
		plugin:   p,
		versions: []int{3},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *PluginServe) bufPluginDialer(context.Context, string) (net.Conn, error) {
	return s.testListenerConn.Dial()
}

func (s *PluginServe) Serve(ctx context.Context) error {
	if err := types.RegisterAllExtensions(); err != nil {
		return err
	}
	defer func() {
		if err := types.UnregisterAllExtensions(); err != nil {
			log.Error().Err(err).Msg("failed to unregister all extensions")
		}
	}()
	cmd := s.newCmdPluginRoot()
	if s.args != nil {
		cmd.SetArgs(s.args)
	}
	return cmd.ExecuteContext(ctx)
}

func (s *PluginServe) newCmdPluginServe() *cobra.Command {
	var address string
	var network string
	var noSentry bool
	var otelEndpoint string
	var otelEndpointHeaders []string
	var otelEndpointInsecure bool
	var otelEndpointURLPath string
	var licenseFile string
	logLevel := newEnum([]string{"trace", "debug", "info", "warn", "error"}, "info")
	logFormat := newEnum([]string{"text", "json"}, "text")
	telemetryLevel := newEnum([]string{"none", "errors", "stats", "all"}, "all")
	err := telemetryLevel.Set(getEnvOrDefault("CQ_TELEMETRY_LEVEL", telemetryLevel.Value))
	if err != nil {
		fmt.Fprint(os.Stderr, "failed to set telemetry level: "+err.Error())
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: servePluginShort,
		Long:  servePluginShort,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			zerologLevel, err := zerolog.ParseLevel(logLevel.String())
			if err != nil {
				return err
			}
			var logger zerolog.Logger
			if logFormat.String() == "json" {
				logger = zerolog.New(os.Stdout).Level(zerologLevel)
			} else {
				logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerologLevel)
			}

			shutdown, err := setupOtel(cmd.Context(), logger, s.plugin, otelEndpoint, otelEndpointInsecure, otelEndpointHeaders, otelEndpointURLPath)
			if err != nil {
				return fmt.Errorf("failed to setup OpenTelemetry: %w", err)
			}
			if shutdown != nil {
				logger = logger.Hook(newOTELLoggerHook())
				defer shutdown()
			}

			licenseClient, err := premium.NewLicenseClient(cmd.Context(), logger, premium.WithMeta(s.plugin.Meta()), premium.WithLicenseFileOrDirectory(licenseFile))
			if err != nil {
				return fmt.Errorf("failed to create license client: %w", err)
			}
			switch err := licenseClient.ValidateLicense(cmd.Context()); err {
			case nil:
				s.plugin.SetSkipUsageClient(true)
			case premium.ErrLicenseNotApplicable:
				// no-op: Treat as if no license was provided
			default:
				return fmt.Errorf("failed to validate license: %w", err)
			}

			var listener net.Listener
			if s.testListener {
				listener = s.testListenerConn
			} else {
				listener, err = net.Listen(network, address)
				if err != nil {
					return fmt.Errorf("failed to listen %s:%s: %w", network, address, err)
				}
			}
			defer listener.Close()
			// source plugins can only accept one connection at a time
			// unlike destination plugins that can accept multiple connections
			// limitListener := netutil.LimitListener(listener, 1)
			// See logging pattern https://github.com/grpc-ecosystem/go-grpc-middleware/blob/v2/providers/zerolog/examples_test.go
			grpcServer := grpc.NewServer(
				grpc.ChainUnaryInterceptor(
					logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger)),
				),
				grpc.ChainStreamInterceptor(
					logging.StreamServerInterceptor(grpczerolog.InterceptorLogger(logger)),
				),
				grpc.MaxRecvMsgSize(MaxMsgSize),
				grpc.MaxSendMsgSize(MaxMsgSize),
			)
			s.plugin.SetLogger(logger)
			pbv3.RegisterPluginServer(grpcServer, &serversv3.Server{
				Plugin: s.plugin,
				Logger: logger,
			})
			if s.destinationV0V1Server {
				pbDestinationV1.RegisterDestinationServer(grpcServer, &serverDestinationV1.Server{
					Plugin: s.plugin,
					Logger: logger,
				})
				pbDestinationV0.RegisterDestinationServer(grpcServer, &serverDestinationV0.Server{
					Plugin: s.plugin,
					Logger: logger,
				})
			}
			pbdiscoveryv0.RegisterDiscoveryServer(grpcServer, &discoveryServerV0.Server{
				Versions: []string{"v0", "v1", "v2", "v3"},
			})
			pbdiscoveryv1.RegisterDiscoveryServer(grpcServer, &discoveryServerV1.Server{
				Versions: []int32{0, 1, 2, 3},
			})

			ctx := cmd.Context()
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			defer func() {
				signal.Stop(c)
			}()

			go func() {
				select {
				case sig := <-c:
					logger.Info().Str("address", listener.Addr().String()).Str("signal", sig.String()).Msg("Got stop signal. Plugin server shutting down")
					grpcServer.Stop()
				case <-ctx.Done():
					logger.Info().Str("address", listener.Addr().String()).Msg("Context cancelled. Plugin server shutting down")
					grpcServer.Stop()
				}
			}()

			logger.Info().Str("address", listener.Addr().String()).Msg("Plugin server listening")
			if err := grpcServer.Serve(listener); err != nil {
				return fmt.Errorf("failed to serve: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&address, "address", "localhost:7777", "address to serve on. can be tcp: `localhost:7777` or unix socket: `/tmp/plugin.rpc.sock`")
	cmd.Flags().StringVar(&network, "network", "tcp", `the network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket"`)
	cmd.Flags().Var(logLevel, "log-level", fmt.Sprintf("log level. one of: %s", strings.Join(logLevel.Allowed, ",")))
	cmd.Flags().Var(logFormat, "log-format", fmt.Sprintf("log format. one of: %s", strings.Join(logFormat.Allowed, ",")))
	cmd.Flags().StringVar(&otelEndpoint, "otel-endpoint", "", "Open Telemetry HTTP collector endpoint")
	cmd.Flags().StringVar(&otelEndpointURLPath, "otel-endpoint-urlpath", "", "Open Telemetry HTTP collector endpoint URL path")
	cmd.Flags().StringArrayVar(&otelEndpointHeaders, "otel-endpoint-headers", []string{}, "Open Telemetry HTTP collector endpoint headers")
	cmd.Flags().BoolVar(&otelEndpointInsecure, "otel-endpoint-insecure", false, "use Open Telemetry HTTP endpoint (for development only)")
	cmd.Flags().BoolVar(&noSentry, "no-sentry", false, "disable sentry")
	cmd.Flags().StringVar(&licenseFile, "license", "", "Path to offline license file or directory")

	return cmd
}

func (s *PluginServe) newCmdPluginRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use: fmt.Sprintf("%s <command>", s.plugin.Name()),
	}
	cmd.AddCommand(s.newCmdPluginServe())
	cmd.AddCommand(s.newCmdPluginDoc())
	cmd.AddCommand(s.newCmdPluginPackage())
	cmd.AddCommand(s.newCmdPluginInfo())
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.Version = s.plugin.Version()
	return cmd
}
