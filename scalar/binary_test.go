package scalar

import "testing"

func TestBinarySet(t *testing.T) {
	var nilPointerByteArray *[]byte
	var nilPointerString *string

	successfulTests := []struct {
		source any
		result Binary
	}{
		{source: []byte{1, 2, 3}, result: Binary{Value: []byte{1, 2, 3}, Valid: true}},
		{source: []byte{}, result: Binary{Value: []byte{}, Valid: true}},
		{source: []byte(nil), result: Binary{}},
		{source: _byteSlice{1, 2, 3}, result: Binary{Value: []byte{1, 2, 3}, Valid: true}},
		{source: _byteSlice(nil), result: Binary{}},
		{source: nilPointerByteArray, result: Binary{}},
		{source: nilPointerString, result: Binary{}},
	}

	for i, tt := range successfulTests {
		var r Binary
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
