package scalar

import (
	"testing"
)

func TestBoolSet(t *testing.T) {
	var nilPointerBool *bool
	var nilPointerString *string

	successfulTests := []struct {
		source any
		result Bool
	}{
		{source: true, result: Bool{Value: true, Valid: true}},
		{source: false, result: Bool{Value: false, Valid: true}},
		{source: "true", result: Bool{Value: true, Valid: true}},
		{source: "false", result: Bool{Value: false, Valid: true}},
		{source: "t", result: Bool{Value: true, Valid: true}},
		{source: "f", result: Bool{Value: false, Valid: true}},
		{source: _bool(true), result: Bool{Value: true, Valid: true}},
		{source: _bool(false), result: Bool{Value: false, Valid: true}},
		{source: &Bool{Value: true, Valid: true}, result: Bool{Value: true, Valid: true}},
		{source: nil, result: Bool{}},
		{source: nilPointerBool, result: Bool{}},
		{source: nilPointerString, result: Bool{}},
	}

	for i, tt := range successfulTests {
		var r Bool
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
