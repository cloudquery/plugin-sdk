package schema

import (
	"testing"
)

func TestUUIDArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result UUIDArray
	}{
		{
			source: nil,
			result: UUIDArray{Status: Null},
		},
		{
			source: [][16]byte{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
			result: UUIDArray{
				Elements:   []UUID{{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [][16]byte{},
			result: UUIDArray{Status: Present},
		},
		{
			source: ([][16]byte)(nil),
			result: UUIDArray{Status: Null},
		},
		{
			source: [][]byte{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
			result: UUIDArray{
				Elements:   []UUID{{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [][]byte{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
				nil,
				{32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47},
			},
			result: UUIDArray{
				Elements: []UUID{
					{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
					{Bytes: [16]byte{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}, Status: Present},
					{Status: Null},
					{Bytes: [16]byte{32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47}, Status: Present},
				},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 4}},
				Status:     Present},
		},
		{
			source: [][]byte{},
			result: UUIDArray{Status: Present},
		},
		{
			source: ([][]byte)(nil),
			result: UUIDArray{Status: Null},
		},
		{
			source: []string{"00010203-0405-0607-0809-0a0b0c0d0e0f"},
			result: UUIDArray{
				Elements:   []UUID{{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: []string{},
			result: UUIDArray{Status: Present},
		},
		{
			source: ([]string)(nil),
			result: UUIDArray{Status: Null},
		},
		{
			source: [][][16]byte{{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				{{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}}},
			result: UUIDArray{
				Elements: []UUID{
					{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
					{Bytes: [16]byte{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [][][][]string{
				{{{
					"00010203-0405-0607-0809-0a0b0c0d0e0f",
					"10111213-1415-1617-1819-1a1b1c1d1e1f",
					"20212223-2425-2627-2829-2a2b2c2d2e2f"}}},
				{{{
					"30313233-3435-3637-3839-3a3b3c3d3e3f",
					"40414243-4445-4647-4849-4a4b4c4d4e4f",
					"50515253-5455-5657-5859-5a5b5c5d5e5f"}}}},
			result: UUIDArray{
				Elements: []UUID{
					{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
					{Bytes: [16]byte{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}, Status: Present},
					{Bytes: [16]byte{32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47}, Status: Present},
					{Bytes: [16]byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63}, Status: Present},
					{Bytes: [16]byte{64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79}, Status: Present},
					{Bytes: [16]byte{80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95}, Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
		{
			source: [2][1][16]byte{{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
				{{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}}},
			result: UUIDArray{
				Elements: []UUID{
					{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
					{Bytes: [16]byte{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}, Status: Present}},
				Dimensions: []ArrayDimension{{LowerBound: 1, Length: 2}, {LowerBound: 1, Length: 1}},
				Status:     Present},
		},
		{
			source: [2][1][1][3]string{
				{{{
					"00010203-0405-0607-0809-0a0b0c0d0e0f",
					"10111213-1415-1617-1819-1a1b1c1d1e1f",
					"20212223-2425-2627-2829-2a2b2c2d2e2f"}}},
				{{{
					"30313233-3435-3637-3839-3a3b3c3d3e3f",
					"40414243-4445-4647-4849-4a4b4c4d4e4f",
					"50515253-5455-5657-5859-5a5b5c5d5e5f"}}}},
			result: UUIDArray{
				Elements: []UUID{
					{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
					{Bytes: [16]byte{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}, Status: Present},
					{Bytes: [16]byte{32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47}, Status: Present},
					{Bytes: [16]byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63}, Status: Present},
					{Bytes: [16]byte{64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79}, Status: Present},
					{Bytes: [16]byte{80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95}, Status: Present}},
				Dimensions: []ArrayDimension{
					{LowerBound: 1, Length: 2},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 1},
					{LowerBound: 1, Length: 3}},
				Status: Present},
		},
	}

	for i, tt := range successfulTests {
		var r UUIDArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
