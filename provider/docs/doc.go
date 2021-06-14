// Package docs helps create provider documentation
package docs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

// GenerateDocs creates table documentation for the provider based on it's ResourceMap
func GenerateDocs(p *provider.Provider, outputPath string) error {
	for _, table := range p.ResourceMap {
		if err := renderAllTables(table, outputPath); err != nil {
			fmt.Printf("render table error: %s", err)
			return err
		}
	}
	return nil
}

func renderAllTables(t *schema.Table, outputPath string) error {
	if err := renderTable(t, outputPath); err != nil {
		return err
	}
	for _, r := range t.Relations {
		if err := renderAllTables(r, outputPath); err != nil {
			return err
		}
	}
	return nil
}

func renderTable(table *schema.Table, path string) error {
	t := template.New("").Funcs(map[string]interface{}{
		"pgType": schema.GetPgTypeFromType,
		"removeLineBreaks": func(text string) string {
			return strings.ReplaceAll(text, "\n", " ")
		},
	})
	t, err := t.New("").Parse(tableTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err := t.Execute(&buf, table); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(path, "tables", fmt.Sprintf("%s.md", table.Name)), buf.Bytes(), 0644)
}

const tableTmpl = `
# Table: {{.Name}}
{{ $.Description }}
## Columns
| Name        | Type           | Description  |
| ------------- | ------------- | -----  |
{{- range $column := $.Columns }}
|{{$column.Name}}|{{$column.Type|pgType}}|{{$column.Description|removeLineBreaks}}|
{{- end }}
`
