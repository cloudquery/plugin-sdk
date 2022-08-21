{
  Name:        "{{.Name}}",
  Type:        schema.{{.Type}},
  {{- if .Resolver}}
  Resolver:     {{.Resolver}},
  {{- end}}
},
