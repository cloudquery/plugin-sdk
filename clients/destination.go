package clients

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DestinationClient struct {
	pbClient       pb.DestinationClient
	cmd            *exec.Cmd
	logger         zerolog.Logger
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

func NewDestinationClient(cc grpc.ClientConnInterface) *DestinationClient {
	return &DestinationClient{
		pbClient: pb.NewDestinationClient(cc),
	}
}

// NewManagedDestinationClient starts a new destination plugin process, connects to it via gRPC server
// and returns a new DestinationClient
func NewManagedDestinationClient(ctx context.Context, path string, opts ...DestinationClientOption) (*DestinationClient, error) {
	c := &DestinationClient{
		logger: log.Logger,
	}
	for _, opt := range opts {
		opt(c)
	}
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
