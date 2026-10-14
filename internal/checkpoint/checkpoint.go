// Package checkpoint lets state clients order their checkpoints after the resources they summarize.
package checkpoint

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/message"
)

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
