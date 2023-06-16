package plugin

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/apache/arrow/go/v13/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

type Validator func(t *testing.T, plugin *Plugin, resources []message.Message)

func TestPluginSync(t *testing.T, plugin *Plugin, spec []byte, options SyncOptions, opts ...TestPluginOption) {
	t.Helper()

	o := &testPluginOptions{
		parallel:   true,
		validators: []Validator{validatePlugin},
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.parallel {
		t.Parallel()
	}

	resourcesChannel := make(chan message.Message)
	var syncErr error

	if err := plugin.Init(context.Background(), spec); err != nil {
		t.Fatal(err)
	}

	go func() {
		defer close(resourcesChannel)
		syncErr = plugin.Sync(context.Background(), options, resourcesChannel)
	}()

	syncedResources := make([]message.Message, 0)
	for resource := range resourcesChannel {
		syncedResources = append(syncedResources, resource)
	}
	if syncErr != nil {
		t.Fatal(syncErr)
	}
	for _, validator := range o.validators {
		validator(t, plugin, syncedResources)
	}
}

type TestPluginOption func(*testPluginOptions)

func WithTestPluginNoParallel() TestPluginOption {
	return func(f *testPluginOptions) {
		f.parallel = false
	}
}

func WithTestPluginAdditionalValidators(v Validator) TestPluginOption {
	return func(f *testPluginOptions) {
		f.validators = append(f.validators, v)
	}
}

type testPluginOptions struct {
	parallel   bool
	validators []Validator
}

func getTableResources(t *testing.T, table *schema.Table, messages []message.Message) []arrow.Record {
	t.Helper()

	tableResources := make([]arrow.Record, 0)
	for _, msg := range messages {
		switch v := msg.(type) {
		case *message.Insert:
			md := v.Record.Schema().Metadata()
			tableName, ok := md.GetValue(schema.MetadataTableName)
			if !ok {
				t.Errorf("Expected table name to be set in metadata")
			}
			if tableName == table.Name {
				tableResources = append(tableResources, v.Record)
			}
		default:
			t.Errorf("Unexpected message type %T", v)
		}
	}

	return tableResources
}

func validateTable(t *testing.T, table *schema.Table, messages []message.Message) {
	t.Helper()
	tableResources := getTableResources(t, table, messages)
	if len(tableResources) == 0 {
		t.Errorf("Expected table %s to be synced but it was not found", table.Name)
		return
	}
	validateResources(t, table, tableResources)
}

func validatePlugin(t *testing.T, plugin *Plugin, resources []message.Message) {
	t.Helper()
	tables, err := plugin.Tables(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, table := range tables.FlattenTables() {
		validateTable(t, table, resources)
	}
}

// Validates that every column has at least one non-nil value.
// Also does some additional validations.
func validateResources(t *testing.T, table *schema.Table, resources []arrow.Record) {
	t.Helper()

	// A set of column-names that have values in at least one of the resources.
	columnsWithValues := make([]bool, len(table.Columns))

	for _, resource := range resources {
		for _, arr := range resource.Columns() {
			for i := 0; i < arr.Len(); i++ {
				if arr.IsValid(i) {
					columnsWithValues[i] = true
				}
			}
		}
	}

	// Make sure every column has at least one value.
	for i, hasValue := range columnsWithValues {
		col := table.Columns[i]
		emptyExpected := col.Name == "_cq_parent_id" && table.Parent == nil
		if !hasValue && !emptyExpected && !col.IgnoreInTests {
			t.Errorf("table: %s column %s has no values", table.Name, table.Columns[i].Name)
		}
	}
}

func RecordDiff(l arrow.Record, r arrow.Record) string {
	var sb strings.Builder
	if l.NumCols() != r.NumCols() {
		return fmt.Sprintf("different number of columns: %d vs %d", l.NumCols(), r.NumCols())
	}
	if l.NumRows() != r.NumRows() {
		return fmt.Sprintf("different number of rows: %d vs %d", l.NumRows(), r.NumRows())
	}
	for i := 0; i < int(l.NumCols()); i++ {
		edits, err := array.Diff(l.Column(i), r.Column(i))
		if err != nil {
			panic(fmt.Sprintf("left: %v, right: %v, error: %v", l.Column(i).DataType(), r.Column(i).DataType(), err))
		}
		diff := edits.UnifiedDiff(l.Column(i), r.Column(i))
		if diff != "" {
			sb.WriteString(l.Schema().Field(i).Name)
			sb.WriteString(": ")
			sb.WriteString(diff)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
