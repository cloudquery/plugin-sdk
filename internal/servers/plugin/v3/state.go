package plugin

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/ipc"
	"github.com/apache/arrow/go/v13/arrow/memory"
	pbDiscovery "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	pbPlugin "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/state"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
)

const stateTablePrefix = "cq_state_"
const keyColumn = "key"
const valueColumn = "value"

type ClientV3 struct {
	client pbPlugin.PluginClient
	mem    map[string]string
	keys   []string
	values []string
}

func newStateClient(ctx context.Context, conn *grpc.ClientConn, spec *pbPlugin.StateBackendSpec) (state.Client, error) {
	discoveryClient := pbDiscovery.NewDiscoveryClient(conn)
	versions, err := discoveryClient.GetVersions(ctx, &pbDiscovery.GetVersions_Request{})
	if err != nil {
		return nil, err
	}
	if !slices.Contains(versions.Versions, 3) {
		return nil, fmt.Errorf("please upgrade your state backend plugin")
	}

	c := &ClientV3{
		client: pbPlugin.NewPluginClient(conn),
		mem:    make(map[string]string),
		keys:   make([]string, 0),
		values: make([]string, 0),
	}
	name := spec.Name
	table := &schema.Table{
		Name: stateTablePrefix + name,
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
	tableBytes, err := table.ToArrowSchemaBytes()
	if err != nil {
		return nil, err
	}

	if _, err := c.client.Init(ctx, &pbPlugin.Init_Request{
		Spec: spec.Spec,
	}); err != nil {
		return nil, err
	}

	writeClient, err := c.client.Write(ctx)
	if err != nil {
		return nil, err
	}

	if err := writeClient.Send(&pbPlugin.Write_Request{
		Message: &pbPlugin.Write_Request_MigrateTable{
			MigrateTable: &pbPlugin.MessageMigrateTable{
				Table: tableBytes,
			},
		},
	}); err != nil {
		return nil, err
	}

	syncClient, err := c.client.Sync(ctx, &pbPlugin.Sync_Request{
		Tables: []string{stateTablePrefix + name},
	})
	if err != nil {
		return nil, err
	}
	for {
		res, err := syncClient.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		insertMessage := res.GetInsert()
		if insertMessage == nil {
			return nil, fmt.Errorf("unexpected message type %T", res)
		}
		rdr, err := ipc.NewReader(bytes.NewReader(insertMessage.Record))
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
			keys := record.Columns()[0].(*array.String)
			values := record.Columns()[1].(*array.String)
			for i := 0; i < keys.Len(); i++ {
				c.mem[keys.Value(i)] = values.Value(i)
			}
		}
	}
	return c, nil
}

func (c *ClientV3) SetKey(ctx context.Context, key string, value string) error {
	c.mem[key] = value
	return nil
}

func (c *ClientV3) flush(ctx context.Context) error {
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, nil)
	for k, v := range c.mem {
		bldr.Field(0).(*array.StringBuilder).Append(k)
		bldr.Field(1).(*array.StringBuilder).Append(v)
	}
	rec := bldr.NewRecord()
	var buf bytes.Buffer
	wrtr := ipc.NewWriter(&buf, ipc.WithSchema(rec.Schema()))
	if err := wrtr.Write(rec); err != nil {
		return err
	}
	if err := wrtr.Close(); err != nil {
		return err
	}
	writeClient, err := c.client.Write(ctx)
	if err != nil {
		return err
	}
	if err := writeClient.Send(&pbPlugin.Write_Request{
		Message: &pbPlugin.Write_Request_Insert{
			Insert: &pbPlugin.MessageInsert{
				Record: buf.Bytes(),
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

func (c *ClientV3) GetKey(ctx context.Context, key string) (string, error) {
	if val, ok := c.mem[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}
