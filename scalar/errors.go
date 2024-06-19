package scalar

import (
	"fmt"

	"github.com/apache/arrow/go/v16/arrow"
)

const (
	noConversion = "no conversion available"
)

type ValidationError struct {
	Err   error
	Msg   string
	Type  arrow.DataType
	Value any
}

func (e *ValidationError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("cannot set `%s` with value `%v`: %s", e.Type, e.Value, e.Msg)
	}
	return fmt.Sprintf("cannot set `%s` with value `%v`: %s (%s)", e.Type, e.Value, e.Msg, e.Err)
}

// this prints the error without the value
func (e *ValidationError) MaskedError() string {
	if e.Err == nil {
		return fmt.Sprintf("cannot set `%s`: %s", e.Type, e.Msg)
	}
	return fmt.Sprintf("cannot set `%s`: %s (%s)", e.Type, e.Msg, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}
