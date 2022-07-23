package spec

import "github.com/xeipuuv/gojsonschema"

func ValidateSpec(schema string, spec interface{}) (*gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	specLoader := gojsonschema.NewGoLoader(spec)
	return gojsonschema.Validate(schemaLoader, specLoader)
}
