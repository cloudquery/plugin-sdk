package configtype_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/configtype"
	"github.com/google/go-cmp/cmp"
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
		var d configtype.Duration
		err := json.Unmarshal([]byte(`"`+tc.give+`"`), &d)
		if err != nil {
			t.Fatalf("error calling Unmarshal(%q): %v", tc.give, err)
		}
		if d.Duration() != tc.want {
			t.Errorf("Unmarshal(%q) = %v, want %v", tc.give, d.Duration(), tc.want)
		}
	}
}

func TestComparability(t *testing.T) {
	cases := []struct {
		give    configtype.Duration
		compare configtype.Duration
		equal   bool
	}{
		{configtype.NewDuration(0), configtype.NewDuration(0), true},
		{configtype.NewDuration(0), configtype.NewDuration(1), false},
	}
	for _, tc := range cases {
		if (tc.give == tc.compare) != tc.equal {
			t.Errorf("comparing %v and %v should be %v", tc.give, tc.compare, tc.equal)
		}

		diff := cmp.Diff(tc.give, tc.compare)
		if tc.equal && diff != "" {
			t.Errorf("comparing %v and %v should be equal, but diff is %s", tc.give, tc.compare, diff)
		} else if !tc.equal && diff == "" {
			t.Errorf("comparing %v and %v should not be equal, but diff is empty", tc.give, tc.compare)
		}
	}
}
