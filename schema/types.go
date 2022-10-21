package schema

import (
	"fmt"
	"strings"
)

type CQType interface{
	Type() ValueType
	Scan(v interface{}) error
	Equal(CQType) bool
}

type CQTypes []CQType

func (c CQTypes) Len() int {
	return len(c)
}

func (c CQTypes) Equal(other CQTypes) bool {
	if other == nil {
		return false
	}
	if len(c) != len(other) {
		return false
	}
	for i := range c {
		if c[i] == nil {
			if other[i] != nil {
				return false
			}
		} else {
			if !c[i].Equal(other[i]) {
				return false
			}
		}

	}
	return true
}

func (c CQTypes) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range c {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteString("]")
	return sb.String()
}