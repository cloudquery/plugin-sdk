package scheduler_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestStrategy_JSONSchema(t *testing.T) {
	sc := (&jsonschema.Reflector{RequiredFromJSONSchemaTags: true}).Reflect(scheduler.StrategyDFS)
	schema, err := json.MarshalIndent(sc, "", "  ")
	require.NoError(t, err)

	validator, err := plugin.JSONSchemaValidator(string(schema))
	require.NoError(t, err)

	type testCase struct {
		Name string
		Spec string
		Err  bool
	}

	for _, tc := range []testCase{
		{
			Name: "dfs scheduler",
			Spec: `"dfs"`,
		},
		{
			Name: "round-robin scheduler",
			Spec: `"round-robin"`,
		},
		{
			Name: "shuffle scheduler",
			Spec: `"shuffle"`,
		},
		{
			Name: "empty scheduler",
			Err:  true,
			Spec: `""`,
		},
		{
			Name: "bad scheduler",
			Err:  true,
			Spec: `"bad"`,
		},
		{
			Name: "bad scheduler type",
			Err:  true,
			Spec: `123`,
		},
		{
			Name: "null scheduler type",
			Err:  true,
			Spec: `null`,
		},
	} {
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
