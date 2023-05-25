package scalar

import (
	"net"
	"testing"
)

func TestMacaddrSet(t *testing.T) {
	successfulTests := []struct {
		source any
		result Mac
	}{
		{
			source: mustParseMacaddr(t, "01:23:45:67:89:ab"),
			result: Mac{Value: mustParseMacaddr(t, "01:23:45:67:89:ab"), Valid: true},
		},
		{
			source: "01:23:45:67:89:ab",
			result: Mac{Value: mustParseMacaddr(t, "01:23:45:67:89:ab"), Valid: true},
		},
		{
			source: &Mac{Value: mustParseMacaddr(t, "01:23:45:67:89:ab"), Valid: true},
			result: Mac{Value: mustParseMacaddr(t, "01:23:45:67:89:ab"), Valid: true},
		},
	}

	for i, tt := range successfulTests {
		var r Mac
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !r.Equal(&tt.result) {
			t.Errorf("%d: %v != %v", i, r, tt.result)
		}
	}
}

// nolint:unparam
func mustParseMacaddr(t testing.TB, s string) net.HardwareAddr {
	addr, err := net.ParseMAC(s)
	if err != nil {
		t.Fatal(err)
	}

	return addr
}
