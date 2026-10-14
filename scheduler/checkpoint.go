package scheduler

import (
	"context"
	"maps"
	"slices"
	"sync"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/schema"

	"github.com/cloudquery/plugin-sdk/v4/internal/checkpoint"
	"github.com/cloudquery/plugin-sdk/v4/message"
)

// syncCheckpoints sends checkpoints on the sync message channel after the resources handed to the batcher before them.
// As a checkpoint.Scope, it holds checkpoints until the sync completes, for strategies that cannot track what a table
// resolver produced.
type syncCheckpoints struct {
	res      chan<- message.SyncMessage
	requests chan []message.SyncMessage
	// closed when Sync stops reading requests
	loopDone chan struct{}
	// closed once the batcher is closed
	finished chan struct{}

	mu       sync.Mutex
	deferred []message.SyncMessage
	// state rows held until the sync completes, by table, so repeated checkpoints of a key keep only the newest row
	deferredState map[string]*heldState
	done          bool
	completed     bool
}

var _ checkpoint.Scope = (*syncCheckpoints)(nil)

// state table columns, see the state client
const (
	stateKeyColumn     = "key"
	stateValueColumn   = "value"
	stateVersionColumn = "version"
)

func newSyncCheckpoints(res chan<- message.SyncMessage) *syncCheckpoints {
	return &syncCheckpoints{
		res:      res,
		requests: make(chan []message.SyncMessage),
		loopDone: make(chan struct{}),
		finished: make(chan struct{}),
	}
}

