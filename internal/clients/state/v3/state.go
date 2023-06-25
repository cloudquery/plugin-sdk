package state

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/ipc"
	"github.com/apache/arrow/go/v13/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

const keyColumn = "key"
const valueColumn = "value"

type Client struct {
	client    pb.PluginClient
	tableName string
	mem       map[string]string
	mutex     *sync.RWMutex
	keys      []string
	values    []string
	schema    *arrow.Schema
}

func NewClient(ctx context.Context, pbClient pb.PluginClient, tableName string) (*Client, error) {
	c := &Client{
		client:    pbClient,
		tableName: tableName,
		mem:       make(map[string]string),
		mutex:     &sync.RWMutex{},
		keys:      make([]string, 0),
		values:    make([]string, 0),
	}
	table := &schema.Table{
		Name: tableName,
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
		Message: &pb.Write_Request_Options{
			Options: &pb.WriteOptions{MigrateForce: false},
		},
	}); err != nil {
		return nil, err
	}
	if err := writeClient.Send(&pb.Write_Request{
		Message: &pb.Write_Request_MigrateTable{
			MigrateTable: &pb.MessageMigrateTable{
				Table: tableBytes,
			},
		},
	}); err != nil {
		return nil, err
	}

	syncClient, err := c.client.Sync(ctx, &pb.Sync_Request{
		Tables: []string{tableName},
	})
	if err != nil {
		return nil, err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for {
		res, err := syncClient.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		var insertMessage *pb.Sync_Response_Insert
		switch m := res.Message.(type) {
		case *pb.Sync_Response_Delete:
			continue
		case *pb.Sync_Response_MigrateTable:
			continue
		case *pb.Sync_Response_Insert:
			insertMessage = m
		}
		rdr, err := ipc.NewReader(bytes.NewReader(insertMessage.Insert.Record))
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
			for i := 0; i < keys.Len(); i++ {
				c.mem[keys.Value(i)] = values.Value(i)
			}
		}
	}
	return c, nil
}

func (c *Client) SetKey(_ context.Context, key string, value string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.mem[key] = value
	return nil
}

func (c *Client) Flush(ctx context.Context) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, c.schema)
	for k, v := range c.mem {
		bldr.Field(0).(*array.StringBuilder).Append(k)
		bldr.Field(1).(*array.StringBuilder).Append(v)
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
		Message: &pb.Write_Request_Options{},
	}); err != nil {
		return err
	}
	if err := writeClient.Send(&pb.Write_Request{
		Message: &pb.Write_Request_Insert{
			Insert: &pb.MessageInsert{
				Record: recordBytes,
			},
		},
	}); err != nil {
		return err
	}
	if _, err := writeClient.CloseAndRecv(); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetKey(_ context.Context, key string) (string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.mem[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}
