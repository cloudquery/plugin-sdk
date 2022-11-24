{{- if not .Parent }}
- [{{.Name}}]({{.Name}}.md)
{{- else}}
{{. | indentToDepth}}  - [{{.Name}}]({{.Name}}.md)
{{- end}}
{{- range $index, $rel := .Relations}}
{{- template "all_tables_entry.md.go.tpl" $rel}}
{{- end}}