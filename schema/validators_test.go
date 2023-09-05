package schema

import (
	"testing"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/stretchr/testify/assert"
)

func TestTableValidators(t *testing.T) {
	var testTableValidators = Table{
		Name: "test_table_validator",
		Columns: []Column{
			{
				Name: "zero_bool",
				Type: arrow.FixedWidthTypes.Boolean,
			},
			{
				Name: "zero_int",
				Type: arrow.PrimitiveTypes.Int64,
			},
			{
				Name: "not_zero_bool",
				Type: arrow.FixedWidthTypes.Boolean,
			},
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
