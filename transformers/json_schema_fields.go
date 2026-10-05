package transformers

import (
	"reflect"
	"slices"
	"strings"
)

type jsonSchemaField struct {
	field  reflect.StructField
	depth  int
	tagged bool
}

func (t *structTransformer) jsonSchemaFields(structType reflect.Type) map[string]reflect.StructField {
	candidates := make(map[string][]jsonSchemaField)
	visited := make(map[reflect.Type]bool)
	level := []reflect.Type{structType}
	for depth := 0; len(level) > 0; depth++ {
		var next []reflect.Type
		for _, levelType := range level {
			for i := 0; i < levelType.NumField(); i++ {
				structField := levelType.Field(i)
				if isTypeIgnored(structField.Type) {
					continue
				}
				fieldType := structField.Type
				if fieldType.Name() == "" && fieldType.Kind() == reflect.Pointer {
					fieldType = fieldType.Elem()
				}
				tagName, _, _ := strings.Cut(structField.Tag.Get("json"), ",")
				isEmbeddedStruct := structField.Anonymous && fieldType.Kind() == reflect.Struct
				if isEmbeddedStruct && tagName == "" {
					next = append(next, fieldType)
					continue
				}
				if !structField.IsExported() && !isEmbeddedStruct {
					continue
				}
				name, err := t.jsonSchemaNameTransformer(structField)
				if err != nil || name == "" {
					continue
				}
				candidates[name] = append(candidates[name], jsonSchemaField{
					field:  structField,
					depth:  depth,
					tagged: tagName != "",
				})
			}
		}
		for _, levelType := range level {
			visited[levelType] = true
		}
		level = slices.DeleteFunc(next, func(nextType reflect.Type) bool { return visited[nextType] })
	}

	fields := make(map[string]reflect.StructField, len(candidates))
	for name, sameName := range candidates {
		if winner, ok := dominantJSONSchemaField(sameName); ok {
			fields[name] = winner.field
		}
	}
	return fields
}

func dominantJSONSchemaField(sameName []jsonSchemaField) (jsonSchemaField, bool) {
	shallowest := sameName[0].depth
	for _, f := range sameName {
		shallowest = min(shallowest, f.depth)
	}
	var atShallowest, tagged []jsonSchemaField
	for _, f := range sameName {
		if f.depth != shallowest {
			continue
		}
		atShallowest = append(atShallowest, f)
		if f.tagged {
			tagged = append(tagged, f)
		}
	}
	switch {
	case len(atShallowest) == 1:
		return atShallowest[0], true
	case len(tagged) == 1:
		return tagged[0], true
	default:
		return jsonSchemaField{}, false
	}
}
