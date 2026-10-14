// Package checkpoint lets state clients order their checkpoints after the resources they summarize.
package checkpoint

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// ErrSyncAborted is returned for checkpoints that were dropped because the sync was cancelled.
var ErrSyncAborted = errors.New("sync aborted before the checkpoint could be sent")

// Scope is the part of a sync that a table resolver drives: the resources it produces and their relations.
type Scope interface {
	// Checkpoint passes msg to Deliver after every resource the scope produced before the call is sent.
	// It may return before that.
	Checkpoint(ctx context.Context, msg message.SyncMessage) error
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
		if err := Deliver(ctx, res, m.Message); err != nil {
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

// Barrier travels through scopes and sent checkpoints, and calls reached each time one of them gets to it.
type Barrier struct {
	reached func(ctx context.Context, res chan<- message.SyncMessage) error
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

// Sent tracks whether a checkpoint was sent, so that later checkpoints of the same keys can wait for it.
type Sent struct {
	mu      sync.Mutex
	done    bool
	waiting []*Barrier
	// set when the checkpoint was merged into another one, which is then sent in its place
	mergedInto *Sent
}

// Done reports whether the checkpoint was sent.
func (s *Sent) Done() bool {
	s.mu.Lock()
	mergedInto, done := s.mergedInto, s.done
	s.mu.Unlock()
	if mergedInto != nil {
		return mergedInto.Done()
	}
	return done
}

// MergeInto marks s as sent once into is sent. into must not be sent yet.
func (s *Sent) MergeInto(into *Sent) {
	s.mu.Lock()
	waiting := s.waiting
	s.waiting = nil
	s.mergedInto = into
	s.mu.Unlock()
	for _, barrier := range waiting {
		into.wait(barrier)
	}
}

// wait makes barrier wait for s. It reports true, without waiting, if s was sent.
func (s *Sent) wait(barrier *Barrier) bool {
	s.mu.Lock()
	if s.mergedInto != nil {
		into := s.mergedInto
		s.mu.Unlock()
		return into.wait(barrier)
	}
	defer s.mu.Unlock()
	if s.done {
		return true
	}
	s.waiting = append(s.waiting, barrier)
	return false
}

// Finish marks s as sent, and releases the checkpoints waiting for it.
func (s *Sent) Finish(ctx context.Context, res chan<- message.SyncMessage) error {
	s.mu.Lock()
	s.done = true
	waiting := s.waiting
	s.waiting = nil
	s.mu.Unlock()
	var errs []error
	for _, barrier := range waiting {
		if err := barrier.reached(ctx, res); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Send passes msg to Deliver once every one of scope's resources produced before the call is sent, and every
// checkpoint in after is sent. Without a scope or pending checkpoints to wait for, it sends msg on res.
// The returned Sent tracks msg.
func Send(ctx context.Context, res chan<- message.SyncMessage, scope Scope, after []*Sent, msg message.SyncMessage) (*Sent, error) {
	sent := &Sent{}
	tracked := &Tracked{Message: msg, Sent: sent}
	scopes := flatten(scope)
	if len(after) == 0 {
		switch len(scopes) {
		case 0:
			return sent, Deliver(ctx, res, tracked)
		case 1:
			return sent, scopes[0].Checkpoint(ctx, tracked)
		}
	}

	var remaining atomic.Int64
	remaining.Store(int64(len(scopes) + len(after)))
	barrier := &Barrier{reached: func(ctx context.Context, res chan<- message.SyncMessage) error {
		if remaining.Add(-1) > 0 {
			return nil
		}
		return Deliver(ctx, res, tracked)
	}}
	var errs []error
	for _, previous := range after {
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
	return sent, errors.Join(errs...)
}

// Join returns a scope whose checkpoints follow the resources of all the given scopes. It skips nil scopes, and
// returns nil if all are nil.
func Join(scopes ...Scope) Scope {
	var unique []Scope
	for _, scope := range scopes {
		for _, s := range flatten(scope) {
			if !contains(unique, s) {
				unique = append(unique, s)
			}
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

func contains(scopes []Scope, scope Scope) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}
