package faker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFakerStruct struct {
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
	assert.Empty(t, a.E) // empty interfaces are not faked
}

type customType string

type fakerStructWithCustomType struct {
	A customType
	B map[string]customType
	C map[customType]string
	D []customType
}

func TestFakerWithCustomType(t *testing.T) {
	a := fakerStructWithCustomType{}
	if err := FakeObject(&a); err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, a.A)
	assert.NotEmpty(t, a.B)
	assert.NotEmpty(t, a.C)
	assert.NotEmpty(t, a.D)
}
