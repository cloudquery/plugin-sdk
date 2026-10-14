package scheduler

import (
	"context"
	"errors"
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/internal/checkpoint"
	"github.com/cloudquery/plugin-sdk/v4/message"
)

var errSyncAborted = errors.New("sync aborted before the checkpoint could be sent")

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

	mu        sync.Mutex
	deferred  []message.SyncMessage
	done      bool
	completed bool
}

var _ checkpoint.Scope = (*syncCheckpoints)(nil)

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
		return errSyncAborted
	}
	return sendMessages(ctx, c.res, msgs)
}

func (c *syncCheckpoints) Checkpoint(ctx context.Context, msg message.SyncMessage) error {
	c.mu.Lock()
	if !c.done {
		c.deferred = append(c.deferred, msg)
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
	c.deferred = nil
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
		select {
		case res <- msg:
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
	return nil
}
