package destination

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/internal/logging"
	pbBase "github.com/cloudquery/plugin-sdk/internal/pb/base/v0"
	pb "github.com/cloudquery/plugin-sdk/internal/pb/destination/v0"
	"github.com/cloudquery/plugin-sdk/internal/random"
	"github.com/cloudquery/plugin-sdk/plugins/destination"
	"github.com/cloudquery/plugin-sdk/registry"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	pbClient       pb.DestinationClient
	directory      string
	cmd            *exec.Cmd
	logger         zerolog.Logger
	userConn       *grpc.ClientConn
	conn           *grpc.ClientConn
	grpcSocketName string
	noSentry       bool
	wg             *sync.WaitGroup
}

type ClientOption func(*Client)

func WithLogger(logger zerolog.Logger) func(*Client) {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithDirectory(directory string) func(*Client) {
	return func(c *Client) {
		c.directory = directory
	}
}

func WithGrpcConn(userConn *grpc.ClientConn) func(*Client) {
	return func(c *Client) {
		// we use a different variable here because we don't want to close a connection that wasn't created by us.
		c.userConn = userConn
	}
}

func WithNoSentry() func(*Client) {
	return func(c *Client) {
		c.noSentry = true
	}
}

func NewClient(ctx context.Context, registrySpec specs.Registry, path string, version string, opts ...ClientOption) (*Client, error) {
	var err error
	c := &Client{
		directory: registry.DefaultDownloadDir,
		wg:        &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt(c)
	}
	switch registrySpec {
	case specs.RegistryGrpc:
		if c.userConn == nil {
			c.conn, err = grpc.DialContext(ctx, path,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultCallOptions(
					grpc.MaxCallRecvMsgSize(maxMsgSize),
					grpc.MaxCallSendMsgSize(maxMsgSize),
				),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to dial grpc source plugin at %s: %w", path, err)
			}
			c.pbClient = pb.NewDestinationClient(c.conn)
		} else {
			c.pbClient = pb.NewDestinationClient(c.userConn)
		}
		return c, nil
	case specs.RegistryLocal:
		if err := c.newManagedClient(ctx, path); err != nil {
			return nil, err
		}
	case specs.RegistryGithub:
		pathSplit := strings.Split(path, "/")
		if len(pathSplit) != 2 {
			return nil, fmt.Errorf("invalid github plugin path: %s. format should be owner/repo", path)
		}
		org, name := pathSplit[0], pathSplit[1]
		localPath := filepath.Join(c.directory, "plugins", string(registry.PluginTypeDestination), org, name, version, "plugin")
		localPath = registry.WithBinarySuffix(localPath)
		if err := registry.DownloadPluginFromGithub(ctx, localPath, org, name, version, registry.PluginTypeDestination); err != nil {
			return nil, err
		}
		if err := c.newManagedClient(ctx, localPath); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported registry %s", registrySpec)
	}

	return c, nil
}

// newManagedClient starts a new destination plugin process from local file, connects to it via gRPC server
// and returns a new Client
func (c *Client) newManagedClient(ctx context.Context, path string) error {
	c.grpcSocketName = random.GenerateRandomUnixSocketName()
	// spawn the plugin first and then connect
	args := []string{"serve", "--network", "unix", "--address", c.grpcSocketName,
		"--log-level", c.logger.GetLevel().String(), "--log-format", "json"}
	if c.noSentry {
		args = append(args, "--no-sentry")
	}
	cmd := exec.CommandContext(ctx, path, args...)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start destination plugin %s: %w", path, err)
	}

	c.cmd = cmd

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		lr := logging.NewLogReader(reader)
		for {
			line, err := lr.NextLine()
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, logging.ErrLogLineToLong) {
				c.logger.Info().Str("line", string(line)).Msg("truncated destination plugin log line")
				continue
			}
			if err != nil {
				c.logger.Err(err).Msg("failed to read log line from destination plugin")
				break
			}
			var structuredLogLine map[string]any
			if err := json.Unmarshal(line, &structuredLogLine); err != nil {
				c.logger.Err(err).Str("line", string(line)).Msg("failed to unmarshal log line from destination plugin")
			} else {
				logging.JSONToLog(c.logger, structuredLogLine)
			}
		}
	}()

	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		d := &net.Dialer{}
		return d.DialContext(ctx, "unix", addr)
	}
	c.conn, err = grpc.DialContext(ctx, c.grpcSocketName, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithContextDialer(dialer))
	if err != nil {
		if err := cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msg("failed to kill plugin process")
		}
		return err
	}
	c.pbClient = pb.NewDestinationClient(c.conn)
	return nil
}

func (c *Client) GetProtocolVersion(ctx context.Context) (uint64, error) {
	res, err := c.pbClient.GetProtocolVersion(ctx, &pbBase.GetProtocolVersion_Request{})
	if err != nil {
		s := status.Convert(err)
		if s.Code() != codes.Unimplemented {
			return 0, err
		}
		c.logger.Warn().Err(err).Msg("plugin does not support protocol version. assuming protocol version 1")
		return 1, nil
	}
	return res.Version, nil
}