// send sends msgs after every resource Sync received so far.
func (c *syncCheckpoints) send(ctx context.Context, msgs []message.SyncMessage) error {
	if len(msgs) == 0 {
		return nil
	}
	select {
	case c.requests <- msgs:
		return nil
	case <-c.loopDone:
		return c.sendAfterSync(ctx, msgs)
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}

func (c *syncCheckpoints) sendAfterSync(ctx context.Context, msgs []message.SyncMessage) error {
	select {
	case <-c.finished:
	case <-ctx.Done():
		return context.Cause(ctx)
	}
	if !c.completed {
		return checkpoint.ErrSyncAborted
	}
	return sendMessages(ctx, c.res, msgs)
}

func (c *syncCheckpoints) Checkpoint(ctx context.Context, msg message.SyncMessage) error {
	c.mu.Lock()
	if !c.done {
		if !c.holdState(msg) {
			c.deferred = append(c.deferred, msg)
		}
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()
	return c.sendAfterSync(ctx, []message.SyncMessage{msg})
}

// complete is called once the batcher is closed after sending every resource.
func (c *syncCheckpoints) complete(ctx context.Context) error {
	c.mu.Lock()
	deferred := c.deferred
	for _, state := range c.deferredState {
		deferred = append(deferred, &message.SyncCheckpoint{Record: state.record()})
	}
	c.deferred = nil
	c.deferredState = nil
	c.done = true
	c.completed = true
	c.mu.Unlock()
	close(c.finished)
	return sendMessages(ctx, c.res, deferred)
}

// abort is called once the batcher is closed when some resources may not have been sent.
func (c *syncCheckpoints) abort() {
	c.mu.Lock()
	c.deferred = nil
	c.deferredState = nil
	c.done = true
	c.mu.Unlock()
	close(c.finished)
}

// tableScope sends checkpoints after the resources its table resolver produced before them, including their relations.
// The goroutine reading the resolver's results forwards the checkpoints once it has processed everything received
// before them.
type tableScope struct {
	checkpoints *syncCheckpoints
	notify      chan struct{}

	mu      sync.Mutex
	pending []message.SyncMessage
	ended   bool
}

var _ checkpoint.Scope = (*tableScope)(nil)

func newTableScope(checkpoints *syncCheckpoints) *tableScope {
	return &tableScope{checkpoints: checkpoints, notify: make(chan struct{}, 1)}
}

func (t *tableScope) Checkpoint(ctx context.Context, msg message.SyncMessage) error {
	t.mu.Lock()
	if t.ended {
		t.mu.Unlock()
		return t.checkpoints.send(ctx, []message.SyncMessage{msg})
	}
	t.pending = append(t.pending, msg)
	t.mu.Unlock()
	select {
	case t.notify <- struct{}{}:
	default:
	}
	return nil
}

// take returns the pending checkpoints.
func (t *tableScope) take() []message.SyncMessage {
	t.mu.Lock()
	defer t.mu.Unlock()
	pending := t.pending
	t.pending = nil
	return pending
}

// end returns the pending checkpoints. Later checkpoints go straight to the sync.
func (t *tableScope) end() []message.SyncMessage {
	t.mu.Lock()
	defer t.mu.Unlock()
	pending := t.pending
	t.pending = nil
	t.ended = true
	return pending
}

func sendMessages(ctx context.Context, res chan<- message.SyncMessage, msgs []message.SyncMessage) error {
	for _, msg := range msgs {
		if barrier, ok := msg.(*checkpoint.Barrier); ok {
			if err := barrier.Reached(ctx, res); err != nil {
				return err
			}
			continue
		}
		select {
		case res <- msg:
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
	return nil
}

type stateRow struct {
	value   string
	version uint64
}

type heldState struct {
	schema *arrow.Schema
	rows   map[string]stateRow
}

// holdState merges a state checkpoint into the held state, and reports false for other messages.
// Must be called with c.mu held.
func (c *syncCheckpoints) holdState(msg message.SyncMessage) bool {
	m, ok := msg.(*message.SyncCheckpoint)
	if !ok {
		return false
	}
	tableName, _ := m.Record.Schema().Metadata().GetValue(schema.MetadataTableName)
	keys, keysOK := column[*array.String](m.Record, stateKeyColumn)
	values, valuesOK := column[*array.String](m.Record, stateValueColumn)
	versions, versioned := column[*array.Uint64](m.Record, stateVersionColumn)
	if !keysOK || !valuesOK || (!versioned && m.Record.Schema().HasField(stateVersionColumn)) {
		return false
	}
	state := c.deferredState[tableName]
	if state == nil || !state.schema.Equal(m.Record.Schema()) {
		if state != nil {
			c.deferred = append(c.deferred, &message.SyncCheckpoint{Record: state.record()})
		}
		state = &heldState{schema: m.Record.Schema(), rows: make(map[string]stateRow)}
		if c.deferredState == nil {
			c.deferredState = make(map[string]*heldState)
		}
		c.deferredState[tableName] = state
	}
	for i := range keys.Len() {
		row := stateRow{value: values.Value(i)}
		if versioned {
			row.version = versions.Value(i)
			if held, ok := state.rows[keys.Value(i)]; ok && held.version > row.version {
				continue
			}
		}
		state.rows[keys.Value(i)] = row
	}
	return true
}

func (s *heldState) record() arrow.RecordBatch {
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, s.schema)
	defer bldr.Release()
	for _, key := range slices.Sorted(maps.Keys(s.rows)) {
		row := s.rows[key]
		for i, field := range s.schema.Fields() {
			switch field.Name {
			case stateKeyColumn:
				bldr.Field(i).(*array.StringBuilder).Append(key)
			case stateValueColumn:
				bldr.Field(i).(*array.StringBuilder).Append(row.value)
			case stateVersionColumn:
				bldr.Field(i).(*array.Uint64Builder).Append(row.version)
			default:
				bldr.Field(i).AppendNull()
			}
		}
	}
	return bldr.NewRecordBatch()
}

func column[T arrow.Array](record arrow.RecordBatch, name string) (T, bool) {
	var zero T
	indices := record.Schema().FieldIndices(name)
	if len(indices) != 1 {
		return zero, false
	}
	col, ok := record.Column(indices[0]).(T)
	return col, ok
}
