package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudquery/plugin-sdk/schema"
)

const tableTmpl = `
# Table: {{.Name}}
{{ $.Description }}
## Columns
| Name          | Type          |
| ------------- | ------------- |
{{- range $column := $.Columns }}
|{{$column.Name}}|{{$column.Type | formatType}}|
{{- end }}
`

const allTablesTpml = `
# Source Plugin: {{.Name}}
## Tables
| Name          | Description   |
| ------------- | ------------- |
{{- range $table := $.Tables }}
|{{$table.Name}}|{{$table.Description }}|
{{- end }}
`

// GenerateSourcePluginDocs creates table documentation for the source plugin based on its list of tables
func (p *SourcePlugin) GenerateSourcePluginDocs(dir string) error {
	for _, table := range p.Tables() {
		if err := renderAllTables(table, dir); err != nil {
			fmt.Printf("render table %s error: %s", table.Name, err)
			return err
		}
	}
	t, err := template.New("").Parse(allTablesTpml)
	if err != nil {
		return fmt.Errorf("failed to parse template for README.md: %v", err)
	}

	outputPath := filepath.Join(dir, "README.md")
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", outputPath, err)
	}
	defer f.Close()
	if err := t.Execute(f, p); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	return nil
}

func renderAllTables(t *schema.Table, dir string) error {
	if err := renderTable(t, dir); err != nil {
		return err
	}
	for _, r := range t.Relations {
		if err := renderAllTables(r, dir); err != nil {
			return err
		}
	}
	return nil
}

func renderTable(table *schema.Table, dir string) error {
	t := template.New("").Funcs(map[string]interface{}{
		"formatType": formatType,
		"removeLineBreaks": func(text string) string {
			return strings.ReplaceAll(text, "\n", " ")
		},
	})
	t, err := t.New("").Parse(tableTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	outputPath := filepath.Join(dir, fmt.Sprintf("%s.md", table.Name))
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", outputPath, err)
	}
	defer f.Close()
	if err := t.Execute(f, table); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	return nil
}

func formatType(v schema.ValueType) string {
	return strings.TrimPrefix(v.String(), "Type")
}
