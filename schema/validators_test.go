package schema

import (
	"testing"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/stretchr/testify/assert"
)

func TestTableValidators(t *testing.T) {
	var testTableValidators = Table{
		Name: "test_table_validator",
		Columns: []Column{
			{Field: arrow.Field{Name: "zero_bool", Type: arrow.FixedWidthTypes.Boolean, Nullable: true}},
			{Field: arrow.Field{Name: "zero_int", Type: arrow.PrimitiveTypes.Int64, Nullable: true}},
			{Field: arrow.Field{Name: "not_zero_bool", Type: arrow.FixedWidthTypes.Boolean, Nullable: true}},
		},
	}

	// table has passed all validators
	err := ValidateTable(&testTableValidators)
	assert.Nil(t, err)

	// table name is too long
	tableWithLongName := testTableValidators
	tableWithLongName.Name = "WithLongNametableWithLongNametableWithLongNametableWithLongNamet"
	err = ValidateTable(&tableWithLongName)
	assert.Error(t, err)

	// column name is too long
	tableWithLongColumnName := testTableValidators
	tableWithLongName.Columns[0].Name = "tableWithLongColumnNametableWithLongColumnNametableWithLongColumnName"
	err = ValidateTable(&tableWithLongColumnName)
	assert.Error(t, err)
}
