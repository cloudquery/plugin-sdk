package configtype_test

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/configtype"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/google/go-cmp/cmp"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	cases := []struct {
		give string
		want time.Duration
	}{
		{"1ns", 1 * time.Nanosecond},
		{"20s", 20 * time.Second},
		{"-50m30s", -50*time.Minute - 30*time.Second},
		{"25 minute", 25 * time.Minute},
		{"50 minutes", 50 * time.Minute},
		{"10 years ago", -10 * 365 * 24 * time.Hour},
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

func TestDuration_JSONSchema(t *testing.T) {
	sc := (&jsonschema.Reflector{RequiredFromJSONSchemaTags: true}).Reflect(configtype.Duration{})
	schema, err := json.MarshalIndent(sc, "", "  ")
	require.NoError(t, err)

	validator, err := plugin.JSONSchemaValidator(string(schema))
	require.NoError(t, err)

	type testCase struct {
		Name string
		Spec string
		Err  bool
	}

	for _, tc := range append([]testCase{
		{
			Name: "empty",
			Err:  true,
			Spec: `""`,
		},
		{
			Name: "null",
			Err:  true,
			Spec: `null`,
		},
		{
			Name: "bad type",
			Err:  true,
			Spec: `false`,
		},
		{
			Name: "bad format",
			Err:  true,
			Spec: `false`,
		},
	},
		func() []testCase {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			const (
				cases      = 20
				maxDur     = int64(100 * time.Hour)
				maxDurHalf = maxDur / 2
			)
			result := make([]testCase, cases)
			for i := 0; i < cases; i++ {
				val := rnd.Int63n(maxDur) - maxDurHalf
				d := configtype.NewDuration(time.Duration(val))

				data, err := d.MarshalJSON()
				require.NoError(t, err)
				result[i] = testCase{
					Name: string(data),
					Spec: string(data),
				}
			}
			return result
		}()...,
	) {
		t.Run(tc.Name, func(t *testing.T) {
			var val any
			err := json.Unmarshal([]byte(tc.Spec), &val)
			require.NoError(t, err)
			if tc.Err {
				require.Error(t, validator.Validate(val))
			} else {
				require.NoError(t, validator.Validate(val))
			}
		})
	}
}
