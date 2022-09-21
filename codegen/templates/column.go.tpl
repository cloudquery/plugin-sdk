{
  Name:        "{{.Name}}",
  Type:        schema.{{.Type}},
  {{- if .Resolver}}
  Resolver:     {{.Resolver}},
  {{- end}}
  {{- if .Description}}
  Description:     `{{.Description}}`,
  {{- end}}
  {{- if .Options.PrimaryKey}}
  CreationOptions: schema.ColumnCreationOptions{
    {{- if .Options.PrimaryKey}}
      PrimaryKey: true,
    {{- end }}
  },
  {{- end}}
  {{- if .IgnoreInTests}}
  IgnoreInTests: true,
  {{- end}}
},
