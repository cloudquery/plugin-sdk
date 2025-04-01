package transformers

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

func TestTransformTablesErrorOnPKFieldsAndPKComponentFields(t *testing.T) {
	type testStructWithPKComponent struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	testTable := &schema.Table{
		Name:      "test",
		Transform: TransformWithStruct(testStructWithPKComponent{}, WithPrimaryKeys("id")),
		Columns: []schema.Column{
			{Name: "name", Type: arrow.BinaryTypes.String, PrimaryKeyComponent: true},
		},
	}

	err := TransformTables([]*schema.Table{testTable})
	require.Error(t, err, "primary keys and primary key components cannot both be set for table \"test\"")
}
