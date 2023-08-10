package docs

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func TestDefaultTitleTransformer(t *testing.T) {
	type args struct {
		table *schema.Table
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default",
			args: args{
				table: &schema.Table{
					Name: "test_table",
				},
			},
			want: "Test Table",
		},
		{
			name: "existing title",
			args: args{
				table: &schema.Table{
					Name:  "test_table",
					Title: "My Custom Title",
				},
			},
			want: "My Custom Title",
		},
		{
			name: "abbreviations",
			args: args{
				table: &schema.Table{
					Name: "test_acls",
				},
			},
			want: "Test ACLs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultTitleTransformer(tt.args.table); got != tt.want {
				t.Errorf("DefaultTitleTransformer() = %v, want %v", got, tt.want)
			}
		})
	}
}
