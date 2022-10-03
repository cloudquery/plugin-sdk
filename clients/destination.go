package clients

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DestinationClient struct {
	pbClient       pb.DestinationClient
	directory      string
	writers        []io.Writer
	cmd            *exec.Cmd
	logger         zerolog.Logger
	userConn       *grpc.ClientConn
	conn           *grpc.ClientConn
	grpcSocketName string
	cmdWaitErr     error
}

type DestinationClientOption func(*DestinationClient)

func WithDestinationLogger(logger zerolog.Logger) func(*DestinationClient) {
	return func(c *DestinationClient) {
		c.logger = logger
	}
}

func WithDestinationDirectory(directory string) func(*DestinationClient) {
	return func(c *DestinationClient) {
		c.directory = directory
	}
}

// WithDestinationWithWriters adds writers when downloading plugins from github
func WithDestinationWithWriters(writers ...io.Writer) func(*DestinationClient) {
	return func(c *DestinationClient) {
		c.writers = writers
	}
}

func WithDestinationGrpcConn(userConn *grpc.ClientConn) func(*DestinationClient) {
	return func(c *DestinationClient) {
		// we use a different variable here because we don't want to close a connection that wasn't created by us.
		c.userConn = userConn
	}
}

func NewDestinationClient(ctx context.Context, registry specs.Registry, path string, version string, opts ...DestinationClientOption) (*DestinationClient, error) {
	var err error
	c := &DestinationClient{
		directory: DefaultDownloadDir,
	}
	for _, opt := range opts {
		opt(c)
	}
	switch registry {
	case specs.RegistryGrpc:
		if c.userConn == nil {
			c.conn, err = grpc.Dial(path, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, fmt.Errorf("failed to dial grpc source plugin at %s: %w", path, err)
			}
			c.pbClient = pb.NewDestinationClient(c.conn)
		} else {
			c.pbClient = pb.NewDestinationClient(c.userConn)
		}
		return c, nil
	case specs.RegistryLocal:
		return c.newManagedClient(ctx, path)
	case specs.RegistryGithub:
		pathSplit := strings.Split(path, "/")
		if len(pathSplit) != 2 {
			return nil, fmt.Errorf("invalid github plugin path: %s. format should be owner/repo", path)
		}
		org, name := pathSplit[0], pathSplit[1]
		localPath := filepath.Join(c.directory, "plugins", string(PluginTypeDestination), org, name, version, "plugin")
		localPath = withBinarySuffix(localPath)
		if err := DownloadPluginFromGithub(ctx, localPath, org, name, version, PluginTypeDestination, c.writers...); err != nil {
			return nil, err
		}
		return c.newManagedClient(ctx, localPath)
	default:
		return nil, fmt.Errorf("unsupported registry %s", registry)
	}
}

// newManagedClient starts a new destination plugin process from local file, connects to it via gRPC server
// and returns a new DestinationClient
func (c *DestinationClient) newManagedClient(ctx context.Context, path string) (*DestinationClient, error) {
	c.grpcSocketName = generateRandomUnixSocketName()
	// spawn the plugin first and then connect
	cmd := exec.CommandContext(ctx, path, "serve", "--network", "unix", "--address", c.grpcSocketName,
		"--log-level", c.logger.GetLevel().String(), "--log-format", "json")
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start plugin %s: %w", path, err)
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			c.cmdWaitErr = err
			c.logger.Error().Err(err).Str("plugin", path).Msg("plugin exited")
		}
	}()
	c.cmd = cmd

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			var structuredLogLine map[string]interface{}
			b := scanner.Bytes()
			if err := json.Unmarshal(b, &structuredLogLine); err != nil {
				c.logger.Err(err).Str("line", string(b)).Msg("failed to unmarshal log line from plugin")
			} else {
				jsonToLog(c.logger, structuredLogLine)
			}
		}
	}()

	c.conn, err = grpc.DialContext(ctx, "unix://"+c.grpcSocketName, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		if err := cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msg("failed to kill plugin process")
		}
		return c, err
	}
	c.pbClient = pb.NewDestinationClient(c.conn)
	return c, nil
}

func (c *DestinationClient) Name(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to get name: %w", err)
	}
	return res.Name, nil
}

func (c *DestinationClient) Version(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}
	return res.Version, nil
}

func (c *DestinationClient) Initialize(ctx context.Context, spec specs.Destination) error {
	b, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("destination configure: failed to marshal spec: %w", err)
	}
	_, err = c.pbClient.Configure(ctx, &pb.Configure_Request{
		Config: b,
	})
	if err != nil {
		return fmt.Errorf("destination configure: failed to configure: %w", err)
	}
	return nil
}

func (c *DestinationClient) Migrate(ctx context.Context, tables []*schema.Table) error {
	b, err := json.Marshal(tables)
	if err != nil {
		return fmt.Errorf("destination migrate: failed to marshal plugin: %w", err)
	}
	_, err = c.pbClient.Migrate(ctx, &pb.Migrate_Request{Tables: b})
	if err != nil {
		return fmt.Errorf("destination migrate: failed to migrate: %w", err)
	}
	return nil
}

// Write writes rows as they are received from the channel to the destination plugin.
// resources is marshaled schema.Resource. We are not marshalling this inside the function
// because usually it is alreadun marshalled from the source plugin.
func (c *DestinationClient) Write(ctx context.Context, resources <-chan []byte) (uint64, error) {
	saveClient, err := c.pbClient.Write(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to create save client: %w", err)
	}
	for resource := range resources {
		if err := saveClient.Send(&pb.Write_Request{
			Resource: resource,
		}); err != nil {
			return 0, err
		}
	}
	res, err := saveClient.CloseAndRecv()
	if err != nil {
		return 0, fmt.Errorf("failed to CloseAndRecv client: %w", err)
	}

	return res.FailedWrites, nil
}

// Close is used only in conjunction with NewManagedDestinationClient.
// It closes the connection it created, kills the spawned process and removes the socket file.
func (c *DestinationClient) Close() error {
	if c.grpcSocketName != "" {
		defer os.Remove(c.grpcSocketName)
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error().Err(err).Msg("failed to close gRPC connection")
		}
	}
	if c.cmd != nil && c.cmd.Process != nil {
		if err := c.cmd.Process.Kill(); err != nil {
			return err
		}
	}
	return nil
}

func (c *DestinationClient) GetWaitError() error {
	return c.cmdWaitErr
}
