package helpers

import (
	"testing"
)

type testStruct struct {
	test string
}

func TestToPointer(t *testing.T) {
	// passing string should return pointer to string
	give := "test"
	got := ToPointer(give)
	if _, ok := got.(*string); !ok {
		t.Errorf("ToPointer(%q) returned %q, expected type *string", give, got)
	}

	// passing struct by value should return pointer to (copy of the) struct
	giveObj := testStruct{
		test: "value",
	}
	gotStruct := ToPointer(giveObj)
	if _, ok := gotStruct.(*testStruct); !ok {
		t.Errorf("ToPointer(%q) returned %q, expected type *testStruct", giveObj, gotStruct)
	}

	// passing pointer should return the original pointer
	ptr := &giveObj
	gotPointer := ToPointer(ptr)
	if gotPointer != ptr {
		t.Errorf("ToPointer(%q) returned %q, expected %q", ptr, gotPointer, ptr)
	}

	// passing nil should return nil back without panicking
	gotNil := ToPointer(nil)
	if gotNil != nil {
		t.Errorf("ToPointer(%v) returned %q, expected nil", nil, gotNil)
	}

	// passing number should return pointer to number
	giveNumber := int64(0)
	gotNumber := ToPointer(giveNumber)
	if v, ok := gotNumber.(*int64); !ok {
		t.Errorf("ToPointer(%q) returned %q, expected type *int64", giveNumber, gotNumber)
		if *v != 0 {
			t.Errorf("ToPointer(%q) returned %q, expected 0", giveNumber, gotNumber)
		}
	}
}
