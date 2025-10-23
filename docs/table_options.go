package docs

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"

	schemaDocs "github.com/cloudquery/codegen/jsonschema/docs"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	invoschema "github.com/invopop/jsonschema"
)

func transformDescription(table *schema.Table, tableNamesToOptionsDocs map[string]string) {
	if tableNamesToOptionsDocs[table.Name] != "" {
		table.Description = table.Description + "\n\n" + tableNamesToOptionsDocs[table.Name]
	}
	for _, rel := range table.Relations {
		transformDescription(rel, tableNamesToOptionsDocs)
	}
}

func TableOptionsDescriptionTransformer(tableOptions any, jsonSchema string) (schema.Transform, error) {
	var sc invoschema.Schema
	if err := json.Unmarshal([]byte(jsonSchema), &sc); err != nil {
		return nil, err
	}
	tableNamesToOptionsDocs := make(map[string]string)
	tableOptionsType := reflect.ValueOf(tableOptions).Elem().Type()
	for i := range tableOptionsType.NumField() {
		field := tableOptionsType.Field(i)
		fieldType := field.Type.String()
		if strings.Contains(fieldType, ".") {
			fieldType = strings.Split(fieldType, ".")[1]
		}
		defValue, ok := sc.Definitions[fieldType]
		if !ok {
			return nil, errors.New("definition not found for " + field.Name)
		}
		tableName := strings.Split(field.Tag.Get("json"), ",")[0]
		if tableName == "" {
			return nil, errors.New("json tag not found for table " + field.Name)
		}
		newRoot := sc
		newRoot.ID = "Table Options"
		newRoot.Ref = "#/$defs/" + "Table Options"
		newRoot.Definitions["Table Options"] = defValue
		sch, _ := json.Marshal(newRoot)
		doc, _ := schemaDocs.Generate(sch, 1)
		tocRegex := regexp.MustCompile(`# Table of contents[\s\S]+?##`)
		tableNamesToOptionsDocs[tableName] = tocRegex.ReplaceAllString(doc, "##")
	}

	return func(table *schema.Table) error {
		transformDescription(table, tableNamesToOptionsDocs)
		return nil
	}, nil
}
