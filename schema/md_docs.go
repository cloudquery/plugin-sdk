package schema

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// We are taking a similar approach to Cobra command generation for tables
// https://github.com/spf13/cobra/blob/main/doc/md_docs.go
// We also use this in our CLI cloudquery/cloudquery

var tableMdFuncs = map[string]interface{}{
	"removeLineBreaks": func(text string) string {
		return strings.ReplaceAll(text, "\n", " ")
	},
}

func GenerateMarkdownTree(tables []*Table, dir string) error {
	for _, table := range tables {
		if err := generateMarkdownTree(table, dir); err != nil {
			return err
		}
	}
	return nil
}

func generateMarkdownTree(table *Table, dir string) error {
	for _, child := range table.Relations {
		if err := generateMarkdownTree(child, dir); err != nil {
			return err
		}
	}
	t := template.New("").Funcs(tableMdFuncs)
	t, err := t.New("table_md").Parse(tableTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, table); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", table.Name)), buf.Bytes(), 0644)
}

const tableTmpl = `
# Table: {{.Name}}
{{ $.Description }}
## Columns
| Name          | Typ           | Description  |
| ------------- | ------------- | -----------  |
{{- range $column := $.Columns }}
|{{$column.Name}}|{{$column.Type}}|{{$column.Description|removeLineBreaks}}|
{{- end }}
`
