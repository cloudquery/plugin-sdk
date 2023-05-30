package serve

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	pbv0 "github.com/cloudquery/plugin-pb-go/pb/destination/v0"
	pbv1 "github.com/cloudquery/plugin-pb-go/pb/destination/v1"
	pbdiscoveryv0 "github.com/cloudquery/plugin-pb-go/pb/discovery/v0"
	servers "github.com/cloudquery/plugin-sdk/v3/internal/servers/destination/v0"
	serversv1 "github.com/cloudquery/plugin-sdk/v3/internal/servers/destination/v1"
	discoveryServerV0 "github.com/cloudquery/plugin-sdk/v3/internal/servers/discovery/v0"
	"github.com/cloudquery/plugin-sdk/v3/plugins/destination"
	"github.com/cloudquery/plugin-sdk/v3/types"
	"github.com/getsentry/sentry-go"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type destinationServe struct {
	plugin    *destination.Plugin
	sentryDSN string
}

type DestinationOption func(*destinationServe)

func WithDestinationSentryDSN(dsn string) DestinationOption {
	return func(s *destinationServe) {
		s.sentryDSN = dsn
	}
}

var testDestinationListener *bufconn.Listener
var testDestinationListenerLock sync.Mutex

const serveDestinationShort = `Start destination plugin server`

func Destination(plugin *destination.Plugin, opts ...DestinationOption) {
	s := &destinationServe{
		plugin: plugin,
	}
	for _, opt := range opts {
		opt(s)
	}
	if err := newCmdDestinationRoot(s).Execute(); err != nil {
		sentry.CaptureMessage(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}

// nolint:dupl
func newCmdDestinationServe(serve *destinationServe) *cobra.Command {
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
		Short: serveDestinationShort,
		Long:  serveDestinationShort,
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

			var listener net.Listener
			if network == "test" {
				testDestinationListenerLock.Lock()
				listener = bufconn.Listen(testBufSize)
				testDestinationListener = listener.(*bufconn.Listener)
				testDestinationListenerLock.Unlock()
			} else {
				listener, err = net.Listen(network, address)
				if err != nil {
					return fmt.Errorf("failed to listen %s:%s: %w", network, address, err)
				}
			}
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
			pbv0.RegisterDestinationServer(s, &servers.Server{
				Plugin: serve.plugin,
				Logger: logger,
			})
			pbv1.RegisterDestinationServer(s, &serversv1.Server{
				Plugin: serve.plugin,
				Logger: logger,
			})
			pbdiscoveryv0.RegisterDiscoveryServer(s, &discoveryServerV0.Server{
				Versions: []string{"v0", "v1"},
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

			if err := types.RegisterAllExtensions(); err != nil {
				return err
			}
			defer func() {
				if err := types.UnregisterAllExtensions(); err != nil {
					logger.Error().Err(err).Msg("Failed to unregister extensions")
				}
			}()

			ctx := cmd.Context()
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			defer func() {
				signal.Stop(c)
			}()

			go func() {
				select {
				case sig := <-c:
					logger.Info().Str("address", listener.Addr().String()).Str("signal", sig.String()).Msg("Got stop signal. Destination plugin server shutting down")
					s.Stop()
				case <-ctx.Done():
					logger.Info().Str("address", listener.Addr().String()).Msg("Context cancelled. Destination plugin server shutting down")
					s.Stop()
				}
			}()

			logger.Info().Str("address", listener.Addr().String()).Msg("Destination plugin server listening")
			if err := s.Serve(listener); err != nil {
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

func newCmdDestinationRoot(serve *destinationServe) *cobra.Command {
	cmd := &cobra.Command{
		Use: fmt.Sprintf("%s <command>", serve.plugin.Name()),
	}
	cmd.AddCommand(newCmdDestinationServe(serve))
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.Version = serve.plugin.Version()
	return cmd
}
