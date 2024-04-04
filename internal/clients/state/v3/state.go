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

type Client struct {
	client  pb.PluginClient
	mem     map[string]map[uint64]string
	changes map[string]struct{} // changed keys
	latest  map[string]uint64   // latest versions
	mutex   *sync.RWMutex
	schema  *arrow.Schema
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

func NewClient(ctx context.Context, pbClient pb.PluginClient, tableName string) (*Client, error) {
	return NewClientWithTable(ctx, pbClient, Table(tableName))
}

func NewClientWithTable(ctx context.Context, pbClient pb.PluginClient, table *schema.Table) (*Client, error) {
	c := &Client{
		client:  pbClient,
		mem:     make(map[string]map[uint64]string), // key vs. version vs. value
		changes: make(map[string]struct{}),
		latest:  make(map[string]uint64),
		mutex:   &sync.RWMutex{},
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
			versions := record.Columns()[2].(*array.Uint64)
			for i := 0; i < keys.Len(); i++ {
				k, val := keys.Value(i), values.Value(i)

				if _, ok := c.mem[k]; !ok {
					c.mem[k] = make(map[uint64]string)
				}
				var ver uint64
				if versions.IsValid(i) {
					ver = versions.Value(i)
				}
				c.mem[k][ver] = val
			}
		}
	}

	for k, v := range c.mem {
		var maxVer uint64
		for ver := range v {
			if ver > maxVer {
				maxVer = ver
			}
		}
		c.latest[k] = maxVer
	}

	return c, nil
}

func (c *Client) SetKey(_ context.Context, key string, value string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.latest[key]++
	if _, ok := c.mem[key]; !ok {
		c.mem[key] = make(map[uint64]string)
	}
	c.mem[key][c.latest[key]] = value
	c.changes[key] = struct{}{}
	return nil
}

func (c *Client) Flush(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, c.schema)
	for k := range c.changes {
		ver := c.latest[k]
		val := c.mem[k][ver]
		bldr.Field(0).(*array.StringBuilder).Append(k)
		bldr.Field(1).(*array.StringBuilder).Append(val)
		bldr.Field(2).(*array.Uint64Builder).Append(ver)
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
	if ver, ok := c.latest[key]; ok {
		return c.mem[key][ver], nil
	}
	return "", nil
}
