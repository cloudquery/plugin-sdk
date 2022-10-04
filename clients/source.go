package clients

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
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

// SourceClient
type SourceClient struct {
	pbClient       pb.SourceClient
	directory      string
	cmd            *exec.Cmd
	logger         zerolog.Logger
	userConn       *grpc.ClientConn
	conn           *grpc.ClientConn
	grpcSocketName string
	cmdWaitErr     error
}

type FetchResultMessage struct {
	Resource []byte
}

type SourceClientOption func(*SourceClient)

func WithSourceLogger(logger zerolog.Logger) func(*SourceClient) {
	return func(c *SourceClient) {
		c.logger = logger
	}
}

func WithSourceDirectory(directory string) func(*SourceClient) {
	return func(c *SourceClient) {
		c.directory = directory
	}
}

func WithSourceGRPCConnection(userConn *grpc.ClientConn) func(*SourceClient) {
	return func(c *SourceClient) {
		// we use a different variable here because we don't want to close a connection that wasn't created by us.
		c.userConn = userConn
	}
}

// NewSourceClient connect to gRPC server running source plugin and returns a new SourceClient
func NewSourceClient(ctx context.Context, registry specs.Registry, path string, version string, opts ...SourceClientOption) (*SourceClient, error) {
	var err error
	c := &SourceClient{
		directory: DefaultDownloadDir,
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
			c.pbClient = pb.NewSourceClient(c.conn)
		} else {
			c.pbClient = pb.NewSourceClient(c.userConn)
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
		localPath := filepath.Join(c.directory, "plugins", string(PluginTypeSource), org, name, version, "plugin")
		localPath = withBinarySuffix(localPath)
		if err := DownloadPluginFromGithub(ctx, localPath, org, name, version, PluginTypeSource); err != nil {
			return nil, err
		}
		return c.newManagedClient(ctx, localPath)
	default:
		return nil, fmt.Errorf("unsupported registry %s", registry)
	}
}

// newManagedClient starts a new source plugin process from local path, connects to it via gRPC server
// and returns a new SourceClient
func (c *SourceClient) newManagedClient(ctx context.Context, path string) (*SourceClient, error) {
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
	c.pbClient = pb.NewSourceClient(c.conn)
	return c, nil
}

func (c *SourceClient) Name(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to call GetName: %w", err)
	}
	return res.Name, nil
}

func (c *SourceClient) Version(ctx context.Context) (string, error) {
	res, err := c.pbClient.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		return "", fmt.Errorf("failed to call GetVersion: %w", err)
	}
	return res.Version, nil
}

func (c *SourceClient) GetTables(ctx context.Context) ([]*schema.Table, error) {
	res, err := c.pbClient.GetTables(ctx, &pb.GetTables_Request{})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetTables: %w", err)
	}
	var tables []*schema.Table
	if err := json.Unmarshal(res.Tables, &tables); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tables: %w", err)
	}
	return tables, nil
}

// Sync start syncing for the source client per the given spec and returning the results
// in the given channel. res is marshaled schema.Resource. We are not unmarshalling this for performance reasons
// as usually this is sent over-the-wire anyway to a destination plugin
func (c *SourceClient) Sync(ctx context.Context, spec specs.Source, res chan<- []byte) error {
	b, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal source spec: %w", err)
	}
	stream, err := c.pbClient.Sync(ctx, &pb.Sync_Request{
		Spec: b,
	})
	if err != nil {
		return fmt.Errorf("failed to call Sync: %w", err)
	}
	for {
		r, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to fetch resources from stream: %w", err)
		}
		res <- r.Resource
	}
}

// Terminate is used only in conjunction with NewManagedSourceClient.
// It closes the connection it created, kills the spawned process and removes the socket file.
func (c *SourceClient) Terminate() error {
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

func (c *SourceClient) GetSyncSummary(ctx context.Context) (*schema.SyncSummary, error) {
	res, err := c.pbClient.GetSyncSummary(ctx, &pb.GetSyncSummary_Request{})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetSyncSummary: %w", err)
	}
	var summary schema.SyncSummary
	if err := json.Unmarshal(res.Summary, &summary); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sync summary: %w", err)
	}
	return &summary, nil
}

func (c *SourceClient) GetWaitError() error {
	return c.cmdWaitErr
}
