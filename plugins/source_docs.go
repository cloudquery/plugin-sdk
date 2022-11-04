package plugins

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudquery/plugin-sdk/schema"
)

//go:embed templates/*.go.tpl
var templatesFS embed.FS

type SourceDocsFormat string

const (
	SourceDocsFormatMarkdown = "markdown"
	SourceDocsFormatJSON     = "json"
)

var SourceDocsFormats = []SourceDocsFormat{
	SourceDocsFormatMarkdown,
	SourceDocsFormatJSON,
}

func (s SourceDocsFormat) String() string {
	return string(s)
}

func (s SourceDocsFormat) Validate() error {
	for _, f := range SourceDocsFormats {
		if s == f {
			return nil
		}
	}
	return fmt.Errorf("invalid format: %v", s.String())
}

// GenerateSourcePluginDocs creates table documentation for the source plugin based on its list of tables
func (p *SourcePlugin) GenerateSourcePluginDocs(dir string, format SourceDocsFormat) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	SetDestinationManagedCqColumns(p.Tables())

	switch format {
	case SourceDocsFormatMarkdown:
		return p.renderTablesAsMarkdown(dir)
	case SourceDocsFormatJSON:
		return p.renderTablesAsJSON(dir)
	default:
		return fmt.Errorf("unsupported format: %v", format.String())
	}
}

type jsonTable struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Columns     []jsonColumn `json:"columns"`
	Relations   []jsonTable  `json:"relations"`
}

type jsonColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (p *SourcePlugin) renderTablesAsJSON(dir string) error {
	tables := p.jsonifyTables(p.Tables())
	b, err := json.MarshalIndent(tables, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tables as json: %v", err)
	}
	outputPath := filepath.Join(dir, "__tables.json")
	return os.WriteFile(outputPath, b, 0644)
}

func (p *SourcePlugin) jsonifyTables(tables schema.Tables) []jsonTable {
	jsonTables := make([]jsonTable, len(tables))
	for i, table := range tables {
		jsonColumns := make([]jsonColumn, len(table.Columns))
		for c, col := range table.Columns {
			jsonColumns[c] = jsonColumn{
				Name: col.Name,
				Type: col.Type.String(),
			}
		}
		jsonTables[i] = jsonTable{
			Name:        table.Name,
			Description: table.Description,
			Columns:     jsonColumns,
			Relations:   p.jsonifyTables(table.Relations),
		}
	}
	return jsonTables
}

func (p *SourcePlugin) renderTablesAsMarkdown(dir string) error {
	for _, table := range p.Tables() {
		if err := renderAllTables(table, dir); err != nil {
			return err
		}
	}
	t, err := template.New("all_tables.md.go.tpl").ParseFS(templatesFS, "templates/all_tables.md.go.tpl")
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
	})
	t, err := t.New("table.md.go.tpl").ParseFS(templatesFS, "templates/table.md.go.tpl")
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
