package scheduler_test

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

//go:embed strategy.json
var jsonSchema string

func TestStrategy_JSONSchema(t *testing.T) {
	sc := (&jsonschema.Reflector{RequiredFromJSONSchemaTags: true}).ReflectFromType(reflect.TypeOf(scheduler.StrategyDFS))
	data, err := json.MarshalIndent(sc, "", "  ")
	require.NoError(t, err)
	require.JSONEq(t, string(data)+"\n", jsonSchema)
}
