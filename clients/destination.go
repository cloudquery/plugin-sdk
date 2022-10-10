package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cloudquery/plugin-sdk/internal/pb"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DestinationClient struct {
	pbClient       pb.DestinationClient
	directory      string
	cmd            *exec.Cmd
	logger         zerolog.Logger
	userConn       *grpc.ClientConn
	conn           *grpc.ClientConn
	grpcSocketName string
	cmdWaitErr     error
	wg             *sync.WaitGroup
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
		wg:        &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt(c)
	}
	switch registry {
	case specs.RegistryGrpc:
		if c.userConn == nil {
			c.conn, err = grpc.DialContext(ctx, path, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		if err := DownloadPluginFromGithub(ctx, localPath, org, name, version, PluginTypeDestination); err != nil {
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

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		if err := cmd.Wait(); err != nil {
			c.cmdWaitErr = err
			c.logger.Error().Err(err).Str("plugin", path).Msg("plugin exited")
		}
	}()
	c.cmd = cmd

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		lr := newLogReader(reader)
		for {
			line, err := lr.NextLine()
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, errLogLineToLong) {
				c.logger.Err(err).Str("line", string(line)).Msg("skipping too long log line")
				continue
			}
			if err != nil {
				c.logger.Err(err).Msg("failed to read log line from plugin")
				break
			}
			var structuredLogLine map[string]interface{}
			if err := json.Unmarshal(line, &structuredLogLine); err != nil {
				c.logger.Err(err).Str("line", string(line)).Msg("failed to unmarshal log line from plugin")
			} else {
				jsonToLog(c.logger, structuredLogLine)
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
		return c, err
	}
	c.pbClient = pb.NewDestinationClient(c.conn)
	return c, nil
}

func (c *DestinationClient) Name(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to call GetName: %w", err)
	}
	return res.Name, nil
}

func (c *DestinationClient) Version(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to call GetVersion: %w", err)
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
		return fmt.Errorf("destination configure: failed to call Configure: %w", err)
	}
	return nil
}

func (c *DestinationClient) Validate(ctx context.Context, spec specs.Destination) (warnings, errors []string, err error) {
	b, err := json.Marshal(spec)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal destination spec: %w", err)
	}
	resp, err := c.pbClient.Validate(ctx, &pb.ValidateDestination_Request{
		Spec: b,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unimplemented {
			// Backwards-compatibility with older plugin versions that don't support Validate().
			// In this case, we only return one warning: that the plugin should be updated.
			return []string{"the version of this plugin is outdated and should be updated"}, nil, nil
		}
		return nil, nil, fmt.Errorf("failed to call Validate: %w", err)
	}
	return resp.Warnings, resp.Errors, nil
}

func (c *DestinationClient) Migrate(ctx context.Context, tables []*schema.Table) error {
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
// because usually it is alreadun marshalled from the source plugin.
func (c *DestinationClient) Write(ctx context.Context, source string, syncTime time.Time, resources <-chan []byte) (uint64, error) {
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
			return 0, fmt.Errorf("failed to call Write.Send: %w", err)
		}
	}
	res, err := saveClient.CloseAndRecv()
	if err != nil {
		return 0, fmt.Errorf("failed to CloseAndRecv client: %w", err)
	}

	return res.FailedWrites, nil
}

func (c *DestinationClient) Close(ctx context.Context) error {
	if _, err := c.pbClient.Close(ctx, &pb.Close_Request{}); err != nil {
		return fmt.Errorf("failed to close destination: %w", err)
	}
	return nil
}

func (c *DestinationClient) DeleteStale(ctx context.Context, tables schema.Tables, source string, timestamp time.Time) error {
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

// Terminate is used only in conjunction with NewManagedDestinationClient.
// It closes the connection it created, kills the spawned process and removes the socket file.
func (c *DestinationClient) Terminate() error {
	// wait for log streaming to complete before returning from this function
	defer c.wg.Wait()

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
			c.logger.Error().Err(err).Msg("failed to kill process")
			return err
		}
	}

	return nil
}

func (c *DestinationClient) GetWaitError() error {
	return c.cmdWaitErr
}
