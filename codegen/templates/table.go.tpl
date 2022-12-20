// Code generated by codegen; DO NOT EDIT.

package {{.Service}}

import (
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/cloudquery/plugins/source/{{.PluginName}}/client"
)

func {{.SubService | ToCamel}}() *schema.Table {
  return &schema.Table{
    Name:         "{{.Name}}",
    {{- if .Description}}
    Description:     `{{.Description}}`,
    {{- end}}
    {{- if .Resolver}}
    Resolver:     {{.Resolver}},
    {{- end}}
    {{- if .PreResourceResolver}}
    PreResourceResolver:     {{.PreResourceResolver}},
    {{- end}}
    {{- if .PostResourceResolver}}
    PostResourceResolver:     {{.PostResourceResolver}},
    {{- end}}
    {{- if .Multiplex}}
    Multiplex:    {{.Multiplex}},
    {{- end}}
    Columns: []schema.Column{
    {{range .Columns}}{{template "column.go.tpl" .}}{{end}}
    },
    {{with .Relations}}
    Relations: []*schema.Table{
    {{range .}}{{.}},
    {{end}}
		},
    {{end}}
  }
}