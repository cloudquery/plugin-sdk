package scalar

import "testing"

type SomeUUIDWrapper struct {
	SomeUUIDType
}

type SomeUUIDType [16]byte
type StringUUIDType string

func (s StringUUIDType) String() string {
	return string(s)
}

func TestUUIDSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result UUID
	}{
		{
			source: nil,
			result: UUID{},
		},
		{
			source: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: UUID{Value: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Valid: true},
		},
		{
			source: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: UUID{Value: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Valid: true},
		},
		{
			source: SomeUUIDType{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: UUID{Value: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Valid: true},
		},
		{
			source: ([]byte)(nil),
			result: UUID{},
		},
		{
			source: "00010203-0405-0607-0809-0a0b0c0d0e0f",
			result: UUID{Value: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Valid: true},
		},
		{
			source: "000102030405060708090a0b0c0d0e0f",
			result: UUID{Value: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Valid: true},
		},
		{
			source: StringUUIDType("00010203-0405-0607-0809-0a0b0c0d0e0f"),
			result: UUID{Value: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Valid: true},
		},
		{
			source: &UUID{Value: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Valid: true},
			result: UUID{Value: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Valid: true},
		},
	}

	for i, tt := range successfulTests {
		var r UUID
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
