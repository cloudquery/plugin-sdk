package scalar

import (
	"strconv"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/apache/arrow/go/v13/arrow/memory"
	"github.com/stretchr/testify/assert"
)

func TestStructMissingKeys(t *testing.T) {
	tl := []struct {
		dt          *arrow.StructType
		input       any
		expectPanic bool
	}{
		{dt: arrow.StructOf(arrow.Field{Name: "i64", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "s", Type: arrow.BinaryTypes.String}), input: `{"i64": 1, "s": "foo"}`, expectPanic: false},
		{dt: arrow.StructOf(arrow.Field{Name: "i64", Type: arrow.PrimitiveTypes.Int64}, arrow.Field{Name: "s", Type: arrow.BinaryTypes.String}), input: `{"i64": 1, "s": "foo", "extra":"bar"}`, expectPanic: true},
	}

	for idx, tc := range tl {
		tc := tc
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			panicked := false
			defer func() {
				if t.Failed() {
					return
				}

				if panicked && !tc.expectPanic {
					t.Errorf("unexpected panic")
				}
				if !panicked && tc.expectPanic {
					t.Errorf("expected panic")
				}
			}()

			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()

			bldr := array.NewBuilder(memory.DefaultAllocator, tc.dt)
			defer bldr.Release()

			s := NewScalar(tc.dt)
			if !arrow.TypeEqual(s.DataType(), tc.dt) {
				t.Fatalf("expected %v, got %v", tc.dt, s.DataType())
			}

			assert.NoError(t, s.Set(tc.input))
			if t.Failed() {
				return
			}

			assert.True(t, s.IsValid())
			AppendToBuilder(bldr, s)
		})
	}
}
