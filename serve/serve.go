package serve

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/cloudquery/cq-provider-sdk/internal/pb"
	"github.com/cloudquery/cq-provider-sdk/internal/servers"
	"github.com/cloudquery/cq-provider-sdk/plugins"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type Options struct {
	// Required: Provider is the actual provider that will be served.
	SourcePlugin      *plugins.SourcePlugin
	DestinationPlugin plugins.DestinationPlugin
}

const (
	serveShort = `Start plugin server`
)

func newCmdServe(opts Options) *cobra.Command {
	var address string
	var network string
	logLevel := newEnum([]string{"trace", "debug", "info", "warn", "error"}, "info")
	logFormat := newEnum([]string{"text", "json"}, "text")
	cmd := &cobra.Command{
		Use:   "serve",
		Short: serveShort,
		Long:  serveShort,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			zerologLevel, err := zerolog.ParseLevel(logLevel.String())
			if err != nil {
				return err
			}
			var logger zerolog.Logger
			if logFormat.String() == "json" {
				logger = zerolog.New(os.Stderr).Level(zerologLevel)
			} else {
				logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerologLevel)
			}
			// opts.Plugin.Logger = logger
			listener, err := net.Listen(network, address)
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}
			// See logging pattern https://github.com/grpc-ecosystem/go-grpc-middleware/blob/v2/providers/zerolog/examples_test.go
			s := grpc.NewServer(
				middleware.WithUnaryServerChain(
					logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger)),
				),
				middleware.WithStreamServerChain(
					logging.StreamServerInterceptor(grpczerolog.InterceptorLogger(logger)),
				),
				// grpc.ChainStreamInterceptor(grpc_zero),
				// grpc.ChainUnaryInterceptor(),
			)

			if opts.SourcePlugin != nil {
				opts.SourcePlugin.Logger = logger
				pb.RegisterSourceServer(s, &servers.SourceServer{Plugin: opts.SourcePlugin})
			}
			if opts.DestinationPlugin != nil {
				opts.SourcePlugin.Logger = logger
				pb.RegisterDestinationServer(s, &servers.DestinationServer{Plugin: opts.DestinationPlugin})
			}

			logger.Info().Str("address", listener.Addr().String()).Msg("server listening")
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
	return cmd
}

func newCmdRoot(opts Options) *cobra.Command {
	cmd := &cobra.Command{
		Use: "plugin <command>",
	}
	cmd.AddCommand(newCmdServe(opts))
	cmd.AddCommand(newCmdDoc(opts))
	return cmd
}

func Serve(opts Options) {
	if err := newCmdRoot(opts).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
