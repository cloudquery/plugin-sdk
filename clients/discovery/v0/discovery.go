package discovery

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

	"github.com/cloudquery/plugin-sdk/v2/internal/logging"
	pb "github.com/cloudquery/plugin-sdk/v2/internal/pb/discovery/v0"
	"github.com/cloudquery/plugin-sdk/v2/internal/random"
	"github.com/cloudquery/plugin-sdk/v2/registry"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pbClient       pb.DiscoveryClient
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

func NewClient(ctx context.Context, registrySpec specs.Registry, pluginType registry.PluginType, path string, version string, opts ...ClientOption) (*Client, error) {
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
			c.pbClient = pb.NewDiscoveryClient(c.conn)
		} else {
			c.pbClient = pb.NewDiscoveryClient(c.userConn)
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
		localPath := filepath.Join(c.directory, "plugins", string(pluginType), org, name, version, "plugin")
		localPath = registry.WithBinarySuffix(localPath)
		if err := registry.DownloadPluginFromGithub(ctx, localPath, org, name, version, pluginType); err != nil {
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

// newManagedClient starts a new discovery plugin process from local file, connects to it via gRPC server
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
	c.pbClient = pb.NewDiscoveryClient(c.conn)
	return nil
}

func (c *Client) GetVersions(ctx context.Context) ([]string, error) {
	res, err := c.pbClient.GetVersions(ctx, &pb.GetVersions_Request{})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetVersions: %w", err)
	}
	return res.Versions, nil
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
