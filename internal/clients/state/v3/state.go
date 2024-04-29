package state

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/ipc"
	"github.com/apache/arrow/go/v15/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

const (
	keyColumn     = "key"
	valueColumn   = "value"
	versionColumn = "version"
)

type Closer interface {
	Close() error
}

type Client struct {
	client        pb.PluginClient
	mem           map[string]versionedValue
	changes       map[string]struct{} // changed keys
	mutex         *sync.RWMutex
	schema        *arrow.Schema
	versionedMode bool
	conn          Closer
}

type versionedValue struct {
	value   string
	version uint64
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
		},
	}
}

func VersionedTable(name string) *schema.Table {
	t := Table(name)
	t.Columns = append(t.Columns, schema.Column{
		// Not defined as PrimaryKey to enable single keys if the destination supports PKs
		Name: versionColumn,
		Type: arrow.PrimitiveTypes.Uint64,
	})
	return t
}

func NewClient(ctx context.Context, pbClient pb.PluginClient, tableName string, conn Closer) (*Client, error) {
	return NewClientWithTable(ctx, pbClient, Table(tableName), conn)
}

func NewClientWithTable(ctx context.Context, pbClient pb.PluginClient, table *schema.Table, conn Closer) (*Client, error) {
	c := &Client{
		conn:          conn,
		client:        pbClient,
		mem:           make(map[string]versionedValue),
		changes:       make(map[string]struct{}),
		mutex:         &sync.RWMutex{},
		versionedMode: table.Column(versionColumn) != nil,
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

func (c *Client) SetKey(_ context.Context, key string, value string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.mem[key].value == value {
		return nil // don't update if the value is the same
	}
	c.mem[key] = versionedValue{
		value:   value,
		version: c.mem[key].version + 1,
	}
	c.changes[key] = struct{}{}
	return nil
}

func (c *Client) Flush(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, c.schema)
	for k := range c.changes {
		val := c.mem[k]
		bldr.Field(0).(*array.StringBuilder).Append(k)
		bldr.Field(1).(*array.StringBuilder).Append(val.value)
		if c.versionedMode {
			bldr.Field(2).(*array.Uint64Builder).Append(val.version)
		}
	}
	rec := bldr.NewRecord()
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

	c.changes = make(map[string]struct{})
	return nil
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
