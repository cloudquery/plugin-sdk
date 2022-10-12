package faker

import (
	"math/rand"
	"testing"
	"time"

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

func TestHelpers(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	vals := []interface{}{
		Name(),
		Word(),
		RandomUnixTime(),
		Timestamp(),
		UUIDHyphenated(),
		UUIDDigit(),
	}
	for i, v := range vals {
		assert.NotEmptyf(t, v, "Helper %d failed", i)
	}
}
