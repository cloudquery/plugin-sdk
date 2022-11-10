package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableValidators(t *testing.T) {
	var testTableValidators = Table{
		Name: "test_table_validator",
		Columns: []Column{
			{
				Name: "zero_bool",
				Type: TypeBool,
			},
			{
				Name: "zero_int",
				Type: TypeInt,
			},
			{
				Name: "not_zero_bool",
				Type: TypeBool,
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
