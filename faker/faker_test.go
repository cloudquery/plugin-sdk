package faker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFakerStruct struct {
	A int
	B string
	C *string
}

type testFakerStructWithEFace struct {
	A int
	B string
	C *string
	E interface{}
}

func TestFaker(t *testing.T) {
	a := testFakerStruct{}
	if err := FakeObject(&a); err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, a.A)
	assert.NotEmpty(t, a.B)
	assert.NotEmpty(t, a.C)
}

func TestFakerWithSkipFields(t *testing.T) {
	a := testFakerStruct{}
	if err := FakeObject(&a, WithSkipFields("B", "C")); err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, a.A)
	assert.Empty(t, a.B)
	assert.Empty(t, a.C)
}

func TestFakerWithSkipEFace(t *testing.T) {
	a := testFakerStructWithEFace{}
	if err := FakeObject(&a, WithSkipFields("B"), WithSkipEFace()); err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, a.A)
	assert.Empty(t, a.B)
	assert.NotEmpty(t, a.C)
	assert.Empty(t, a.E)
}
