package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceSlice(t *testing.T) {
	someStringPtr := "test"
	cases := []struct {
		Name  string
		Value any
		Want  []any
	}{
		{Name: "base", Value: []string{"a", "b", "c"}, Want: []any{"a", "b", "c"}},
		{Name: "nil", Value: nil, Want: nil},
		{Name: "empty", Value: []any{}, Want: []any{}},
		{Name: "empty_string_array", Value: []string{}, Want: []any{}},
		{Name: "string_ptr_array", Value: []*string{&someStringPtr}, Want: []any{&someStringPtr}},
		{Name: "string_array_ptr", Value: &[]string{"a"}, Want: []any{"a"}},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Want, InterfaceSlice(tc.Value))
		})
	}
}
