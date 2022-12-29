package schema

import "testing"

type SomeUUIDWrapper struct {
	SomeUUIDType
}

type SomeUUIDType [16]byte
type StringUUIDType string
type GetterUUIDType string

func (s StringUUIDType) String() string {
	return string(s)
}

func (s GetterUUIDType) Get() any {
	return string(s)
}

func TestUUIDSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result UUID
	}{
		{
			source: nil,
			result: UUID{Status: Null},
		},
		{
			source: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
		},
		{
			source: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
		},
		{
			source: SomeUUIDType{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
		},
		{
			source: ([]byte)(nil),
			result: UUID{Status: Null},
		},
		{
			source: "00010203-0405-0607-0809-0a0b0c0d0e0f",
			result: UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
		},
		{
			source: "000102030405060708090a0b0c0d0e0f",
			result: UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
		},
		{
			source: StringUUIDType("00010203-0405-0607-0809-0a0b0c0d0e0f"),
			result: UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
		}, {
			source: GetterUUIDType("00010203-0405-0607-0809-0a0b0c0d0e0f"),
			result: UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: Present},
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

func TestUUID_LessThan(t *testing.T) {
	uuid1 := [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	uuid2 := [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 16}

	cases := []struct {
		a    UUID
		b    UUID
		want bool
	}{
		{a: UUID{Bytes: uuid1, Status: Present}, b: UUID{Bytes: uuid2, Status: Present}, want: true},
		{a: UUID{Bytes: uuid2, Status: Present}, b: UUID{Bytes: uuid1, Status: Present}, want: false},
		{a: UUID{Bytes: uuid1, Status: Present}, b: UUID{Bytes: uuid1, Status: Present}, want: false},
		{a: UUID{Bytes: uuid1, Status: Null}, b: UUID{Bytes: uuid2, Status: Present}, want: true},
		{a: UUID{Bytes: uuid2, Status: Null}, b: UUID{Bytes: uuid1, Status: Present}, want: true},
		{a: UUID{Bytes: uuid1, Status: Null}, b: UUID{Bytes: uuid1, Status: Present}, want: true},
		{a: UUID{Bytes: uuid1, Status: Null}, b: UUID{Bytes: uuid2, Status: Null}, want: true},
	}

	for _, tt := range cases {
		if got := tt.a.LessThan(&tt.b); got != tt.want {
			t.Errorf("UUID.LessThan(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
