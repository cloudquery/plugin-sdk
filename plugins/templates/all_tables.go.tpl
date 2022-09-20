# Source Plugin: {{.Name}}
## Tables
| Name          | Description   |
| ------------- | ------------- |
{{- range $table := $.Tables }}
|{{$table.Name}}|{{$table.Description }}|
{{- end }}