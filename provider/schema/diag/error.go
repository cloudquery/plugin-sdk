package diag

// ExecutionError is a generic error returned when execution is run, ExecutionError satisfies
type ExecutionError struct {
	// Err is the underlying go error this diagnostic wraps
	Err error

	// Resource indicates the resource that failed in the execution
	resource string

	// Severity indicates the level of the Diagnostic. Currently, can be set to
	// either Error/Warning/Ignore
	severity Severity

	// Summary is a short description of the problem
	summary string

	// Detail is an optional second message, typically used to communicate a potential fix to the user.
	detail string

	// DiagnosticType indicates the classification family of this diagnostic
	diagnosticType DiagnosticType
}

func (e ExecutionError) Severity() Severity {
	return e.severity
}

func (e ExecutionError) Description() Description {
	return Description{
		e.resource,
		e.summary,
		e.detail,
	}
}

func (e ExecutionError) Type() DiagnosticType {
	return e.diagnosticType
}

func (e ExecutionError) Error() string {
	return e.Err.Error()
}

// FromError creates an ExecutionError from given error
func FromError(err error, severity Severity, dt DiagnosticType, resource, summary, details string) *ExecutionError {
	return &ExecutionError{
		Err:            err,
		severity:       severity,
		resource:       resource,
		summary:        summary,
		detail:         details,
		diagnosticType: dt,
	}
}
