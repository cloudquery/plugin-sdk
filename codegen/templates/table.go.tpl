{
		Name:         "{{.Name}}",
		{{- if .Description}}
        Description:     "{{.Description}}",
    {{- end}}
		{{- if .Resolver}}
    Resolver:     {{.Resolver}},
    {{- end}}
		{{- if .Multiplex}}
    Multiplex:    {{.Multiplex}},
    {{- end}}
		{{- if .IgnoreError}}
    IgnoreError:  {{.IgnoreError}},
    {{- end}}
		Columns: []schema.Column{
{{range .Columns}}{{template "column.go.tpl" .}}{{end}}
		},
{{with .Relations}}
		Relations: []*schema.Table{
{{range .}}{{template "table.go.tpl" .}}{{end}}
		},
{{end}}
}