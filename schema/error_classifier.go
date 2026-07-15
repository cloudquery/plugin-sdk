package schema

import "context"

// ErrorPhase identifies the resolver stage that produced an error.
type ErrorPhase int

const (
	ErrorPhaseTableResolver ErrorPhase = iota
	ErrorPhasePreResourceChunkResolver
	ErrorPhasePreResourceResolver
	ErrorPhaseColumnResolver
	ErrorPhasePostResourceResolver
)

func (p ErrorPhase) String() string {
	switch p {
	case ErrorPhaseTableResolver:
		return "table_resolver"
	case ErrorPhasePreResourceChunkResolver:
		return "pre_resource_chunk_resolver"
	case ErrorPhasePreResourceResolver:
		return "pre_resource_resolver"
	case ErrorPhaseColumnResolver:
		return "column_resolver"
	case ErrorPhasePostResourceResolver:
		return "post_resource_resolver"
	default:
		return "unknown"
	}
}

// ErrorEvent describes the context in which a resolver error occurred.
type ErrorEvent struct {
	Table  *Table
	Client ClientMeta
	Phase  ErrorPhase
	// Column is set only when Phase == ErrorPhaseColumnResolver.
	Column *Column
}

// ErrorClassifier reports whether a resolver error should be suppressed rather than
// raised. Suppressed errors are logged at debug level and are not counted in error
// metrics or emitted as a SyncError message. A nil ErrorClassifier raises every error.
// It is not consulted for primary key calculation or validation errors.
type ErrorClassifier func(ctx context.Context, err error, event ErrorEvent) bool

// Suppress is a nil-safe call: a nil ErrorClassifier returns false.
func (c ErrorClassifier) Suppress(ctx context.Context, err error, event ErrorEvent) bool {
	if c == nil {
		return false
	}
	return c(ctx, err, event)
}
