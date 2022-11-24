
{{. | indentToDepth}}- [{{.Name}}]({{.Name}}.md)
{{- range $index, $rel := .Relations}}
{{- template "all_tables_entry.md.go.tpl" $rel}}
{{- end}}