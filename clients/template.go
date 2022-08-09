package clients

import (
	"strings"
	"text/template"
)

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"indent": indent,
	}
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}
