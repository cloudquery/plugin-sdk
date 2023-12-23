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
	"github.com/cloudquery/plugin-sdk/v4/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"

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
	"github.com/getsentry/sentry-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type PluginServe struct {
	plugin                *plugin.Plugin
	args                  []string
	destinationV0V1Server bool
	sentryDSN             string
	testListener          bool
	testListenerConn      *bufconn.Listener
	versions              []int
}

type PluginOption func(*PluginServe)

func WithPluginSentryDSN(dsn string) PluginOption {
	return func(s *PluginServe) {
		s.sentryDSN = dsn
	}
}

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
	logLevel := newEnum([]string{"trace", "debug", "info", "warn", "error"}, "info")
	logFormat := newEnum([]string{"text", "json"}, "text")
	telemetryLevel := newEnum([]string{"none", "errors", "stats", "all"}, "all")
	err := telemetryLevel.Set(getEnvOrDefault("CQ_TELEMETRY_LEVEL", telemetryLevel.Value))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to set telemetry level: "+err.Error())
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: servePluginShort,
		Long:  servePluginShort,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			if otelEndpoint != "" {
				resource := newResource(s.plugin)
				opts := []otlptracehttp.Option{
					otlptracehttp.WithEndpoint(otelEndpoint),
				}
				if otelEndpointInsecure {
					opts = append(opts, otlptracehttp.WithInsecure())
				}
				if len(otelEndpointHeaders) > 0 {
					headers := make(map[string]string, len(otelEndpointHeaders))
					for _, h := range otelEndpointHeaders {
						parts := strings.SplitN(h, ":", 2)
						if len(parts) != 2 {
							return fmt.Errorf("invalid header %q", h)
						}
						headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
					}
					opts = append(opts, otlptracehttp.WithHeaders(headers))
				}
				if otelEndpointURLPath != "" {
					opts = append(opts, otlptracehttp.WithURLPath(otelEndpointURLPath))
				}

				client := otlptracehttp.NewClient(opts...)
				exp, err := otlptrace.New(cmd.Context(), client)
				if err != nil {
					return fmt.Errorf("creating OTLP trace exporter: %w", err)
				}
				tp := trace.NewTracerProvider(
					trace.WithBatcher(exp),
					trace.WithResource(resource),
				)
				defer func() {
					if err := tp.Shutdown(context.Background()); err != nil {
						logger.Error().Err(err).Msg("failed to shutdown OTLP trace exporter")
					}
				}()
				otel.SetTracerProvider(tp)
			}

			// opts.Plugin.Logger = logger
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
				Plugin:   s.plugin,
				Logger:   logger,
				NoSentry: noSentry,
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

			version := s.plugin.Version()

			if s.sentryDSN != "" && !strings.EqualFold(version, "development") && !noSentry {
				err = sentry.Init(sentry.ClientOptions{
					Dsn:              s.sentryDSN,
					Debug:            false,
					AttachStacktrace: false,
					Release:          version,
					Transport:        sentry.NewHTTPSyncTransport(),
					ServerName:       "oss", // set to "oss" on purpose to avoid sending any identifying information
					// https://docs.sentry.io/platforms/go/configuration/options/#removing-default-integrations
					Integrations: func(integrations []sentry.Integration) []sentry.Integration {
						var filteredIntegrations []sentry.Integration
						for _, integration := range integrations {
							if integration.Name() == "Modules" {
								continue
							}
							filteredIntegrations = append(filteredIntegrations, integration)
						}
						return filteredIntegrations
					},
				})
				if err != nil {
					log.Error().Err(err).Msg("Error initializing sentry")
				}
			}

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
	sendErrors := funk.ContainsString([]string{"all", "errors"}, telemetryLevel.String())
	if !sendErrors {
		noSentry = true
	}

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
