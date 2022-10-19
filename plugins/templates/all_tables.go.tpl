# Source Plugin: {{.Name}}
## Tables
| Name          |
| ------------- |
{{- range $table := $.Tables }}
| [{{$table.Name}}]({{$table.Name}}.md) |
{{- range $index, $rel := $table.Relations}}
| ↳ [{{$rel.Name}}]({{$rel.Name}}.md) |
{{- end}}
{{- end }}