func (c *Client) GetMetrics(ctx context.Context) (*destination.Metrics, error) {
	res, err := c.pbClient.GetMetrics(ctx, &pb.GetDestinationMetrics_Request{})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetMetrics: %w", err)
	}
	var stats destination.Metrics
	if err := json.Unmarshal(res.Metrics, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal destination metrics: %w", err)
	}
	return &stats, nil
}

func (c *Client) Name(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetName(ctx, &pbBase.GetName_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to call GetName: %w", err)
	}
	return res.Name, nil
}

func (c *Client) Version(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetVersion(ctx, &pbBase.GetVersion_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to call GetVersion: %w", err)
	}
	return res.Version, nil
}

func (c *Client) Initialize(ctx context.Context, spec specs.Destination) error {
	b, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("destination configure: failed to marshal spec: %w", err)
	}
	_, err = c.pbClient.Configure(ctx, &pbBase.Configure_Request{
		Config: b,
	})
	if err != nil {
		return fmt.Errorf("destination configure: failed to call Configure: %w", err)
	}
	return nil
}

func (c *Client) Migrate(ctx context.Context, tables []*schema.Table) error {
	b, err := json.Marshal(tables)
	if err != nil {
		return fmt.Errorf("destination migrate: failed to marshal plugin: %w", err)
	}
	_, err = c.pbClient.Migrate(ctx, &pb.Migrate_Request{Tables: b})
	if err != nil {
		return fmt.Errorf("failed to call Migrate: %w", err)
	}
	return nil
}

// Write writes rows as they are received from the channel to the destination plugin.
// resources is marshaled schema.Resource. We are not marshalling this inside the function
// because usually it is already marshalled from the destination plugin.
func (c *Client) Write(ctx context.Context, source string, syncTime time.Time, resources <-chan []byte) (uint64, error) {
	saveClient, err := c.pbClient.Write(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to call Write: %w", err)
	}
	for resource := range resources {
		if err := saveClient.Send(&pb.Write_Request{
			Resource:  resource,
			Source:    source,
			Timestamp: timestamppb.New(syncTime),
		}); err != nil {
			if err == io.EOF {
				// don't send write request if the channel is closed
				break
			}
			return 0, fmt.Errorf("failed to call Write.Send: %w", err)
		}
	}
	res, err := saveClient.CloseAndRecv()
	if err != nil {
		return 0, fmt.Errorf("failed to CloseAndRecv client: %w", err)
	}

	return res.FailedWrites, nil
}

func (c *Client) Write2(ctx context.Context, sourceSpec specs.Source, tables schema.Tables, syncTime time.Time, resources <-chan []byte) error {
	saveClient, err := c.pbClient.Write2(ctx)
	if err != nil {
		return fmt.Errorf("failed to call Write2: %w", err)
	}
	b, err := json.Marshal(tables)
	if err != nil {
		return fmt.Errorf("failed to marshal tables: %w", err)
	}
	sourceSpecBytes, err := json.Marshal(sourceSpec)
	if err != nil {
		return fmt.Errorf("failed to marshal source spec: %w", err)
	}
	if err := saveClient.Send(&pb.Write2_Request{
		Tables:     b,
		Source:     sourceSpec.Name,
		Timestamp:  timestamppb.New(syncTime),
		SourceSpec: sourceSpecBytes,
	}); err != nil {
		return fmt.Errorf("failed to send tables: %w", err)
	}
	for resource := range resources {
		if err := saveClient.Send(&pb.Write2_Request{
			Resource: resource,
		}); err != nil {
			if err == io.EOF {
				// don't send write request if the channel is closed
				break
			}
			return fmt.Errorf("failed to call Write2.Send: %w", err)
		}
	}
	_, err = saveClient.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to CloseAndRecv client: %w", err)
	}

	return nil
}

func (c *Client) Close(ctx context.Context) error {
	if _, err := c.pbClient.Close(ctx, &pb.Close_Request{}); err != nil {
		return fmt.Errorf("failed to close destination: %w", err)
	}
	return nil
}

func (c *Client) DeleteStale(ctx context.Context, tables schema.Tables, source string, timestamp time.Time) error {
	b, err := json.Marshal(tables)
	if err != nil {
		return fmt.Errorf("destination delete stale: failed to marshal plugin: %w", err)
	}
	if _, err := c.pbClient.DeleteStale(ctx, &pb.DeleteStale_Request{
		Source:    source,
		Timestamp: timestamppb.New(timestamp),
		Tables:    b,
	}); err != nil {
		return fmt.Errorf("failed to call DeleteStale: %w", err)
	}
	return nil
}

// Terminate is used only in conjunction with NewManagedClient.
// It closes the connection it created, kills the spawned process and removes the socket file.
func (c *Client) Terminate() error {
	// wait for log streaming to complete before returning from this function
	defer c.wg.Wait()

	if c.grpcSocketName != "" {
		defer func() {
			if err := os.RemoveAll(c.grpcSocketName); err != nil {
				c.logger.Error().Err(err).Msg("failed to remove destination socket file")
			}
		}()
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error().Err(err).Msg("failed to close gRPC connection to destination plugin")
		}
		c.conn = nil
	}
	if c.cmd != nil && c.cmd.Process != nil {
		if err := c.terminateProcess(); err != nil {
			return err
		}
	}

	return nil
}
