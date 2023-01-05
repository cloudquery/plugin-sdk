package schema

import (
	"testing"
)

func TestTextArraySet(t *testing.T) {
	successfulTests := []struct {
		source any
		result TextArray
	}{
		{
			source: []string{"foo"},
			result: TextArray{
				Elements:   []Text{{Str: "foo", Status: Present}},
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
				Elements:   []Text{{Str: "foo", Status: Present}, {Str: "bar", Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [][][][]string{{{{"foo", "bar", "baz"}}}, {{{"wibble", "wobble", "wubble"}}}},
			result: TextArray{
				Elements: []Text{
					{Str: "foo", Status: Present},
					{Str: "bar", Status: Present},
					{Str: "baz", Status: Present},
					{Str: "wibble", Status: Present},
					{Str: "wobble", Status: Present},
					{Str: "wubble", Status: Present}},
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
				Elements:   []Text{{Str: "foo", Status: Present}, {Str: "bar", Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [2][1][1][3]string{{{{"foo", "bar", "baz"}}}, {{{"wibble", "wobble", "wubble"}}}},
			result: TextArray{
				Elements: []Text{
					{Str: "foo", Status: Present},
					{Str: "bar", Status: Present},
					{Str: "baz", Status: Present},
					{Str: "wibble", Status: Present},
					{Str: "wobble", Status: Present},
					{Str: "wubble", Status: Present}},
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

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

func TestTextArray_Size(t *testing.T) {
	cases := []struct {
		source TextArray
		result int
	}{
		{
			source: TextArray{
				Elements:   []Text{{Str: "foo", Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
			result: 3,
		},
		{
			source: TextArray{Status: Null},
			result: 0,
		},
		{
			source: TextArray{
				Elements:   []Text{{Str: "foo", Status: Present}, {Str: "bar", Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
			result: 6,
		},
		{
			source: TextArray{
				Elements: []Text{
					{Str: "foo", Status: Present},
					{Str: "bar", Status: Present},
					{Str: "baz", Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
			result: 9,
		},
	}

	for i, tt := range cases {
		result := tt.source.Size()
		if result != tt.result {
			t.Errorf("%d: %v != %v", i, result, tt.result)
		}
	}
}
