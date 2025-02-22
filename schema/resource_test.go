package schema

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/stretchr/testify/require"
)

func TestResource_Validate(t *testing.T) {
	tests := []struct {
		name        string
		resource    *Resource
		valueSetter func(resource *Resource) error
		err         error
	}{
		{
			name:     "valid resource without primary keys or primary key components",
			resource: NewResourceData(&Table{Name: "test", Columns: ColumnList{{Name: "col1", Type: arrow.BinaryTypes.String}}}, nil, nil),
			err:      nil,
		},
		{
			name:     "valid resource with primary keys",
			resource: NewResourceData(&Table{Name: "test", Columns: ColumnList{{Name: "col1", Type: arrow.BinaryTypes.String, PrimaryKey: true}}}, nil, nil),
			err:      nil,
			valueSetter: func(resource *Resource) error {
				return resource.Set("col1", "test")
			},
		},
		{
			name:     "valid resource with primary key components",
			resource: NewResourceData(&Table{Name: "test", Columns: ColumnList{{Name: "col1", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true}}}, nil, nil),
			err:      nil,
			valueSetter: func(resource *Resource) error {
				return resource.Set("col1", "test")
			},
		},
		{
			name:     "invalid resource with primary keys",
			resource: NewResourceData(&Table{Name: "test", Columns: ColumnList{{Name: "col1", Type: arrow.BinaryTypes.String, PrimaryKey: true}}}, nil, nil),
			err:      &PKError{MissingPKs: []string{"col1"}},
		},
		{
			name:     "invalid resource with primary key components",
			resource: NewResourceData(&Table{Name: "test", Columns: ColumnList{{Name: "col1", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true}}}, nil, nil),
			err:      &PKComponentError{MissingPKComponents: []string{"col1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valueSetter != nil {
				require.NoError(t, tt.valueSetter(tt.resource))
			}
			validationError := tt.resource.Validate()
			require.Equal(t, tt.err, validationError)
		})
	}
}
