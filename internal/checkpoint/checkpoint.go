// Package checkpoint lets state clients order their checkpoints after the resources they summarize.
package checkpoint

import (
	"context"
	"errors"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// ErrSyncAborted is returned for checkpoints that were dropped because the sync was cancelled.
var ErrSyncAborted = errors.New("sync aborted before the checkpoint could be sent")

// Scope is the part of a sync that a table resolver drives: the resources it produces and their relations.
type Scope interface {
	// Checkpoint passes msg to Deliver after every resource the scope produced before the call is sent, or to Drop if
	// the sync is cancelled first. It may return before that.
	Checkpoint(ctx context.Context, msg message.SyncMessage) error
}

// Covering is implemented by scopes that never send a checkpoint before an unsent one they already hold.
type Covering interface {
	// Covers reports whether a checkpoint passed to the scope now will be sent no earlier than sent.
	Covers(sent *Sent) bool
}

// Collapsing is implemented by scopes that can become equivalent to a scope they share with others, such as a table
// resolver's scope once the resolver ends.
type Collapsing interface {
	// Collapsed returns the shared scope, or the scope itself.
	Collapsed() Scope
}

type scopeKey struct{}

func WithScope(ctx context.Context, scope Scope) context.Context {
	return context.WithValue(ctx, scopeKey{}, scope)
}

// FromContext returns the scope ctx belongs to, or nil outside table resolvers.
func FromContext(ctx context.Context) Scope {
	scope, _ := ctx.Value(scopeKey{}).(Scope)
	return scope
}

// Deliver sends msg on res. A *Barrier is marked as reached instead, and a *Tracked message is sent and marked as sent.
func Deliver(ctx context.Context, res chan<- message.SyncMessage, msg message.SyncMessage) error {
	switch m := msg.(type) {
	case *Barrier:
		return m.reached(ctx, res)
	case *Tracked:
		if m.Sent.Aborted() {
			return nil
		}
		if err := Deliver(ctx, res, m.Message); err != nil {
			m.Sent.Abort()
			return err
		}
		return m.Sent.Finish(ctx, res)
	default:
		select {
		case res <- msg:
			return nil
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}

// Drop marks msg as never to be sent: a dropped *Tracked message, or the message a dropped *Barrier leads to, is
// aborted, and so are the checkpoints waiting for it.
func Drop(msg message.SyncMessage) {
	switch m := msg.(type) {
	case *Barrier:
		m.abort()
	case *Tracked:
		m.Sent.Abort()
	}
}

// Barrier travels through scopes and sent checkpoints, and calls reached each time one of them gets to it.
type Barrier struct {
	reached func(ctx context.Context, res chan<- message.SyncMessage) error
	abort   func()
}

func (*Barrier) GetTable() *schema.Table { return nil }
func (*Barrier) IsSyncMessage() bool     { return true }

// Tracked is a message that marks Sent once Deliver sends it.
type Tracked struct {
	Message message.SyncMessage
	Sent    *Sent
}

func (*Tracked) GetTable() *schema.Table { return nil }
func (*Tracked) IsSyncMessage() bool     { return true }

// Sent tracks whether a checkpoint was sent or aborted, so that later checkpoints of the same keys can wait for it.
type Sent struct {
	mu      sync.Mutex
	done    bool
	aborted bool
	waiting []*Barrier
	// set when the checkpoint was merged into another one, which is then sent in its place
	mergedInto *Sent
	// for a checkpoint that others are merged into: the scope holding it
	heldBy Scope
	// set when the checkpoint is aborted
	abortFlags map[*atomic.Bool]struct{}
}

// NewHeld returns a Sent for a checkpoint that holder merges other checkpoints into.
func NewHeld(holder Scope) *Sent {
	return &Sent{heldBy: holder}
}

// root returns the checkpoint s is sent with.
func (s *Sent) root() *Sent {
	for {
		s.mu.Lock()
		into := s.mergedInto
		s.mu.Unlock()
		if into == nil {
			return s
		}
		s = into
	}
}

// lockRoot returns the checkpoint s is sent with, locked. It is still the root once locked, as MergeInto runs under the
// lock of the checkpoint it merges.
func (s *Sent) lockRoot() *Sent {
	for {
		r := s.root()
		r.mu.Lock()
		if r.mergedInto == nil {
			return r
		}
		r.mu.Unlock()
	}
}

// Done reports whether the checkpoint was sent.
func (s *Sent) Done() bool {
	r := s.lockRoot()
	defer r.mu.Unlock()
	return r.done
}

// Aborted reports whether the checkpoint was dropped.
func (s *Sent) Aborted() bool {
	r := s.lockRoot()
	defer r.mu.Unlock()
	return r.aborted
}

func (s *Sent) holder() Scope {
	r := s.lockRoot()
	defer r.mu.Unlock()
	return r.heldBy
}

// NotifyAbort sets flag if the checkpoint is aborted, or now if it already was.
func (s *Sent) NotifyAbort(flag *atomic.Bool) {
	r := s.lockRoot()
	defer r.mu.Unlock()
	if r.aborted {
		flag.Store(true)
		return
	}
	if r.abortFlags == nil {
		r.abortFlags = make(map[*atomic.Bool]struct{})
	}
	r.abortFlags[flag] = struct{}{}
}

// MergedInto reports whether s is sent with into.
func (s *Sent) MergedInto(into *Sent) bool {
	return s != into && s.root() == into
}

// MergeInto marks s as sent once into is sent, and aborted if into is aborted. into must not be sent yet.
func (s *Sent) MergeInto(into *Sent) {
	s.mu.Lock()
	waiting, flags := s.waiting, s.abortFlags
	s.waiting, s.abortFlags = nil, nil
	s.mergedInto = into
	s.mu.Unlock()
	for _, barrier := range waiting {
		into.wait(barrier)
	}
	for flag := range flags {
		into.NotifyAbort(flag)
	}
}

// wait makes barrier wait for s. It reports true, without waiting, if s was sent.
func (s *Sent) wait(barrier *Barrier) bool {
	r := s.lockRoot()
	switch {
	case r.done:
		r.mu.Unlock()
		return true
	case r.aborted:
		r.mu.Unlock()
		barrier.abort()
		return false
	default:
		r.waiting = append(r.waiting, barrier)
		r.mu.Unlock()
		return false
	}
}

// Finish marks s as sent, and releases the checkpoints waiting for it.
func (s *Sent) Finish(ctx context.Context, res chan<- message.SyncMessage) error {
	s.mu.Lock()
	if s.done || s.aborted {
		s.mu.Unlock()
		return nil
	}
	s.done = true
	waiting := s.waiting
	s.waiting, s.abortFlags = nil, nil
	s.mu.Unlock()
	var errs []error
	for _, barrier := range waiting {
		if err := barrier.reached(ctx, res); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Abort marks s as dropped, and aborts the checkpoints waiting for it.
func (s *Sent) Abort() {
	r := s.lockRoot()
	if r.done || r.aborted {
		r.mu.Unlock()
		return
	}
	r.aborted = true
	waiting, flags := r.waiting, r.abortFlags
	r.waiting, r.abortFlags = nil, nil
	r.mu.Unlock()
	for flag := range flags {
		flag.Store(true)
	}
	for _, barrier := range waiting {
		barrier.abort()
	}
}

// Send passes msg to Deliver once every one of scope's resources produced before the call is sent, and every
// checkpoint in after is sent. Without a scope or pending checkpoints to wait for, it sends msg on res.
// The returned Sent tracks msg. If a checkpoint in after was aborted, Send returns ErrSyncAborted.
func Send(ctx context.Context, res chan<- message.SyncMessage, scope Scope, after []*Sent, msg message.SyncMessage) (*Sent, error) {
	pending := make([]*Sent, 0, len(after))
	for _, previous := range after {
		switch {
		case previous.Aborted():
			return nil, ErrSyncAborted
		case !previous.Done():
			pending = append(pending, previous)
		}
	}
	if scope == nil {
		// waiting only for checkpoints one scope holds: going through that scope sends no earlier, and lets it merge msg
		scope = commonHolder(pending)
	}
	if covering, ok := scope.(Covering); ok {
		pending = slices.DeleteFunc(pending, covering.Covers)
	}

	sent := &Sent{}
	tracked := &Tracked{Message: msg, Sent: sent}
	scopes := flatten(scope)
	if len(pending) == 0 {
		switch len(scopes) {
		case 0:
			return sent, Deliver(ctx, res, tracked)
		case 1:
			return sent, scopes[0].Checkpoint(ctx, tracked)
		}
	}

	var remaining atomic.Int64
	remaining.Store(int64(len(scopes) + len(pending)))
	barrier := &Barrier{
		reached: func(ctx context.Context, res chan<- message.SyncMessage) error {
			if remaining.Add(-1) > 0 {
				return nil
			}
			return Deliver(ctx, res, tracked)
		},
		abort: sent.Abort,
	}
	var errs []error
	for _, previous := range pending {
		if previous.wait(barrier) {
			if err := barrier.reached(ctx, res); err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, s := range scopes {
		if err := s.Checkpoint(ctx, barrier); err != nil {
			errs = append(errs, err)
		}
	}
	if err := errors.Join(errs...); err != nil {
		sent.Abort()
		return nil, err
	}
	return sent, nil
}

func commonHolder(pending []*Sent) Scope {
	var holder Scope
	for _, sent := range pending {
		h := sent.holder()
		if h == nil || (holder != nil && h != holder) {
			return nil
		}
		holder = h
	}
	return holder
}

// Collapse returns the shared scope a Collapsing scope collapsed into, or scope itself.
func Collapse(scope Scope) Scope {
	if c, ok := scope.(Collapsing); ok {
		return c.Collapsed()
	}
	return scope
}

// Join returns a scope whose checkpoints follow the resources of all the given scopes. It skips nil scopes, merges
// scopes that collapsed into the same one, and returns nil if all are nil.
func Join(scopes ...Scope) Scope {
	seen := make(map[Scope]struct{}, len(scopes))
	var unique []Scope
	for _, scope := range scopes {
		for _, s := range flatten(scope) {
			s = Collapse(s)
			if _, ok := seen[s]; ok {
				continue
			}
			seen[s] = struct{}{}
			unique = append(unique, s)
		}
	}
	switch len(unique) {
	case 0:
		return nil
	case 1:
		return unique[0]
	default:
		return &joined{scopes: unique}
	}
}

type joined struct {
	scopes []Scope
}

func (j *joined) Checkpoint(ctx context.Context, msg message.SyncMessage) error {
	_, err := Send(ctx, nil, j, nil, msg)
	return err
}

func flatten(scope Scope) []Scope {
	switch s := scope.(type) {
	case nil:
		return nil
	case *joined:
		return s.scopes
	default:
		return []Scope{s}
	}
}
