package source

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
	pbBase "github.com/cloudquery/plugin-sdk/v2/internal/pb/base/v0"
	pb "github.com/cloudquery/plugin-sdk/v2/internal/pb/source/v0"
	"github.com/cloudquery/plugin-sdk/v2/internal/random"
	"github.com/cloudquery/plugin-sdk/v2/plugins/source"
	"github.com/cloudquery/plugin-sdk/v2/registry"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/specs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Client
type Client struct {
	pbClient       pb.SourceClient
	directory      string
	cmd            *exec.Cmd
	logger         zerolog.Logger
	userConn       *grpc.ClientConn
	conn           *grpc.ClientConn
	grpcSocketName string
	noSentry       bool
	wg             *sync.WaitGroup
}

type FetchResultMessage struct {
	Resource []byte
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

func WithGRPCConnection(userConn *grpc.ClientConn) func(*Client) {
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

// NewClient connect to gRPC server running source plugin and returns a new Client
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
		if err := c.newManagedClient(ctx, path); err != nil {
			return nil, err
		}
	case specs.RegistryGithub:
		pathSplit := strings.Split(path, "/")
		if len(pathSplit) != 2 {
			return nil, fmt.Errorf("invalid github plugin path: %s. format should be owner/repo", path)
		}
		org, name := pathSplit[0], pathSplit[1]
		localPath := filepath.Join(c.directory, "plugins", string(registry.PluginTypeSource), org, name, version, "plugin")
		localPath = registry.WithBinarySuffix(localPath)
		if err := registry.DownloadPluginFromGithub(ctx, localPath, org, name, version, registry.PluginTypeSource); err != nil {
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

// newManagedClient starts a new source plugin process from local path, connects to it via gRPC server
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
		return fmt.Errorf("failed to start source plugin %s: %w", path, err)
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
				c.logger.Info().Str("line", string(line)).Msg("truncated source plugin log line")
				continue
			}
			if err != nil {
				c.logger.Err(err).Msg("failed to read log line from source plugin")
				break
			}
			var structuredLogLine map[string]any
			if err := json.Unmarshal(line, &structuredLogLine); err != nil {
				c.logger.Err(err).Str("line", string(line)).Msg("failed to unmarshal log line from source plugin")
			} else {
				logging.JSONToLog(c.logger, structuredLogLine)
			}
		}
	}()

	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		d := &net.Dialer{}
		return d.DialContext(ctx, "unix", addr)
	}
	c.conn, err = grpc.DialContext(ctx, c.grpcSocketName,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(dialer),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(pb.MaxMsgSize),
			grpc.MaxCallSendMsgSize(pb.MaxMsgSize),
		),
	)
	if err != nil {
		if err := cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msg("failed to kill plugin process")
		}
		return err
	}
	c.pbClient = pb.NewSourceClient(c.conn)
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

func (c *Client) GetMetrics(ctx context.Context) (*source.Metrics, error) {
	res, err := c.pbClient.GetMetrics(ctx, &pb.GetSourceMetrics_Request{})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetMetrics: %w", err)
	}
	var stats source.Metrics
	if err := json.Unmarshal(res.Metrics, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal source stats: %w", err)
	}
	return &stats, nil
}

func (c *Client) GetTables(ctx context.Context) ([]*schema.Table, error) {
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

func (c *Client) GetTablesForSpec(ctx context.Context, spec *specs.Source) ([]*schema.Table, error) {
	var (
		b   []byte
		err error
	)
	if spec != nil {
		b, err = json.Marshal(spec)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal source spec: %w", err)
		}
	}
	res, err := c.pbClient.GetTablesForSpec(ctx, &pb.GetTablesForSpec_Request{
		Spec: b,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetTablesForSpec: %w", err)
	}
	var tables []*schema.Table
	if err := json.Unmarshal(res.Tables, &tables); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tables: %w", err)
	}
	return tables, nil
}

// Sync start syncing for the source client per the given spec and returning the results
// in the given channel. res is marshaled schema.Resource. We are not unmarshalling this for performance reasons
// as usually this is sent over-the-wire anyway to a source plugin
func (c *Client) Sync(ctx context.Context, spec specs.Source, res chan<- []byte) error {
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
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to fetch resources from stream: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res <- r.Resource:
		}
	}
}

// Sync start syncing for the source client per the given spec and returning the results
// in the given channel. res is marshaled schema.Resource. We are not unmarshalling this for performance reasons
// as usually this is sent over-the-wire anyway to a source plugin
func (c *Client) Sync2(ctx context.Context, spec specs.Source, res chan<- []byte) error {
	b, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal source spec: %w", err)
	}
	stream, err := c.pbClient.Sync2(ctx, &pb.Sync2_Request{
		Spec: b,
	})
	if err != nil {
		return fmt.Errorf("failed to call Sync: %w", err)
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to fetch resources from stream: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res <- r.Resource:
		}
	}
}

// Terminate is used only in conjunction with NewManagedClient.
// It closes the connection it created, kills the spawned process and removes the socket file.
func (c *Client) Terminate() error {
	// wait for log streaming to complete before returning from this function
	defer c.wg.Wait()

	if c.grpcSocketName != "" {
		defer func() {
			if err := os.RemoveAll(c.grpcSocketName); err != nil {
				c.logger.Error().Err(err).Msg("failed to remove source socket file")
			}
		}()
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error().Err(err).Msg("failed to close gRPC connection to source plugin")
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

func (c *Client) GetSyncSummary(ctx context.Context) (*schema.SyncSummary, error) {
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
