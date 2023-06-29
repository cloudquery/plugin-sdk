package configtype

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	cases := []struct {
		give string
		want time.Duration
	}{
		{"1ns", 1 * time.Nanosecond},
		{"20s", 20 * time.Second},
		{"-50m30s", -50*time.Minute - 30*time.Second},
	}
	for _, tc := range cases {
		var d Duration
		err := json.Unmarshal([]byte(`"`+tc.give+`"`), &d)
		if err != nil {
			t.Fatalf("error calling Unmarshal(%q): %v", tc.give, err)
		}
		if d.Duration() != tc.want {
			t.Errorf("Unmarshal(%q) = %v, want %v", tc.give, d.Duration(), tc.want)
		}
	}
}
