package state

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

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
	// changed keys. With checkpoints, each key has the scopes of its changes since the last flush, nil for changes made
	// outside table resolvers.
	changes map[string]*scopeSet
	// the last checkpoint of each key, which the next checkpoint of the key must not overtake
	last map[string]*checkpoint.Delivery
	// set when one of the checkpoints in sent is aborted
	abortSeen     atomic.Bool
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
		changes:       make(map[string]*scopeSet),
		last:          make(map[string]*checkpoint.Delivery),
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

	scopes, pending := c.changes[key]
	if c.mem[key].value == value {
		// the pending change summarizes this caller's resources too
		if pending && c.checkpoints != nil {
			scopes.add(checkpoint.FromContext(ctx))
		}
		return nil // don't update if the value is the same
	}
	c.mem[key] = versionedValue{
		value:   value,
		version: c.mem[key].version + 1,
	}
	if c.checkpoints == nil {
		c.changes[key] = nil
		return nil
	}
	if !pending {
		scopes = newScopeSet()
		c.changes[key] = scopes
	}
	scopes.add(checkpoint.FromContext(ctx))
	return nil
}

// minCompactSize is the size at which a scope set is first compacted.
const minCompactSize = 64

// scopeSet holds the scopes of a key's changes since the last flush.
type scopeSet struct {
	scopes map[checkpoint.Scope]struct{}
	// the size at which the set is compacted next
	compactAt int
}

func newScopeSet() *scopeSet {
	return &scopeSet{scopes: make(map[checkpoint.Scope]struct{}), compactAt: minCompactSize}
}

// add adds scope. Once the set reaches compactAt, it collapses the scopes of ended resolvers into the one they share,
// so that a key that many resolvers change keeps a small set. It compacts again only once the set doubles, so that a
// set of scopes that do not collapse costs constant time per add.
func (s *scopeSet) add(scope checkpoint.Scope) {
	s.scopes[scope] = struct{}{}
	if len(s.scopes) < s.compactAt {
		return
	}
	for scope := range s.scopes {
		if collapsed := checkpoint.Collapse(scope); collapsed != scope {
			delete(s.scopes, scope)
			s.scopes[collapsed] = struct{}{}
		}
	}
	s.compactAt = max(minCompactSize, 2*len(s.scopes))
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

	c.changes = make(map[string]*scopeSet)
	return nil
}

// flushCheckpoints sends one checkpoint per scope, so each waits only for the resources of the scopes that changed its
// keys. Keys whose checkpoint fails stay pending.
func (c *Client) flushCheckpoints(ctx context.Context) error {
	aborted := c.requeueAborted()
	groups := make(map[string]*checkpointGroup)
	for key, scopes := range c.changes {
		if _, ok := aborted[key]; ok {
			continue
		}
		scope, id := joinScopes(scopes.scopes)
		group, ok := groups[id]
		if !ok {
			group = &checkpointGroup{scope: scope}
			groups[id] = group
		}
		group.keys = append(group.keys, key)
	}

	var errs []error
	if len(aborted) > 0 {
		errs = append(errs, fmt.Errorf("%d keys: %w", len(aborted), checkpoint.ErrSyncAborted))
	}
	// keys changed outside table resolvers do not wait for any scope, and their group id is empty, so they go first
	for _, id := range slices.Sorted(maps.Keys(groups)) {
		group := groups[id]
		msg := &message.SyncCheckpoint{Record: c.record(slices.Values(group.keys))}
		sent, err := checkpoint.Schedule(ctx, c.checkpoints, group.scope, c.unsent(group.keys), msg)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		sent.NotifyAbort(&c.abortSeen)
		for _, key := range group.keys {
			delete(c.changes, key)
			c.last[key] = sent
		}
	}
	return errors.Join(errs...)
}

// checkpointGroup is the keys whose changes wait for the same scopes, so one checkpoint covers them.
type checkpointGroup struct {
	scope checkpoint.Scope
	keys  []string
}

// joinScopes joins scopes, and returns an id that is the same for sets that join the same scopes, and empty for none.
func joinScopes(scopes map[checkpoint.Scope]struct{}) (checkpoint.Scope, string) {
	unique := make(map[checkpoint.Scope]struct{}, len(scopes))
	for s := range scopes {
		if s = checkpoint.Collapse(s); s != nil {
			unique[s] = struct{}{}
		}
	}
	ids := make([]string, 0, len(unique))
	for s := range unique {
		ids = append(ids, fmt.Sprintf("%p", s))
	}
	slices.Sort(ids)
	return checkpoint.Join(slices.Collect(maps.Keys(unique))...), strings.Join(ids, ",")
}

// requeueAborted puts the keys whose last checkpoint was dropped by a cancelled sync back among the changes, and
// returns them: they stay pending, as the resources they summarize may have been dropped too. It also forgets the last
// checkpoints that were sent. It only looks once a checkpoint was aborted.
func (c *Client) requeueAborted() map[string]struct{} {
	if !c.abortSeen.Load() {
		return nil
	}
	aborted := make(map[string]struct{})
	for key, sent := range c.last {
		switch {
		case sent.Delivered():
			delete(c.last, key)
		case sent.Aborted():
			aborted[key] = struct{}{}
			if _, ok := c.changes[key]; !ok {
				c.changes[key] = newScopeSet()
			}
		}
	}
	return aborted
}

// unsent returns the previous checkpoints of keys that are not sent yet, as the scopes they wait for may still be
// sending the resources they cover.
func (c *Client) unsent(keys []string) []*checkpoint.Delivery {
	var unsent []*checkpoint.Delivery
	for _, key := range keys {
		sent, ok := c.last[key]
		if !ok {
			continue
		}
		if sent.Delivered() {
			delete(c.last, key)
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
