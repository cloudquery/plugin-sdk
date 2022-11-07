package codegen

import (
	"bytes"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/stretchr/testify/require"
)

func TestGenerateTemplate(t *testing.T) {
	type args struct {
		table *TableDefinition
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "should add comma between relations",
			args: args{
				table: &TableDefinition{
					Name:      "with relations",
					Relations: []string{"relation1", "relation2"},
				},
			},
		},
		{
			name: "should add ignore_in_tests to columns",
			args: args{
				table: &TableDefinition{
					Name: "with relations",
					Columns: []ColumnDefinition{
						{Name: "ignore_in_tests", Type: schema.TypeString, IgnoreInTests: true},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			err := tt.args.table.GenerateTemplate(buf)
			require.NoError(t, err)
			cupaloy.SnapshotT(t, buf.Bytes())
		})
	}
}
