package schema

import "fmt"

type ValidationError struct {
	err   error
	msg   string
	Type  ValueType
	Value any
}

func (e *ValidationError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("cannot convert `%s`: %s", e.Type, e.msg)
	}
	return fmt.Sprintf("cannot convert `%s`: %s (%s)", e.Type, e.msg, e.err)
}

func (e *ValidationError) Unwrap() error {
	return e.err
}
