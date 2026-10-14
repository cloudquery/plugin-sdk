package state

import (
	"bytes"
	"context"
	"errors"
	"io"
	"iter"
	"maps"
	"slices"
	"sync"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/checkpoint"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"google.golang.org/grpc"
)

const (
	keyColumn     = "key"
	valueColumn   = "value"
	versionColumn = "version"
)

type Client struct {
	client pb.PluginClient
	mem    map[string]versionedValue
	// changed keys, with the scope covering every change since the last flush (nil outside table resolvers)
	changes map[string]checkpoint.Scope
	// the last checkpoint of each key, which the next checkpoint of the key must not overtake
	sent          map[string]*checkpoint.Sent
	mutex         *sync.RWMutex
	schema        *arrow.Schema
	versionedMode bool
	conn          *grpc.ClientConn
	// when set, Flush sends checkpoints instead of writing to the backend
	checkpoints chan<- message.SyncMessage
}

type Option func(*Client)

// WithCheckpoints makes Flush send the changes as message.SyncCheckpoint on res instead of writing them to the
// backend. Changes made within a table resolver are sent after the resources it produced before the Flush.
func WithCheckpoints(res chan<- message.SyncMessage) Option {
	return func(c *Client) {
		c.checkpoints = res
	}
}

type versionedValue struct {
	value   string
	version uint64
}

func NewClient(ctx context.Context, conn *grpc.ClientConn, tableName string, opts ...Option) (*Client, error) {
	table := Table(tableName)
	c := &Client{
		conn:          conn,
		client:        pb.NewPluginClient(conn),
		mem:           make(map[string]versionedValue),
		changes:       make(map[string]checkpoint.Scope),
		sent:          make(map[string]*checkpoint.Sent),
		mutex:         &sync.RWMutex{},
		versionedMode: table.Column(versionColumn) != nil,
	}
	for _, opt := range opts {
		opt(c)
	}
	sc := table.ToArrowSchema()
	c.schema = sc
	tableBytes, err := pb.SchemaToBytes(sc)
	if err != nil {
		return nil, err
	}

	writeClient, err := c.client.Write(ctx)
	if err != nil {
		return nil, err
	}
	if err := writeClient.Send(&pb.Write_Request{
		Message: &pb.Write_Request_MigrateTable{
			MigrateTable: &pb.Write_MessageMigrateTable{
				Table: tableBytes,
			},
		},
	}); err != nil {
		return nil, err
	}
	if _, err := writeClient.CloseAndRecv(); err != nil {
		return nil, err
	}

	readClient, err := c.client.Read(ctx, &pb.Read_Request{
		Table: tableBytes,
	})
	if err != nil {
		return nil, err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for {
		res, err := readClient.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		rdr, err := ipc.NewReader(bytes.NewReader(res.Record))
		if err != nil {
			return nil, err
		}
		for {
			record, err := rdr.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			if record.NumRows() == 0 {
				continue
			}
			keys := record.Columns()[0].(*array.String)
			values := record.Columns()[1].(*array.String)

			var versions *array.Uint64
			if c.versionedMode {
				versions = record.Columns()[2].(*array.Uint64)
			}
			for i := 0; i < keys.Len(); i++ {
				k, val := keys.Value(i), values.Value(i)

				var ver uint64
				if versions != nil && versions.IsValid(i) {
					ver = versions.Value(i)
				}
				if cur, ok := c.mem[k]; ok {
					if cur.version > ver {
						continue
					}
				}
				c.mem[k] = versionedValue{
					value:   val,
					version: ver,
				}
			}
		}
	}

	return c, nil
}

func (c *Client) SetKey(ctx context.Context, key string, value string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.mem[key].value == value {
		return nil // don't update if the value is the same
	}
	c.mem[key] = versionedValue{
		value:   value,
		version: c.mem[key].version + 1,
	}
	c.changes[key] = checkpoint.Join(c.changes[key], checkpoint.FromContext(ctx))
	return nil
}

func (c *Client) Flush(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.checkpoints != nil {
		return c.flushCheckpoints(ctx)
	}
	rec := c.record(maps.Keys(c.changes))
	recordBytes, err := pb.RecordToBytes(rec)
	if err != nil {
		return err
	}
	writeClient, err := c.client.Write(ctx)
	if err != nil {
		return err
	}
	if err := writeClient.Send(&pb.Write_Request{
		Message: &pb.Write_Request_Insert{
			Insert: &pb.Write_MessageInsert{
				Record: recordBytes,
			},
		},
	}); err != nil {
		return err
	}
	if _, err := writeClient.CloseAndRecv(); err != nil {
		return err
	}

	c.changes = make(map[string]checkpoint.Scope)
	return nil
}

// flushCheckpoints sends one checkpoint per scope, so each waits only for the resources of the scopes that changed its
// keys. Keys whose checkpoint fails stay pending.
func (c *Client) flushCheckpoints(ctx context.Context) error {
	keysByScope := make(map[checkpoint.Scope][]string)
	for key, scope := range c.changes {
		keysByScope[scope] = append(keysByScope[scope], key)
	}
	// keys changed outside table resolvers do not wait for any scope, so they go first
	scopes := slices.SortedFunc(maps.Keys(keysByScope), func(a, b checkpoint.Scope) int {
		switch {
		case a == nil && b != nil:
			return -1
		case a != nil && b == nil:
			return 1
		default:
			return 0
		}
	})
	var errs []error
	for _, scope := range scopes {
		keys := keysByScope[scope]
		msg := &message.SyncCheckpoint{Record: c.record(slices.Values(keys))}
		sent, err := checkpoint.Send(ctx, c.checkpoints, scope, c.unsent(keys), msg)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, key := range keys {
			delete(c.changes, key)
			c.sent[key] = sent
		}
	}
	return errors.Join(errs...)
}

// unsent returns the previous checkpoints of keys that are not sent yet, as the scopes they wait for may still be
// sending the resources they cover.
func (c *Client) unsent(keys []string) []*checkpoint.Sent {
	var unsent []*checkpoint.Sent
	for _, key := range keys {
		sent, ok := c.sent[key]
		if !ok {
			continue
		}
		if sent.Done() {
			delete(c.sent, key)
		} else if !slices.Contains(unsent, sent) {
			unsent = append(unsent, sent)
		}
	}
	return unsent
}

func (c *Client) record(keys iter.Seq[string]) arrow.RecordBatch {
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, c.schema)
	defer bldr.Release()
	keyBldr, valueBldr := bldr.Field(0).(*array.StringBuilder), bldr.Field(1).(*array.StringBuilder)
	var versionBldr *array.Uint64Builder
	if c.versionedMode {
		versionBldr = bldr.Field(2).(*array.Uint64Builder)
	}
	for k := range keys {
		val := c.mem[k]
		keyBldr.Append(k)
		valueBldr.Append(val.value)
		if versionBldr != nil {
			versionBldr.Append(val.version)
		}
	}
	return bldr.NewRecordBatch()
}

func (c *Client) GetKey(_ context.Context, key string) (string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.mem[key].value, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func Table(name string) *schema.Table {
	return &schema.Table{
		Name: name,
		Columns: []schema.Column{
			{
				Name:       keyColumn,
				Type:       arrow.BinaryTypes.String,
				PrimaryKey: true,
			},
			{
				Name: valueColumn,
				Type: arrow.BinaryTypes.String,
			},
			{
				// Not defined as PrimaryKey to enable single keys if the destination supports PKs
				Name: versionColumn,
				Type: arrow.PrimitiveTypes.Uint64,
			},
		},
	}
}
