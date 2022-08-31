package schema

import (
	"testing"
)

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

func TestTableValidators(t *testing.T) {
	// table has passed all validators
	err := ValidateTable(&testTableValidators)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// table name is too long
	tableWithLongName := testTableValidators
	tableWithLongName.Name = "WithLongNametableWithLongNametableWithLongNametableWithLongNamet"
	err = ValidateTable(&tableWithLongName)
	if err == nil {
		t.Errorf("expected error but got none")
	}

	// column name is too long
	tableWithLongColumnName := testTableValidators
	tableWithLongName.Columns[0].Name = "tableWithLongColumnNametableWithLongColumnNametableWithLongColumnName"
	err = ValidateTable(&tableWithLongColumnName)
	if err == nil {
		t.Errorf("expected error but got none")
	}
}
