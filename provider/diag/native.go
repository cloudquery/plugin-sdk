package diag

type nativeError struct {
	err error
}

func (n nativeError) Severity() Severity {
	return ERROR
}

func (n nativeError) Type() DiagnosticType {
	return INTERNAL
}

func (n nativeError) Description() Description {
	return Description{
		Summary: n.err.Error(),
		Detail:  "",
	}
}
func (n nativeError) Err() error {
	return n.err
}

func (n nativeError) Error() string {
	return n.err.Error()
}
