package mockgen

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"

	"github.com/cloudquery/plugin-sdk/caser"
)

//go:embed templates/*.go.tpl
var templatesFS embed.FS

type Options struct {
	// ShouldInclude tests whether a method should be included in the generated interfaces. If it returns true,
	// the method will be included. MethodHasPrefix and MethodHasSuffix can be used inside a custom function here
	// to customize the behavior.
	ShouldInclude func(reflect.Method) bool
}

func (o *Options) SetDefaults() {
	if o.ShouldInclude == nil {
		o.ShouldInclude = func(reflect.Method) bool { return true }
	}
}

type Option func(*Options)

func WithIncludeFunc(f func(reflect.Method) bool) Option {
	return func(o *Options) {
		o.ShouldInclude = f
	}
}

// GenerateInterfaces generates service interfaces to be used for generating
// mocks. The clients passed in as the first argument should be structs that will be used to
// generate the service interfaces. The second argument, dir, is the path to the output
// directory where the service interface files will be created.
func GenerateInterfaces(clients []any, dir string, opts ...Option) error {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	options.SetDefaults()

	services := make([]serviceInfo, 0)
	for _, client := range clients {
		services = append(services, getServiceInfo(client, options))
	}

	// write individual service files
	serviceTpl, err := template.New("service.go.tpl").ParseFS(templatesFS, "templates/service.go.tpl")
	if err != nil {
		return err
	}

	for _, service := range services {
		buff := bytes.Buffer{}
		if err := serviceTpl.Execute(&buff, service); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
		filePath := path.Join(dir, fmt.Sprintf("%s.go", service.PackageName))
		err := formatAndWriteFile(filePath, buff)
		if err != nil {
			return fmt.Errorf("failed to format and write file for service %v: %w", service.Name, err)
		}
	}

	return nil
}

// Adapted from https://stackoverflow.com/a/54129236
func signature(name string, f any) string {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return "<not a function>"
	}

	buf := strings.Builder{}
	buf.WriteString(name + "(")
	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		if t.IsVariadic() && i == t.NumIn()-1 {
			buf.WriteString("..." + strings.TrimPrefix(t.In(i).String(), "[]"))
		} else {
			buf.WriteString(t.In(i).String())
		}
	}
	buf.WriteString(")")
	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			buf.WriteString(" (")
		} else {
			buf.WriteString(" ")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(t.Out(i).String())
		}
		if numOut > 1 {
			buf.WriteString(")")
		}
	}

	return buf.String()
}

type serviceInfo struct {
	Import      string
	Name        string
	PackageName string
	ClientName  string
	Signatures  []string
}

func getServiceInfo(client any, opts *Options) serviceInfo {
	v := reflect.ValueOf(client)
	t := v.Type()
	pkgPath := t.Elem().PkgPath()
	parts := strings.Split(pkgPath, "/")
	pkgName := parts[len(parts)-1]
	csr := caser.New()
	name := csr.ToPascal(pkgName)
	clientName := name + "Client"
	signatures := make([]string, 0)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if opts.ShouldInclude(method) {
			sig := signature(method.Name, v.Method(i).Interface())
			signatures = append(signatures, sig)
		}
	}
	return serviceInfo{
		Import:      pkgPath,
		Name:        name,
		PackageName: pkgName,
		ClientName:  clientName,
		Signatures:  signatures,
	}
}

func formatAndWriteFile(filePath string, buff bytes.Buffer) error {
	content := buff.Bytes()
	formattedContent, err := format.Source(buff.Bytes())
	if err != nil {
		fmt.Printf("failed to format source: %s: %v\n", filePath, err)
	} else {
		content = formattedContent
	}
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	return nil
}
