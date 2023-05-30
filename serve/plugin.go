package serve

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/cloudquery/plugin-sdk/v4/plugin"

	pbDestinationV0 "github.com/cloudquery/plugin-pb-go/pb/destination/v0"
	pbDestinationV1 "github.com/cloudquery/plugin-pb-go/pb/destination/v1"
	pbdiscoveryv0 "github.com/cloudquery/plugin-pb-go/pb/discovery/v0"
	pbv3 "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	discoveryServerV0 "github.com/cloudquery/plugin-sdk/v4/internal/servers/discovery/v0"

	serverDestinationV0 "github.com/cloudquery/plugin-sdk/v4/internal/servers/destination/v0"
	serverDestinationV1 "github.com/cloudquery/plugin-sdk/v4/internal/servers/destination/v1"
	serversv3 "github.com/cloudquery/plugin-sdk/v4/internal/servers/plugin/v3"
	"github.com/getsentry/sentry-go"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type pluginServe struct {
	plugin    *plugin.Plugin
	destinationV0V1Server bool
	sentryDSN string
}

type PluginOption func(*pluginServe)

func WithPluginSentryDSN(dsn string) PluginOption {
	return func(s *pluginServe) {
		s.sentryDSN = dsn
	}
}

// WithDestinationV0V1Server is used to include destination v0 and v1 server to work
// with older sources
func WithDestinationV0V1Server() PluginOption {
	return func(s *pluginServe) {
		s.destinationV0V1Server = true
	}
}

// lis used for unit testing grpc server and client
var testPluginListener *bufconn.Listener
var testPluginListenerLock sync.Mutex

const servePluginShort = `Start plugin server`

func Plugin(plugin *plugin.Plugin, opts ...PluginOption) {
	s := &pluginServe{
		plugin: plugin,
	}
	for _, opt := range opts {
		opt(s)
	}
	if err := newCmdPluginRoot(s).Execute(); err != nil {
		sentry.CaptureMessage(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}

// nolint:dupl
func newCmdPluginServe(serve *pluginServe) *cobra.Command {
	var address string
	var network string
	var noSentry bool
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

			// opts.Plugin.Logger = logger
			var listener net.Listener
			if network == "test" {
				testPluginListenerLock.Lock()
				listener = bufconn.Listen(testBufSize)
				testPluginListener = listener.(*bufconn.Listener)
				testPluginListenerLock.Unlock()
			} else {
				listener, err = net.Listen(network, address)
				if err != nil {
					return fmt.Errorf("failed to listen %s:%s: %w", network, address, err)
				}
			}
			// source plugins can only accept one connection at a time
			// unlike destination plugins that can accept multiple connections
			limitListener := netutil.LimitListener(listener, 1)
			// See logging pattern https://github.com/grpc-ecosystem/go-grpc-middleware/blob/v2/providers/zerolog/examples_test.go
			s := grpc.NewServer(
				grpc.ChainUnaryInterceptor(
					logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger)),
				),
				grpc.ChainStreamInterceptor(
					logging.StreamServerInterceptor(grpczerolog.InterceptorLogger(logger)),
				),
				grpc.MaxRecvMsgSize(MaxMsgSize),
				grpc.MaxSendMsgSize(MaxMsgSize),
			)
			serve.plugin.SetLogger(logger)
			pbv3.RegisterPluginServer(s, &serversv3.Server{
				Plugin: serve.plugin,
				Logger: logger,
			})
			if serve.destinationV0V1Server {
				pbDestinationV1.RegisterDestinationServer(s, &serverDestinationV1.Server{
					Plugin: serve.plugin,
					Logger: logger,
				})
				pbDestinationV0.RegisterDestinationServer(s, &serverDestinationV0.Server{
					Plugin: serve.plugin,
					Logger: logger,
				})
			}
			pbdiscoveryv0.RegisterDiscoveryServer(s, &discoveryServerV0.Server{
				Versions: []string{"v0", "v1", "v2", "v3"},
			})

			version := serve.plugin.Version()

			if serve.sentryDSN != "" && !strings.EqualFold(version, "development") && !noSentry {
				err = sentry.Init(sentry.ClientOptions{
					Dsn:              serve.sentryDSN,
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
					logger.Info().Str("address", listener.Addr().String()).Str("signal", sig.String()).Msg("Got stop signal. Source plugin server shutting down")
					s.Stop()
				case <-ctx.Done():
					logger.Info().Str("address", listener.Addr().String()).Msg("Context cancelled. Source plugin server shutting down")
					s.Stop()
				}
			}()

			logger.Info().Str("address", listener.Addr().String()).Msg("Source plugin server listening")
			if err := s.Serve(limitListener); err != nil {
				return fmt.Errorf("failed to serve: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&address, "address", "localhost:7777", "address to serve on. can be tcp: `localhost:7777` or unix socket: `/tmp/plugin.rpc.sock`")
	cmd.Flags().StringVar(&network, "network", "tcp", `the network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket"`)
	cmd.Flags().Var(logLevel, "log-level", fmt.Sprintf("log level. one of: %s", strings.Join(logLevel.Allowed, ",")))
	cmd.Flags().Var(logFormat, "log-format", fmt.Sprintf("log format. one of: %s", strings.Join(logFormat.Allowed, ",")))
	cmd.Flags().BoolVar(&noSentry, "no-sentry", false, "disable sentry")
	sendErrors := funk.ContainsString([]string{"all", "errors"}, telemetryLevel.String())
	if !sendErrors {
		noSentry = true
	}

	return cmd
}

const (
	pluginDocShort = "Generate documentation for tables"
	pluginDocLong  = `Generate documentation for tables

If format is markdown, a destination directory will be created (if necessary) containing markdown files.
Example:
doc ./output 

If format is JSON, a destination directory will be created (if necessary) with a single json file called __tables.json.
Example:
doc --format json .
`
)

func newCmdPluginDoc(serve *pluginServe) *cobra.Command {
	format := newEnum([]string{"json", "markdown"}, "markdown")
	cmd := &cobra.Command{
		Use:   "doc <directory>",
		Short: pluginDocShort,
		Long:  pluginDocLong,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pbFormat := pbv3.GenDocs_FORMAT(pbv3.GenDocs_FORMAT_value[format.Value])
			return serve.plugin.GeneratePluginDocs(serve.plugin.StaticTables(), args[0], pbFormat)
		},
	}
	cmd.Flags().Var(format, "format", fmt.Sprintf("output format. one of: %s", strings.Join(format.Allowed, ",")))
	return cmd
}

func newCmdPluginRoot(serve *pluginServe) *cobra.Command {
	cmd := &cobra.Command{
		Use: fmt.Sprintf("%s <command>", serve.plugin.Name()),
	}
	cmd.AddCommand(newCmdPluginServe(serve))
	cmd.AddCommand(newCmdPluginDoc(serve))
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.Version = serve.plugin.Version()
	return cmd
}
