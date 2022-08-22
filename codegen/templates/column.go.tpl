{
  Name:        "{{.Name}}",
  Type:        schema.{{.Type}},
  {{- if .Resolver}}
  Resolver:     {{.Resolver}},
  {{- end}}
  {{- if .Options.PrimaryKey}}
  CreationOptions: schema.ColumnCreationOptions{
    {{- if .Options.PrimaryKey}}
      PrimaryKey: true,
    {{- end }}
  },
  {{- end}}  
},
