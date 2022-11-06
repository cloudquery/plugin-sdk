package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/plugin-sdk/v1/schema"
	"github.com/cloudquery/plugin-sdk/v1/specs"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type NewDestinationClientFunc func(context.Context, zerolog.Logger, specs.Destination) (DestinationClient, error)

type DestinationClient interface {
	schema.CQTypeTransformer
	ReverseTransformValues(table *schema.Table, values []interface{}) (schema.CQTypes, error)
	Migrate(ctx context.Context, tables schema.Tables) error
	Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- []interface{}) error
	Write(ctx context.Context, tables schema.Tables, res <-chan *ClientResource) error
	Metrics() DestinationMetrics
	DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error
	Close(ctx context.Context) error
}

type DestinationPlugin struct {
	// Name of destination plugin i.e postgresql,snowflake
	name string
	// Version of the destination plugin
	version string
	// Called upon configure call to validate and init configuration
	newDestinationClient NewDestinationClientFunc
	// initialized destination client
	client DestinationClient
	// spec the client was initialized with
	spec specs.Destination
	// Logger to call, this logger is passed to the serve.Serve Client, if not define Serve will create one instead.
	logger zerolog.Logger
}

type ClientResource struct {
	TableName string
	Data      []interface{}
}

type ReverseTransformer func(*schema.Table, []interface{}) (schema.CQTypes, error)

const writeWorkers = 1

func NewDestinationPlugin(name string, version string, newDestinationClient NewDestinationClientFunc) *DestinationPlugin {
	p := &DestinationPlugin{
		name:                 name,
		version:              version,
		newDestinationClient: newDestinationClient,
	}
	return p
}

func (p *DestinationPlugin) Name() string {
	return p.name
}

func (p *DestinationPlugin) Version() string {
	return p.version
}

func (p *DestinationPlugin) Metrics() DestinationMetrics {
	return p.client.Metrics()
}

// we need lazy loading because we want to be able to initialize after
func (p *DestinationPlugin) Init(ctx context.Context, logger zerolog.Logger, spec specs.Destination) error {
	var err error
	p.logger = logger
	p.spec = spec
	p.client, err = p.newDestinationClient(ctx, logger, spec)
	if err != nil {
		return err
	}
	return nil
}

// we implement all DestinationClient functions so we can hook into pre-post behavior
func (p *DestinationPlugin) Migrate(ctx context.Context, tables schema.Tables) error {
	SetDestinationManagedCqColumns(tables)
	return p.client.Migrate(ctx, tables)
}

func (p *DestinationPlugin) Read(ctx context.Context, table *schema.Table, sourceName string, res chan<- schema.CQTypes) error {
	SetDestinationManagedCqColumns(schema.Tables{table})
	ch := make(chan []interface{})
	var err error
	go func() {
		defer close(ch)
		err = p.client.Read(ctx, table, sourceName, ch)
	}()
	for resource := range ch {
		r, err := p.client.ReverseTransformValues(table, resource)
		if err != nil {
			return err
		}
		res <- r
	}
	return err
}

func (p *DestinationPlugin) Write(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time, res <-chan *schema.DestinationResource) error {
	SetDestinationManagedCqColumns(tables)
	ch := make(chan *ClientResource)
	eg := &errgroup.Group{}
	// given most destination plugins writing in batch we are using a worker pool to write in parallel
	// it might not generalize well and we might need to move it to each destination plugin implementation.
	for i := 0; i < writeWorkers; i++ {
		eg.Go(func() error {
			return p.client.Write(ctx, tables, ch)
		})
	}
	sourceColumn := &schema.Text{}
	_ = sourceColumn.Set(sourceName)
	syncTimeColumn := &schema.Timestamptz{}
	_ = syncTimeColumn.Set(syncTime)
	for r := range res {
		r.Data = append([]schema.CQType{sourceColumn, syncTimeColumn}, r.Data...)
		clientResource := &ClientResource{
			TableName: r.TableName,
			Data:      p.transformerCqTypes(r.Data),
		}
		ch <- clientResource
	}

	close(ch)
	if err := eg.Wait(); err != nil {
		return err
	}
	if p.spec.WriteMode == specs.WriteModeOverwriteDeleteStale {
		if err := p.DeleteStale(ctx, tables, sourceName, syncTime); err != nil {
			return err
		}
	}
	return nil
}

func (p *DestinationPlugin) DeleteStale(ctx context.Context, tables schema.Tables, sourceName string, syncTime time.Time) error {
	return p.client.DeleteStale(ctx, tables, sourceName, syncTime)
}

func (p *DestinationPlugin) Close(ctx context.Context) error {
	return p.client.Close(ctx)
}

func (p *DestinationPlugin) transformerCqTypes(data schema.CQTypes) []interface{} {
	values := make([]interface{}, 0, len(data))
	for _, v := range data {
		switch v := v.(type) {
		case *schema.Bool:
			values = append(values, p.client.TransformBool(v))
		case *schema.Bytea:
			values = append(values, p.client.TransformBytea(v))
		case *schema.CIDRArray:
			values = append(values, p.client.TransformCIDRArray(v))
		case *schema.CIDR:
			values = append(values, p.client.TransformCIDR(v))
		case *schema.Float8:
			values = append(values, p.client.TransformFloat8(v))
		case *schema.InetArray:
			values = append(values, p.client.TransformInetArray(v))
		case *schema.Inet:
			values = append(values, p.client.TransformInet(v))
		case *schema.Int8:
			values = append(values, p.client.TransformInt8(v))
		case *schema.JSON:
			values = append(values, p.client.TransformJSON(v))
		case *schema.MacaddrArray:
			values = append(values, p.client.TransformMacaddrArray(v))
		case *schema.Macaddr:
			values = append(values, p.client.TransformMacaddr(v))
		case *schema.TextArray:
			values = append(values, p.client.TransformTextArray(v))
		case *schema.Text:
			values = append(values, p.client.TransformText(v))
		case *schema.Timestamptz:
			values = append(values, p.client.TransformTimestamptz(v))
		case *schema.UUIDArray:
			values = append(values, p.client.TransformUUIDArray(v))
		case *schema.UUID:
			values = append(values, p.client.TransformUUID(v))
		default:
			panic(fmt.Sprintf("unknown type %T", v))
		}
	}
	return values
}

// Overwrites or adds the CQ columns that are managed by the destination plugins (_cq_sync_time, _cq_source_name).
func SetDestinationManagedCqColumns(tables []*schema.Table) {
	for _, table := range tables {
		table.OverwriteOrAddColumn(&schema.CqSyncTimeColumn)
		table.OverwriteOrAddColumn(&schema.CqSourceNameColumn)

		SetDestinationManagedCqColumns(table.Relations)
	}
}
