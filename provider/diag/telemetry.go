package diag

func TelemetryFromError(err error, eventType string, opts ...BaseErrorOption) Diagnostic {
	opts = append([]BaseErrorOption{WithSeverity(IGNORE), WithDetails(eventType)}, opts...)
	return NewBaseError(err, TELEMETRY, opts...)
}
