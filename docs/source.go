// Package docs helps create plugin documentation
package docs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cloudquery/plugin-sdk/plugins"
	"github.com/cloudquery/plugin-sdk/schema"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	tablesDir = "tables"
)

// GenerateSourcePluginDocs creates table documentation for the source plugin based on its list of tables
func GenerateSourcePluginDocs(p *plugins.SourcePlugin, outputPath string, deleteOld bool) error {
	if deleteOld {
		if err := deleteOldFiles(outputPath); err != nil {
			fmt.Printf("failed to remove old docs: %s", err)
			return err
		}
	}
	for _, table := range p.Tables() {
		if err := renderAllTables(table, outputPath); err != nil {
			fmt.Printf("render table %s error: %s", table.Name, err)
			return err
		}
	}
	return nil
}

// deleteOldFiles removes old files from tables directory, creates tables directory if it does not exist
func deleteOldFiles(outputPath string) error {
	tablesPath := path.Join(outputPath, tablesDir)
	dir, err := ioutil.ReadDir(tablesPath)
	if err != nil {
		// create directory if it does not exist
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(tablesPath, 0744); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("Failed to generate docs: %s\n", err)
	}
	for _, d := range dir {
		if err := os.RemoveAll(path.Join([]string{tablesPath, d.Name()}...)); err != nil {
			return fmt.Errorf("Failed to remove old docs: %s\n", err)
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
		"formatType": formatType,
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
	return ioutil.WriteFile(filepath.Join(path, tablesDir, fmt.Sprintf("%s.md", table.Name)), buf.Bytes(), 0644)
}

func formatType(v schema.ValueType) string {
	return strings.TrimPrefix(v.String(), "Type")
}

const tableTmpl = `
# Table: {{.Name}}
{{ $.Description }}
## Columns
| Name        | Type           | Description  |
| ------------- | ------------- | -----  |
{{- range $column := $.Columns }}
|{{$column.Name}}|{{$column.Type | formatType}}|{{$column.Description|removeLineBreaks}}|
{{- end }}
`
