# Source Plugin: {{.Name}}
## Tables
| Name          | Relations | Description   |
| ------------- | --------- | ------------- |
{{- range $table := $.Tables }}
| [{{$table.Name}}]({{$table.Name}}.md)| {{range $index, $rel := $table.Relations}}{{if $index}}<br />{{end}}[{{$rel.Name}}]({{$rel.Name}}.md){{end}}| {{$table.Description }}|
{{- end }}