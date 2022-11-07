package schema

import (
	"testing"
)

func TestInt8ArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result Int8Array
	}{
		{
			source: []int64{1},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: []int32{1},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: []int16{1},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: []int{1},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: []uint64{1},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: []uint32{1},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: []uint16{1},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: []uint{1},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: (([]int64)(nil)),
			result: Int8Array{Status: Null},
		},
		{
			source: [][]int64{{1}, {2}},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}, {Int: 2, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [][][][]int64{{{{1, 2, 3}}}, {{{4, 5, 6}}}},
			result: Int8Array{
				Elements: []Int8{
					{Int: 1, Status: Present},
					{Int: 2, Status: Present},
					{Int: 3, Status: Present},
					{Int: 4, Status: Present},
					{Int: 5, Status: Present},
					{Int: 6, Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
		{
			source: [2][1]int64{{1}, {2}},
			result: Int8Array{
				Elements:   []Int8{{Int: 1, Status: Present}, {Int: 2, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [2][1][1][3]int64{{{{1, 2, 3}}}, {{{4, 5, 6}}}},
			result: Int8Array{
				Elements: []Int8{
					{Int: 1, Status: Present},
					{Int: 2, Status: Present},
					{Int: 3, Status: Present},
					{Int: 4, Status: Present},
					{Int: 5, Status: Present},
					{Int: 6, Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
	}

	for i, tt := range successfulTests {
		var r Int8Array
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
