package diag

type nativeError struct {
	err error
}

func (nativeError) Severity() Severity {
	return ERROR
}

func (nativeError) Type() Type {
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
