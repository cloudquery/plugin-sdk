package diag

// nolint:revive
type Severity int

// nolint:revive
type Type int
type Diagnostic interface {
	error
	Severity() Severity
	Type() Type
	Description() Description
}

type Description struct {
	Resource   string
	ResourceID []string

	Summary string
	Detail  string
}

const (
	UNKNOWN Type = iota
	RESOLVING
	ACCESS
	THROTTLE
	DATABASE
	SCHEMA
	INTERNAL
	USER
)

const (
	// IGNORE severity is set for diagnostics that were ignored by the SDK
	IGNORE Severity = iota
	// WARNING severity are diagnostics that should be fixed but aren't fatal to the fetch execution
	WARNING
	// ERROR severity are diagnostics that were fatal in the fetch execution and should be fixed.
	ERROR
	// PANIC severity are diagnostics that are returned from a panic in the underlying code.
	PANIC
)

func (s Severity) String() string {
	switch s {
	case IGNORE:
		return "Ignore"
	case WARNING:
		return "Warning"
	case ERROR:
		return "Error"
	case PANIC:
		return "Panic"
	default:
		return "Unknown"
	}
}

func (d Type) String() string {
	switch d {
	case RESOLVING:
		return "Resolving"
	case ACCESS:
		return "Access"
	case THROTTLE:
		return "Throttle"
	case DATABASE:
		return "Database"
	case USER:
		return "User"
	case INTERNAL:
		return "Internal"
	case UNKNOWN:
		fallthrough
	default:
		return "UNKNOWN"
	}
}
