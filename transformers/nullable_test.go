package transformers_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/cloudquery/plugin-sdk/v3/transformers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type (
	simple struct {
		A, B int
	}
	pointer            *simple
	pointerAlias       = pointer
	simplePointerAlias *simple

	intPtr *int
)

func TestNullable(t *testing.T) {
	cases := []struct {
		name     string
		value    any
		nullable bool
	}{
		{name: "array", value: uuid.UUID{}, nullable: false},
		{name: "chan", value: make(chan int), nullable: true},
		{name: "func", value: transformers.Nullable, nullable: true},
		{name: "interface", value: io.Reader(nil), nullable: true},
		{name: "map", value: map[string]struct{}{}, nullable: true},
		{name: "simple", value: simple{}, nullable: false},
		{name: "new(simple)", value: new(simple), nullable: true},
		{name: "slice", value: []simple{}, nullable: true},
		{name: "intPtr", value: intPtr(nil), nullable: true},
		{name: "pointer", value: pointer(nil), nullable: true},
		{name: "pointerAlias", value: pointerAlias(nil), nullable: true},
		{name: "simplePointerAlias", value: simplePointerAlias(nil), nullable: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			require.Exactly(t, tc.nullable, transformers.Nullable(reflect.TypeOf(tc.value)))
		})
	}
}
