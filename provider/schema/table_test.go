package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tableDefinitionTestCases = []tableTestCase{
	{
		Name: "simpleTable",
		Table: &Table{
			Name: "simple_table",
			Columns: []Column{
				{
					Name: "some_string",
					Type: TypeString,
				},
			},
		},
		ExpectedColumnNames: []string{"cq_id", "cq_meta", "some_string"},
		ExpectedHasId:       false,
	},
	{
		Name: "simpleTableWithId",
		Table: &Table{
			Name: "simple_table_with_id",
			Columns: []Column{
				{
					Name: "some_string",
					Type: TypeString,
				},
				{
					Name: "some_int",
					Type: TypeInt,
				},
			},
		},
		ExpectedColumnNames: []string{"cq_id", "cq_meta", "some_string", "some_int"},
		ExpectedHasId:       true,
	},
	{
		Name: "simpleTableWithEmbedded",
		Table: &Table{
			Name: "simple_embedded_table",
			Columns: []Column{
				{
					Name: "some_string",
					Type: TypeString,
				},
				{
					Name: "some_int",
					Type: TypeInt,
				},
				{
					Name: "embedded_some_string",
					Type: TypeString,
				},
				{
					Name: "embedded_some_int",
					Type: TypeInt,
				},
			},
		},
		ExpectedColumnNames: []string{"cq_id", "cq_meta", "some_string", "some_int", "embedded_some_string", "embedded_some_int"},
	},

	{
		Name: "multiEmbeddedTable",
		Table: &Table{
			Name: "multi_embedded_table",
			Columns: []Column{
				{
					Name: "some_int",
					Type: TypeInt,
				},
				{
					Name: "embedded_some_string",
					Type: TypeString,
				},
				{
					Name: "embedded_inner_some_int",
					Type: TypeInt,
				},
			},
		},
		ExpectedColumnNames: []string{"cq_id", "cq_meta", "some_int", "embedded_some_string", "embedded_inner_some_int"},
	},
	{
		Name: "simpleTableWithEmbedded",
		Table: &Table{
			Name: "simple_embedded_table",
			Columns: []Column{
				{
					Name: "some_string",
					Type: TypeString,
				},
				{
					Name: "some_int",
					Type: TypeInt,
				},
				{
					Name: "some_string_no_prefix",
					Type: TypeString,
				},
				{
					Name: "some_int_no_prefix",
					Type: TypeInt,
				},
			},
		},
		ExpectedColumnNames: []string{"cq_id", "cq_meta", "some_string", "some_int", "some_string_no_prefix", "some_int_no_prefix"},
	},
}

type tableTestCase struct {
	Name                string
	Table               *Table
	ExpectedColumnNames []string
	ExpectedHasId       bool
}

func TestTableDefinitionUseCases(t *testing.T) {
	for _, c := range tableDefinitionTestCases {
		assert.Equal(t, c.ExpectedColumnNames, PostgresDialect{}.Columns(c.Table).Names(), "failed case %s", c.Name)
	}
}

func TestTableSignatures(t *testing.T) {
	signaturesSoFar := make(map[string]struct{})
	for i, c := range tableDefinitionTestCases {
		s := c.Table.Signature(PostgresDialect{})
		if _, ok := signaturesSoFar[s]; ok {
			t.Fatalf("duplicate signature encountered (index #%d)", i+1)
		}
		signaturesSoFar[s] = struct{}{}

		newT := *c.Table
		newT.Columns = append(newT.Columns, Column{Name: "extra_column_in_test", Type: TypeString})
		newSig := newT.Signature(PostgresDialect{})
		assert.NotEqual(t, s, newSig)

		newSig2 := newT.Signature(PostgresDialect{})
		assert.Equal(t, newSig, newSig2)

		newT.Serial = "some_test_serial"
		newSig3 := newT.Signature(PostgresDialect{})
		assert.NotEqual(t, newSig, newSig3)

		newT = *c.Table
		newT.Columns = newT.Columns[1:]
		newSig4 := newT.Signature(PostgresDialect{})
		assert.NotEqual(t, newSig, newSig4)
	}
}
