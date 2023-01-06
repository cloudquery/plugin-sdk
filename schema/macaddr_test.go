package schema

import (
	"testing"
)

func TestMacaddrSet(t *testing.T) {
	successfulTests := []struct {
		source any
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

func TestMacaddr_Size(t *testing.T) {
	tests := []struct {
		name string
		b    Macaddr
		want int
	}{
		{
			name: "present",
			b:    Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: Present},
			want: 6,
		},
		{
			name: "null",
			b:    Macaddr{Status: Null},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Size(); got != tt.want {
				t.Errorf("Macaddr.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
