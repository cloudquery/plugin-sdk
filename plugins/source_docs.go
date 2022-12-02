package plugins

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloudquery/plugin-sdk/schema"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

//go:embed templates/*.go.tpl
var templatesFS embed.FS

// GenerateSourcePluginDocs creates table documentation for the source plugin based on its list of tables
func (p *SourcePlugin) GenerateSourcePluginDocs(dir, format string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	SetDestinationManagedCqColumns(p.Tables())

	switch format {
	case "markdown":
		return p.renderTablesAsMarkdown(dir)
	case "json":
		return p.renderTablesAsJSON(dir)
	default:
		return fmt.Errorf("unsupported format: %v", format)
	}
}

type jsonTable struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Columns     []jsonColumn `json:"columns"`
	Relations   []jsonTable  `json:"relations"`
}

type jsonColumn struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	IsPrimaryKey bool   `json:"is_primary_key,omitempty"`
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
				Name:         col.Name,
				Type:         col.Type.String(),
				IsPrimaryKey: col.CreationOptions.PrimaryKey,
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
	// render all tables
	grp, _ := errgroup.WithContext(context.Background())
	for _, table := range p.Tables() {
		table := table
		grp.Go(func() error {
			return renderAllTables(table, dir)
		})
	}
	if err := grp.Wait(); err != nil {
		return err
	}

	t, err := template.New("all_tables.md.go.tpl").Funcs(template.FuncMap{
		"indentToDepth": indentToDepth,
	}).ParseFS(templatesFS, "templates/all_tables*.md.go.tpl")
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

	// render all relations
	grp, _ := errgroup.WithContext(context.Background())
	for _, table := range t.Relations {
		table := table
		grp.Go(func() error {
			return renderAllTables(table, dir)
		})
	}
	return grp.Wait()
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

	sortTableColumns(table)
	if err := t.Execute(f, table); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	return nil
}

func sortTableColumns(table *schema.Table) {
	const cqPfx = "_cq"
	slices.SortStableFunc(table.Columns, func(a, b schema.Column) bool {
		switch {
		case strings.HasPrefix(a.Name, cqPfx):
			return !strings.HasPrefix(b.Name, cqPfx) || a.Name < b.Name
		case strings.HasPrefix(b.Name, cqPfx):
			return false
		case a.CreationOptions.PrimaryKey:
			return !b.CreationOptions.PrimaryKey || a.Name < b.Name
		case b.CreationOptions.PrimaryKey:
			return false
		default:
			return a.Name < b.Name
		}
	})
}

func formatType(v schema.ValueType) string {
	return strings.TrimPrefix(v.String(), "Type")
}

func indentToDepth(table *schema.Table) string {
	s := ""
	t := table
	for t.Parent != nil {
		s += "  "
		t = t.Parent
	}
	return s
}
