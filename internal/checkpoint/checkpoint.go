// Package checkpoint lets state clients order their checkpoints after the resources they summarize.
package checkpoint

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

// ErrSyncAborted is returned for checkpoints that were dropped because the sync was cancelled.
var ErrSyncAborted = errors.New("sync aborted before the checkpoint could be sent")

// Scope is the part of a sync that a table resolver drives: the resources it produces and their relations.
type Scope interface {
	// Checkpoint sends msg on the sync message channel after every resource the scope produced before the call.
	// It may return before msg is sent.
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

// Barrier travels through scopes like a checkpoint. Where a checkpoint would be sent on the sync message channel,
// the scheduler calls Reached instead.
type Barrier struct {
	reached func(ctx context.Context, res chan<- message.SyncMessage) error
}

func (*Barrier) GetTable() *schema.Table { return nil }
func (*Barrier) IsSyncMessage() bool     { return true }

// Reached is called once every resource sent before the barrier is on res.
func (b *Barrier) Reached(ctx context.Context, res chan<- message.SyncMessage) error {
	return b.reached(ctx, res)
}

// Join returns a scope whose checkpoints follow the resources of all the given scopes. It skips nil scopes, and
// returns nil if all are nil.
func Join(scopes ...Scope) Scope {
	var unique []Scope
	for _, scope := range scopes {
		if j, ok := scope.(*joined); ok {
			unique = appendUnique(unique, j.scopes...)
		} else if scope != nil {
			unique = appendUnique(unique, scope)
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

func appendUnique(scopes []Scope, add ...Scope) []Scope {
	for _, scope := range add {
		found := false
		for _, s := range scopes {
			if s == scope {
				found = true
				break
			}
		}
		if !found {
			scopes = append(scopes, scope)
		}
	}
	return scopes
}

type joined struct {
	scopes []Scope
}

// Checkpoint passes a barrier through every scope, and sends msg when the last one reaches it.
func (j *joined) Checkpoint(ctx context.Context, msg message.SyncMessage) error {
	var remaining atomic.Int64
	remaining.Store(int64(len(j.scopes)))
	barrier := &Barrier{reached: func(ctx context.Context, res chan<- message.SyncMessage) error {
		if remaining.Add(-1) > 0 {
			return nil
		}
		select {
		case res <- msg:
			return nil
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}}
	var errs []error
	for _, scope := range j.scopes {
		if err := scope.Checkpoint(ctx, barrier); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
