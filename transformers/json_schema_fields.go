package transformers

import (
	"reflect"
	"strings"
)

type jsonSchemaField struct {
	field       reflect.StructField
	depth       int
	tagged      bool
	addressable bool
}

type embeddedJSONStruct struct {
	structType  reflect.Type
	count       int
	addressable bool
}

type jsonSchemaFieldsKey struct {
	structType  reflect.Type
	addressable bool
}

func (t *structTransformer) jsonSchemaFields(structType reflect.Type, addressable bool) map[string]jsonSchemaField {
	key := jsonSchemaFieldsKey{structType: structType, addressable: addressable}
	if fields, ok := t.jsonSchemaFieldsCache[key]; ok {
		return fields
	}
	if t.jsonSchemaFieldsCache == nil {
		t.jsonSchemaFieldsCache = make(map[jsonSchemaFieldsKey]map[string]jsonSchemaField)
	}
	fields := t.computeJSONSchemaFields(structType, addressable)
	t.jsonSchemaFieldsCache[key] = fields
	return fields
}

func (t *structTransformer) computeJSONSchemaFields(structType reflect.Type, addressable bool) map[string]jsonSchemaField {
	candidates := make(map[string][]jsonSchemaField)
	visited := make(map[reflect.Type]bool)
	level := []*embeddedJSONStruct{{structType: structType, count: 1, addressable: addressable}}
	for depth := 0; len(level) > 0; depth++ {
		for _, embedded := range level {
			visited[embedded.structType] = true
		}
		var next []*embeddedJSONStruct
		nextByType := make(map[reflect.Type]*embeddedJSONStruct)
		for _, embedded := range level {
			for i := 0; i < embedded.structType.NumField(); i++ {
				structField := embedded.structType.Field(i)
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
					if visited[fieldType] {
						continue
					}
					reached, ok := nextByType[fieldType]
					if !ok {
						reached = &embeddedJSONStruct{
							structType:  fieldType,
							addressable: embedded.addressable || structField.Type.Kind() == reflect.Pointer,
						}
						nextByType[fieldType] = reached
						next = append(next, reached)
					}
					reached.count++
					continue
				}
				if !structField.IsExported() && !isEmbeddedStruct {
					continue
				}
				name, err := t.jsonSchemaNameTransformer(structField)
				if err != nil || name == "" {
					continue
				}
				candidate := jsonSchemaField{
					field:       structField,
					depth:       depth,
					tagged:      tagName != "",
					addressable: embedded.addressable,
				}
				candidates[name] = append(candidates[name], candidate)
				if embedded.count > 1 {
					candidates[name] = append(candidates[name], candidate)
				}
			}
		}
		level = next
	}

	fields := make(map[string]jsonSchemaField, len(candidates))
	for name, sameName := range candidates {
		if winner, ok := dominantJSONSchemaField(sameName); ok {
			fields[name] = winner
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
