//nolint:revive,gocritic,unused
package schema

import (
	"fmt"

	"github.com/apache/arrow/go/v16/arrow"
)

const (
	noConversion                = "no conversion available"
	cannotDecodeString          = "cannot decode from string"
	cannotFindDimensions        = "cannot find dimensions"
	notInterface                = "not an interface"
	expectedElements            = "expected %d elements, but got %d instead"
	expectedElementsInDimension = "expected %d elements in dimension %d, but got %d instead"
	cannotSetIndex              = "cannot set index %d"
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
