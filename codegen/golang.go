package codegen

import (
	"embed"
	"fmt"
	"go/ast"
	"io"
	"reflect"
	"strings"
	"text/template"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/iancoleman/strcase"
	"golang.org/x/tools/go/packages"
)

//go:embed templates/*.go.tpl
var TemplatesFS embed.FS

func valueToSchemaType(v reflect.Type) (schema.ValueType, error) {
	k := v.Kind()
	switch k {
	case reflect.String:
		return schema.TypeString, nil
	case reflect.Bool:
		return schema.TypeBool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return schema.TypeInt, nil
	case reflect.Float32, reflect.Float64:
		return schema.TypeFloat, nil
	case reflect.Map:
		return schema.TypeJSON, nil
	case reflect.Struct:
		t := v.PkgPath() + "." + v.Name()
		if t == "time.Time" {
			return schema.TypeTimestamp, nil
		}
		return schema.TypeJSON, nil
	case reflect.Pointer:
		return valueToSchemaType(v.Elem())
	case reflect.Slice:
		switch v.Elem().Kind() {
		case reflect.String:
			return schema.TypeStringArray, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return schema.TypeIntArray, nil
		default:
			return schema.TypeJSON, nil
		}
	default:
		return schema.TypeInvalid, fmt.Errorf("unsupported type: %s", k)
	}
}

type TableOptions func(*TableDefinition)

func WithNameTransformer(transformer func(string) string) TableOptions {
	return func(t *TableDefinition) {
		t.Name = transformer(t.Name)
	}
}

func WithSkipFields(fields []string) TableOptions {
	return func(t *TableDefinition) {
		t.skipFields = fields
	}
}

func WithOverrideColumns(columns []ColumnDefinition) TableOptions {
	return func(t *TableDefinition) {
		t.overrideColumns = columns
	}
}

func WithDescriptionsEnabled() TableOptions {
	return func(t *TableDefinition) {
		t.descriptionsEnabled = true
	}
}

func defaultTransformer(name string) string {
	return strcase.ToSnake(name)
}

func sliceContains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func NewTableFromStruct(name string, obj interface{}, opts ...TableOptions) (*TableDefinition, error) {
	t := TableDefinition{
		Name:            name,
		nameTransformer: defaultTransformer,
	}
	for _, opt := range opts {
		opt(&t)
	}

	e := reflect.ValueOf(obj)
	if e.Kind() == reflect.Pointer {
		e = e.Elem()
	}
	if e.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", e.Kind())
	}

	comments := make(map[string]string)
	if t.descriptionsEnabled {
		comments = readStructComments(e.Type().PkgPath(), e.Type().Name())
	}

	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)
		if sliceContains(t.skipFields, field.Name) {
			continue
		}

		if t.overrideColumns != nil {
			if col := t.overrideColumns.GetByName(t.nameTransformer(field.Name)); col != nil {
				t.Columns = append(t.Columns, *col)
				continue
			}
		}

		columnType, err := valueToSchemaType(field.Type)
		if err != nil {
			return nil, err
		}

		// generate a PathResolver to use by default
		pathResolver := fmt.Sprintf("schema.PathResolver(%q)", field.Name)
		column := ColumnDefinition{
			Name:        t.nameTransformer(field.Name),
			Type:        columnType,
			Resolver:    pathResolver,
			Description: strings.ReplaceAll(comments[field.Name], "`", "'"),
		}
		t.Columns = append(t.Columns, column)
	}
	return &t, nil
}

func (t *TableDefinition) GenerateTemplate(wr io.Writer) error {
	tpl, err := template.New("table.go.tpl").ParseFS(TemplatesFS, "templates/*")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tpl.Execute(wr, t); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}

// type commentReader struct {
// 	pkgPath  string
// 	//comments is a map of type->comment
// 	comments map[string]string
// }

// func newCommentsReader() {

// }

func readStructComments(pkgPath string, structName string) map[string]string {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		panic(err)
	}
	comments := make(map[string]string, 0)
	for _, p := range pkgs {
		for _, f := range p.Syntax {
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					if st, ok := x.Type.(*ast.StructType); ok {
						if x.Name.Name == structName {
							for _, field := range st.Fields.List {
								if len(field.Names) > 0 {
									comments[field.Names[0].Name] = field.Doc.Text()
								}
							}
						}
					}
				}
				return true
			})
		}

	}
	return comments
}
