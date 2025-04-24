package docs

import (
	_ "embed"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/docs/testdata"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/schema.json
var testSchema string

type testTableOptions struct {
	Dummy testdata.DummyTableOptions `json:"dummy,omitempty"`
}

var testTable = &schema.Table{
	Name:        "dummy",
	Description: "This is a dummy table",
}

func TestTableOptionsDescriptionTransformer(t *testing.T) {
	type args struct {
		tableOptions any
		jsonSchema   string
		table        *schema.Table
	}
	tests := []struct {
		name     string
		args     args
		wantDesc string
		wantErr  bool
	}{
		{
			name:     "adds table options to description",
			args:     args{tableOptions: &testTableOptions{}, jsonSchema: testSchema, table: testTable},
			wantDesc: "This is a dummy table\n\n## <a name=\"Table Options\"></a>Table Options\n\n  DummyTableOptions contains configuration for the dummy table\n\n* `filter` (`string`)\n",
		},
		{
			name:     "leaves description unchanged when table doesn't have options",
			args:     args{tableOptions: &testTableOptions{}, jsonSchema: testSchema, table: &schema.Table{Description: "Foobar"}},
			wantDesc: "Foobar",
		},
		{
			name:    "errors out when table options don't match schema",
			wantErr: true,
			args:    args{tableOptions: &testTableOptions{}, jsonSchema: "", table: testTable},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer, err := TableOptionsDescriptionTransformer(tt.args.tableOptions, tt.args.jsonSchema)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, transformer(tt.args.table))
			require.Equal(t, tt.wantDesc, tt.args.table.Description)
		})
	}
}
