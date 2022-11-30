package faker

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testFakerStruct struct {
	A int
	B string
	C *string
	D time.Time
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
	assert.NotEmpty(t, a.D)
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

type complexType struct {
	IPAddress net.IP
}

func TestFakerCanFakeNetIP(t *testing.T) {
	a := complexType{}
	if err := FakeObject(&a); err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, a.IPAddress)
	assert.Equal(t, "1.1.1.1", a.IPAddress.String())
}
