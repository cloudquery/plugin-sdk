package configtype_test

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/configtype"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	now, _ := time.Parse(time.RFC3339Nano, time.RFC3339Nano)

	cases := []struct {
		give string
		want time.Time
	}{
		{"1ns", now.Add(1 * time.Nanosecond)},
		{"20s", now.Add(20 * time.Second)},
		{"-50m30s", now.Add(-50*time.Minute - 30*time.Second)},
		{"2021-09-01T00:00:00Z", time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)},
		{"2021-09-01T00:00:00.123Z", time.Date(2021, 9, 1, 0, 0, 0, 123000000, time.UTC)},
		{"2021-09-01T00:00:00.123456Z", time.Date(2021, 9, 1, 0, 0, 0, 123456000, time.UTC)},
		{"2021-09-01T00:00:00.123456789Z", time.Date(2021, 9, 1, 0, 0, 0, 123456789, time.UTC)},
		{"2021-09-01T00:00:00.123+02:00", time.Date(2021, 9, 1, 0, 0, 0, 123000000, time.FixedZone("CET", 2*60*60))},
		{"2021-09-01T00:00:00.123456+02:00", time.Date(2021, 9, 1, 0, 0, 0, 123456000, time.FixedZone("CET", 2*60*60))},
		{"2021-09-01T00:00:00.123456789+02:00", time.Date(2021, 9, 1, 0, 0, 0, 123456789, time.FixedZone("CET", 2*60*60))},
		{"2024-09-26T10:18:07.37338-04:00", time.Date(2024, 9, 26, 10, 18, 7, 373380000, time.FixedZone("EDT", -4*60*60))},
		{"2021-09-01", time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)},
		{"now", now},
		{"2 days from now", now.AddDate(0, 0, 2)},
		{"5 months ago", now.AddDate(0, -5, 0)},
	}
	for _, tc := range cases {
		var d configtype.Time
		err := json.Unmarshal([]byte(`"`+tc.give+`"`), &d)
		if err != nil {
			t.Fatalf("error calling Unmarshal(%q): %v", tc.give, err)
		}
		computedTime := d.AsTime(now)
		if !computedTime.Equal(tc.want) {
			t.Errorf("Unmarshal(%q) = %v, want %v", tc.give, computedTime, tc.want)
		}
	}
}

func TestTime_JSONSchema(t *testing.T) {
	sc := (&jsonschema.Reflector{RequiredFromJSONSchemaTags: true}).Reflect(configtype.Time{})
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
		{
			Name: "not relative duration",
			Err:  true,
			Spec: `"10 days"`,
		},
		{
			Name: "relative duration",
			Err:  false,
			Spec: `"10 months from now"`,
		},
		{
			Name: "complex relative duration",
			Err:  false,
			Spec: `"10 months 3 days 4h20m from now"`,
		},
	},
		func() []testCase {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			const (
				cases      = 20
				maxDur     = int64(100 * time.Hour)
				maxDurHalf = maxDur / 2
			)
			now := time.Now()
			var result []testCase
			for i := 0; i < cases; i++ {
				val := rnd.Int63n(maxDur) - maxDurHalf
				dur := must(configtype.ParseTime(time.Duration(val).String()))

				durationData, err := dur.MarshalJSON()
				require.NoError(t, err)
				result = append(result, testCase{
					Name: string(durationData),
					Spec: string(durationData),
				})

				tim := must(configtype.ParseTime(must(marshalString(now.Add(time.Duration(val))))))

				timeData, err := tim.MarshalJSON()
				require.NoError(t, err)
				result = append(result, testCase{
					Name: string(timeData),
					Spec: string(timeData),
				})
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
