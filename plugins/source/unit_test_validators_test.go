package source

import (
	"fmt"
	"testing"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/stretchr/testify/assert"
)

func TestPluralNamesValidator(t *testing.T) {
	testTable := []struct {
		tableName string
		err       error
	}{
		{
			tableName: "test_tables",
			err:       nil,
		},
		{
			tableName: "test_table",
			err:       fmt.Errorf("invalid table name: test_table. must be plural"),
		},
	}
	for _, tc := range testTable {
		tables := []*schema.Table{
			{
				Name:      tc.tableName,
				Transform: transformers.TransformWithStruct(&testTable),
			},
		}

		plugin := NewPlugin("testSourcePlugin", "1.0.0", tables, newTestExecutionClient)
		err := ValidateTableNamePlural(plugin, nil)
		assert.Equal(t, tc.err, err)
	}
}
