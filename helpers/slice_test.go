package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceSlice(t *testing.T) {
	someStringPtr := "test"
	cases := []struct {
		Name  string
		Value interface{}
		Want  []interface{}
	}{
		{Name: "base", Value: []string{"a", "b", "c"}, Want: []interface{}{"a", "b", "c"}},
		{Name: "nil", Value: nil, Want: nil},
		{Name: "empty", Value: []interface{}{}, Want: []interface{}{}},
		{Name: "empty_string_array", Value: []string{}, Want: []interface{}{}},
		{Name: "string_ptr_array", Value: []*string{&someStringPtr}, Want: []interface{}{&someStringPtr}},
		{Name: "string_array_ptr", Value: &[]string{"a"}, Want: []interface{}{"a"}},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Want, InterfaceSlice(tc.Value))
		})
	}
}
