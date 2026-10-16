package scheduler

import (
	"context"
	"errors"
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
	// the scope of every table resolver that ended: they all send checkpoints after the resources Sync received
	ended *endedScope

	mu       sync.Mutex
	deferred []message.SyncMessage
	// state rows held until the sync completes, by table, so repeated checkpoints of a key keep only the newest row
	deferredState map[string]*heldState
	// sent once the held state is sent, to track the checkpoints merged into it
	heldSent *checkpoint.Sent
	// table scopes that did not end, whose pending checkpoints an aborted sync drops
	scopes    map[*tableScope]struct{}
	done      bool
	completed bool

	// set when the caller closes res once Sync returns, as SyncAll does
	closeRes bool
	// held to send on res after Sync, so that Sync returns only once those sends are over
	resMu     sync.RWMutex
	resClosed bool
}

var (
	_ checkpoint.Scope    = (*syncCheckpoints)(nil)
	_ checkpoint.Covering = (*syncCheckpoints)(nil)
)

// state table columns, see the state client
const (
	stateKeyColumn     = "key"
	stateValueColumn   = "value"
	stateVersionColumn = "version"
)

func newSyncCheckpoints(res chan<- message.SyncMessage) *syncCheckpoints {
	c := &syncCheckpoints{
		res:      res,
		requests: make(chan []message.SyncMessage),
		loopDone: make(chan struct{}),
		finished: make(chan struct{}),
	}
	c.ended = &endedScope{checkpoints: c}
	return c
}

// send sends msgs after every resource Sync received so far. It drops msgs if it fails.
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
		dropMessages(msgs)
		return context.Cause(ctx)
	}
}

func (c *syncCheckpoints) sendAfterSync(ctx context.Context, msgs []message.SyncMessage) error {
	select {
	case <-c.finished:
	case <-ctx.Done():
		dropMessages(msgs)
		return context.Cause(ctx)
	}
	c.resMu.RLock()
	defer c.resMu.RUnlock()
	switch {
	case !c.completed:
		dropMessages(msgs)
		return checkpoint.ErrSyncAborted
	case c.resClosed:
		dropMessages(msgs)
		return errResClosed
	}
	return sendMessages(ctx, c.res, msgs)
}

var errResClosed = errors.New("the sync message channel is closed, as the sync returned")

// returned is called when Sync returns. It waits for checkpoints being sent, if the caller then closes res.
func (c *syncCheckpoints) returned() {
	if !c.closeRes {
		return
	}
	c.resMu.Lock()
	defer c.resMu.Unlock()
	c.resClosed = true
}

// endedScope is the scope of table resolvers that ended.
type endedScope struct {
	checkpoints *syncCheckpoints
}

func (e *endedScope) Checkpoint(ctx context.Context, msg message.SyncMessage) error {
	return e.checkpoints.send(ctx, []message.SyncMessage{msg})
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

// Covers reports whether sent is merged into the held state, which a checkpoint passed now is merged into too, or sent
// after.
func (c *syncCheckpoints) Covers(sent *checkpoint.Sent) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return !c.done && c.heldSent != nil && sent.MergedInto(c.heldSent)
}

// complete is called once the batcher is closed after sending every resource.
func (c *syncCheckpoints) complete(ctx context.Context) error {
	c.mu.Lock()
	deferred := c.deferred
	for _, state := range c.deferredState {
		deferred = append(deferred, &message.SyncCheckpoint{Record: state.record()})
	}
	heldSent := c.heldSent
	c.deferred = nil
	c.deferredState = nil
	c.heldSent = nil
	c.done = true
	c.completed = true
	c.mu.Unlock()
	close(c.finished)
	if err := sendMessages(ctx, c.res, deferred); err != nil {
		if heldSent != nil {
			heldSent.Abort()
		}
		return err
	}
	if heldSent == nil {
		return nil
	}
	return heldSent.Finish(ctx, c.res)
}

// abort is called once the batcher is closed when some resources may not have been sent.
func (c *syncCheckpoints) abort() {
	c.mu.Lock()
	deferred, heldSent, scopes := c.deferred, c.heldSent, c.scopes
	c.deferred = nil
	c.deferredState = nil
	c.heldSent = nil
	c.scopes = nil
	c.done = true
	c.mu.Unlock()
	close(c.finished)
	dropMessages(deferred)
	// a resolver goroutine that is stuck after the cancellation never forwards its pending checkpoints
	for scope := range scopes {
		dropMessages(scope.end())
	}
	if heldSent != nil {
		heldSent.Abort()
	}
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

var (
	_ checkpoint.Scope      = (*tableScope)(nil)
	_ checkpoint.Collapsing = (*tableScope)(nil)
)

func newTableScope(checkpoints *syncCheckpoints) *tableScope {
	t := &tableScope{checkpoints: checkpoints, notify: make(chan struct{}, 1)}
	checkpoints.mu.Lock()
	defer checkpoints.mu.Unlock()
	if checkpoints.done {
		// nothing drops the checkpoints of a scope that starts after the sync, so they go straight to the sync
		t.ended = true
		return t
	}
	if checkpoints.scopes == nil {
		checkpoints.scopes = make(map[*tableScope]struct{})
	}
	checkpoints.scopes[t] = struct{}{}
	return t
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

// Collapsed returns the scope shared by all ended table resolvers once the resolver ended.
func (t *tableScope) Collapsed() checkpoint.Scope {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ended {
		return t.checkpoints.ended
	}
	return t
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
	pending := t.pending
	t.pending = nil
	t.ended = true
	t.mu.Unlock()
	t.checkpoints.mu.Lock()
	delete(t.checkpoints.scopes, t)
	t.checkpoints.mu.Unlock()
	return pending
}

// sendMessages delivers msgs on res, and drops the ones it cannot deliver.
func sendMessages(ctx context.Context, res chan<- message.SyncMessage, msgs []message.SyncMessage) error {
	for i, msg := range msgs {
		if err := checkpoint.Deliver(ctx, res, msg); err != nil {
			dropMessages(msgs[i+1:])
			return err
		}
	}
	return nil
}

func dropMessages(msgs []message.SyncMessage) {
	for _, msg := range msgs {
		checkpoint.Drop(msg)
	}
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
	tracked, isTracked := msg.(*checkpoint.Tracked)
	if isTracked {
		msg = tracked.Message
	}
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
	if isTracked {
		if c.heldSent == nil {
			c.heldSent = checkpoint.NewHeld(c)
		}
		tracked.Sent.MergeInto(c.heldSent)
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
