package schema

import (
	"errors"
	"fmt"
)

type TableValidator interface {
	Validate(t *Table) error
}

type LengthTableValidator struct{}

const (
	maxTableName  = 63 // maximum allowed identifier length is 63 bytes https://www.postgresql.org/docs/13/limits.html
	maxColumnName = 63
)

var defaultValidators = []TableValidator{
	LengthTableValidator{},
}

func ValidateTable(t *Table) error {
	for _, validator := range defaultValidators {
		return validator.Validate(t)
	}
	return nil
}

func validateTableAttributesNameLength(t *Table) error {
	// validate table name
	if len(t.Name) > maxTableName {
		return errors.New("table name has exceeded max length")
	}

	// validate table columns
	for _, col := range t.Columns.Names() {
		if len(col) > maxColumnName {
			return fmt.Errorf("column name %s has exceeded max length", col)
		}
	}

	// validate table relations
	for _, rel := range t.Relations {
		err := validateTableAttributesNameLength(rel)
		if err != nil {
			return err
		}
	}
	return nil
}

func (LengthTableValidator) Validate(t *Table) error {
	return validateTableAttributesNameLength(t)
}
