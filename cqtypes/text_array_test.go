package cqtypes

import (
	"reflect"
	"testing"
)

func TestTextArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result TextArray
	}{
		{
			source: []string{"foo"},
			result: TextArray{
				Elements:   []Text{{String: "foo", Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: (([]string)(nil)),
			result: TextArray{Status: Null},
		},
		{
			source: [][]string{{"foo"}, {"bar"}},
			result: TextArray{
				Elements:   []Text{{String: "foo", Status: Present}, {String: "bar", Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [][][][]string{{{{"foo", "bar", "baz"}}}, {{{"wibble", "wobble", "wubble"}}}},
			result: TextArray{
				Elements: []Text{
					{String: "foo", Status: Present},
					{String: "bar", Status: Present},
					{String: "baz", Status: Present},
					{String: "wibble", Status: Present},
					{String: "wobble", Status: Present},
					{String: "wubble", Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
		{
			source: [2][1]string{{"foo"}, {"bar"}},
			result: TextArray{
				Elements:   []Text{{String: "foo", Status: Present}, {String: "bar", Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [2][1][1][3]string{{{{"foo", "bar", "baz"}}}, {{{"wibble", "wobble", "wubble"}}}},
			result: TextArray{
				Elements: []Text{
					{String: "foo", Status: Present},
					{String: "bar", Status: Present},
					{String: "baz", Status: Present},
					{String: "wibble", Status: Present},
					{String: "wobble", Status: Present},
					{String: "wubble", Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
	}

	for i, tt := range successfulTests {
		var r TextArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}
