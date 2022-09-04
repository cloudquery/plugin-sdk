package faker

import (
	"fmt"
	"testing"
)

type testFakerStruct struct {
	A int
	B string
	C *string
}

func TestFaker(t *testing.T) {
	a := testFakerStruct{}
	if err := FakeStruct(&a); err != nil {
		t.Fatal(err)
	}
	fmt.Println(a)
}
