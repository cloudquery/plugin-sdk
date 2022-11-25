# Table: {{.Name}}

{{ $.Description }}
{{ $length := len $.PrimaryKeys -}}
{{ if eq $length 1 }}
The primary key for this table is **{{ index $.PrimaryKeys 0 }}**.
{{ else }}
The composite primary key for this table is ({{ range $index, $pk := $.PrimaryKeys -}}
	{{if $index }}, {{end -}}
		**{{$pk}}**
	{{- end -}}).
{{ end }}

{{- if or ($.Relations) ($.Parent) }}
## Relations
{{- end }}
{{- if $.Parent }}
This table depends on [{{ $.Parent.Name }}]({{ $.Parent.Name }}.md).
{{- end}}
{{ if $.Relations }}
The following tables depend on {{.Name}}:
{{- range $rel := $.Relations }}
  - [{{ $rel.Name }}]({{ $rel.Name }}.md)
{{- end }}
{{- end }}

## Columns
| Name          | Type          |
| ------------- | ------------- |
{{- range $column := $.Columns }}
|{{$column.Name}}{{if $column.CreationOptions.PrimaryKey}} (PK){{end}}|{{$column.Type | formatType}}|
{{- end }}