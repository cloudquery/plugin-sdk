package schema

import (
	"testing"
)

func TestMacaddrSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result Macaddr
	}{
		{
			source: mustParseMacaddr(t, "01:23:45:67:89:ab"),
			result: Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
		},
		{
			source: "01:23:45:67:89:ab",
			result: Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
		},
	}

	for i, tt := range successfulTests {
		var r Macaddr
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}
