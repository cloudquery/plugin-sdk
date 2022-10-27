package cqtypes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBoolSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result Bool
	}{
		{source: true, result: Bool{Bool: true, Status: Present}},
		{source: false, result: Bool{Bool: false, Status: Present}},
		{source: "true", result: Bool{Bool: true, Status: Present}},
		{source: "false", result: Bool{Bool: false, Status: Present}},
		{source: "t", result: Bool{Bool: true, Status: Present}},
		{source: "f", result: Bool{Bool: false, Status: Present}},
		{source: _bool(true), result: Bool{Bool: true, Status: Present}},
		{source: _bool(false), result: Bool{Bool: false, Status: Present}},
		{source: nil, result: Bool{Status: Null}},
	}

	for i, tt := range successfulTests {
		var r Bool
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if diff := cmp.Diff(r, tt.result); diff != "" {
			t.Errorf("%d: got diff:\n%s", i, diff)
		}
	}
}
