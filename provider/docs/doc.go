// Package docs helps create provider documentation
package docs

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudquery/cq-provider-sdk/provider"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

const (
	tablesDir = "tables"
)

// GenerateDocs creates table documentation for the provider based on it's ResourceMap
func GenerateDocs(p *provider.Provider, outputPath string, deleteOld bool) error {
	if deleteOld {
		if err := deleteOldFiles(outputPath); err != nil {
			fmt.Printf("failed to remove old docs: %s", err)
			return err
		}
	}
	for _, table := range p.ResourceMap {
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
			if err := os.Mkdir(tablesPath, 0644); err != nil {
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
	return ioutil.WriteFile(filepath.Join(path, tablesDir, fmt.Sprintf("%s.md", table.Name)), buf.Bytes(), 0644)
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
