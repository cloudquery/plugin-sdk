package docs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type templateData struct {
	PluginName string
	Tables     schema.Tables
}

func (g *Generator) renderTablesAsMarkdown(dir string) error {
	for _, table := range g.tables {
		if err := g.renderAllTables(dir, table); err != nil {
			return err
		}
	}
	t, err := template.New("all_tables.md.go.tpl").Funcs(template.FuncMap{
		"indentToDepth": indentToDepth,
	}).ParseFS(templatesFS, "templates/all_tables*.md.go.tpl")
	if err != nil {
		return fmt.Errorf("failed to parse template for README.md: %v", err)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, templateData{PluginName: g.pluginName, Tables: g.tables}); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	content := formatMarkdown(b.String())
	outputPath := filepath.Join(dir, "README.md")
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", outputPath, err)
	}
	defer f.Close()
	f.WriteString(content)
	return nil
}

func (g *Generator) renderAllTables(dir string, t *schema.Table) error {
	if err := g.renderTable(dir, t); err != nil {
		return err
	}
	for _, r := range t.Relations {
		if err := g.renderAllTables(dir, r); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) renderTable(dir string, table *schema.Table) error {
	t := template.New("").Funcs(map[string]any{
		"title":           g.titleTransformer,
		"colTypeWithCode": colTypeWithCode,
	})
	t, err := t.New("table.md.go.tpl").ParseFS(templatesFS, "templates/table.md.go.tpl")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	outputPath := filepath.Join(dir, fmt.Sprintf("%s.md", table.Name))

	var b bytes.Buffer
	if err := t.Execute(&b, table); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	content := formatMarkdown(b.String())
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", outputPath, err)
	}
	defer f.Close()
	f.WriteString(content)
	return f.Close()
}

func formatMarkdown(s string) string {
	s = reMatchNewlines.ReplaceAllString(s, "\n\n")
	return reMatchHeaders.ReplaceAllString(s, `$1`+"\n\n")
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

func colTypeWithCode(dt arrow.DataType) string {
	typ := colType(dt, 0)
	if strings.Contains(typ, "<br>") {
		// pre block
		return "<pre>" + typ + "</pre>"
	}
	return "`" + typ + "`"
}

func colType(dt arrow.DataType, level int) string {
	nested, ok := dt.(arrow.NestedType)
	if !ok {
		return dt.String()
	}
	switch nested := nested.(type) {
	case *arrow.StructType:
		return structType(nested, level)
	case *arrow.MapType:
		return mapType(nested, level)
	case arrow.ListLikeType:
		return listLikeType(nested, level)
	default:
		return dt.String()
	}
}

func structType(dt *arrow.StructType, level int) string {
	if !needsMultiline(dt) {
		return simpleStruct(dt)
	}

	var buf strings.Builder
	buf.WriteString("<br>") // for starting elems with a newline
	pfx := strings.Repeat("&nbsp;", level+1)
	for i, field := range dt.Fields() {
		if i > 0 {
			buf.WriteString(",<br>")
		}
		buf.WriteString(pfx)
		buf.WriteString(field.Name)
		buf.WriteString(": ")
		buf.WriteString(colType(field.Type, level+1))
		if field.Nullable {
			buf.WriteString("?")
		}
	}
	return "struct<" + buf.String() + "<br>" + strings.Repeat("&nbsp;", level) + ">"
}

func simpleStruct(dt *arrow.StructType) string {
	field := dt.Field(0)
	res := "struct<" + field.Name + ": " + field.Type.String()
	if field.Nullable {
		res += "?"
	}
	return res + ">"
}

func mapType(dt *arrow.MapType, level int) string {
	keys := colType(dt.KeyType(), level)
	items := colType(dt.ItemType(), level) + "?" // always have nullable values in map
	return "map<" + keys + ", " + items + ">"
}

func listLikeType(dt arrow.ListLikeType, level int) string {
	elemField := dt.ElemField()
	nested := needsMultiline(dt)
	if nested {
		level++ // nested types will require additional handling
	}

	var elems string
	var pfx, sfx string
	if nested {
		pfx = "<br>" + strings.Repeat("&nbsp;", level)
		sfx = "<br>"
		if level > 0 {
			sfx += strings.Repeat("&nbsp;", level-1)
		}
	}
	elems = "<" + pfx + colType(elemField.Type, level)
	if dt.ElemField().Nullable {
		elems += "?"
	}
	elems += sfx + ">"
	switch dt := dt.(type) {
	case *arrow.ListType:
		return "list" + elems
	case *arrow.LargeListType:
		return "large_list" + elems
	case *arrow.FixedSizeListType:
		return "fixed_size_list" + elems + "[" + strconv.FormatInt(int64(dt.Len()), 10) + "]"
	case *arrow.ListViewType:
		return "list_view" + elems
	default:
		return dt.Name() + elems
	}
}

func needsMultiline(dt arrow.DataType) bool {
	nested, ok := dt.(arrow.NestedType)
	if !ok {
		return false
	}
	switch nested := nested.(type) {
	case *arrow.MapType:
		return needsMultiline(nested.ItemType()) // keys are presumed to be simple
	case arrow.ListLikeType:
		return needsMultiline(nested.Elem())
	case *arrow.StructType:
		switch nested.NumFields() {
		case 0:
			return false
		case 1:
			return needsMultiline(nested.Field(0).Type)
		default:
			return true
		}
	default:
		return false
	}
}
