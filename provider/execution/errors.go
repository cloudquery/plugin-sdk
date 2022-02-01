package execution

import (
	"fmt"
	"strings"

	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

const (
	// fdLimitMessage defines the message for when a client isn't able to fetch because the open fd limit is hit
	fdLimitMessage = "try increasing number of available file descriptors via `ulimit -n 10240` or by increasing timeout via provider specific parameters"
)

type ErrorClassifier func(meta schema.ClientMeta, resource string, err error) diag.Diagnostics

func defaultErrorClassifier(_ schema.ClientMeta, resource string, err error) diag.Diagnostics {
	if _, ok := err.(diag.Diagnostic); ok {
		return nil
	}
	if _, ok := err.(diag.Diagnostics); ok {
		return nil
	}
	if strings.Contains(err.Error(), ": socket: too many open files") {
		// Return a Diagnostic error so that it can be properly propagated back to the user via the CLI
		return FromError(err, WithResource(resource), WithSummary(fdLimitMessage), WithType(diag.THROTTLE), WithSeverity(diag.WARNING))
	}
	return nil
}

// Error is a generic error returned when execution is run, Error satisfies diagnostic interface
type Error struct {
	// Err is the underlying go error this diagnostic wraps
	err error

	// Resource indicates the resource that failed in the execution
	resource string

	// Severity indicates the level of the Diagnostic. Currently, can be set to
	// either Error/Warning/Ignore
	severity diag.Severity

	// Summary is a short description of the problem
	summary string

	// Detail is an optional second message, typically used to communicate a potential fix to the user.
	detail string

	// DiagnosticType indicates the classification family of this diagnostic
	diagnosticType diag.DiagnosticType
}

// NewError creates an execution Error
func NewError(severity diag.Severity, dt diag.DiagnosticType, resource, summary string, args ...interface{}) *Error {
	return &Error{
		err:            fmt.Errorf(summary, args...),
		severity:       severity,
		resource:       resource,
		summary:        fmt.Sprintf(summary, args...),
		detail:         "",
		diagnosticType: dt,
	}
}

func (e Error) Err() error {
	return e.err
}

func (e Error) Severity() diag.Severity {
	return e.severity
}

func (e Error) Description() diag.Description {
	return diag.Description{
		Resource: e.resource,
		Summary:  e.summary,
		Detail:   e.detail,
	}
}

func (e Error) Type() diag.DiagnosticType {
	return e.diagnosticType
}

func (e Error) Error() string {
	// return original error
	if e.err != nil {
		return e.err.Error()
	}
	return e.summary
}

type Option func(e *Error)

func WithSeverity(s diag.Severity) Option {
	return func(e *Error) {
		e.severity = s
	}
}

func WithType(dt diag.DiagnosticType) Option {
	return func(e *Error) {
		e.diagnosticType = dt
	}
}

func WithSummary(summary string, args ...interface{}) Option {
	return func(e *Error) {
		e.summary = fmt.Sprintf(summary, args...)
	}
}

func WithResource(resource string) Option {
	return func(e *Error) {
		e.resource = resource
	}
}

func WithDetails(detail string) Option {
	return func(e *Error) {
		e.detail = detail
	}
}

func WithErrorClassifier(e *Error) {
	if e.err != nil && strings.Contains(e.err.Error(), ": socket: too many open files") {
		// Return a Diagnostic error so that it can be properly propagated back to the user via the CLI
		e.severity = diag.WARNING
		e.summary = fdLimitMessage
		e.diagnosticType = diag.THROTTLE
	}
}

func FromError(err error, opts ...Option) diag.Diagnostics {
	switch ti := err.(type) {
	case diag.Diagnostic:
		return diag.Diagnostics{ti}
	case diag.Diagnostics:
		return ti
	default:
		e := &Error{
			err:            err,
			severity:       diag.ERROR,
			diagnosticType: diag.RESOLVING,
		}
		for _, o := range opts {
			o(e)
		}
		return diag.Diagnostics{e}
	}
}
