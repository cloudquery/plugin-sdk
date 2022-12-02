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
	IPAddress      net.IP
	IPAddresses    []net.IP
	PtrIPAddress   *net.IP
	PtrIPAddresses []*net.IP
	NestedComplex  struct {
		IPAddress net.IP
	}
}

func TestFakerCanFakeNetIP(t *testing.T) {
	a := complexType{}
	if err := FakeObject(&a); err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, a.IPAddress)
	assert.Equal(t, "1.1.1.1", a.IPAddress.String())

	assert.Equal(t, 1, len(a.IPAddresses))
	assert.Equal(t, "1.1.1.1", a.IPAddresses[0].String())

	assert.NotEmpty(t, a.PtrIPAddress)
	assert.Equal(t, "1.1.1.1", a.PtrIPAddress.String())

	assert.Equal(t, 1, len(a.PtrIPAddresses))
	assert.Equal(t, "1.1.1.1", a.PtrIPAddresses[0].String())

	assert.NotEmpty(t, a.NestedComplex.IPAddress)
	assert.Equal(t, "1.1.1.1", a.NestedComplex.IPAddress.String())
}
