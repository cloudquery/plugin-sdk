package codegen

import (
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

//go:embed templates/*.go.tpl
var TemplatesFS embed.FS

func (t *TableDefinition) GenerateTemplate(wr io.Writer) error {
	tpl, err := template.New("table.go.tpl").Funcs(template.FuncMap{
		"ToCamel": strcase.ToCamel,
		"ToLower": strings.ToLower,
	}).ParseFS(TemplatesFS, "templates/*")

	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tpl.Execute(wr, t); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
