package premium

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/assert"
)

func TestContainsPaidTables(t *testing.T) {
	noPaidTables := schema.Tables{
		&schema.Table{Name: "table1", IsPaid: false},
		&schema.Table{Name: "table2", IsPaid: false},
		&schema.Table{Name: "table3", IsPaid: false},
	}

	paidTables := schema.Tables{
		&schema.Table{Name: "table1", IsPaid: false},
		&schema.Table{Name: "table2", IsPaid: true},
		&schema.Table{Name: "table3", IsPaid: false},
	}

	assert.False(t, ContainsPaidTables(noPaidTables), "no paid tables")
	assert.True(t, ContainsPaidTables(paidTables), "paid tables")
}

func TestMakeAllTablesPaid(t *testing.T) {
	noPaidTables := schema.Tables{
		&schema.Table{Name: "table1", IsPaid: false},
		&schema.Table{Name: "table2", IsPaid: false},
		&schema.Table{Name: "table3", IsPaid: false},
		&schema.Table{Name: "table_with_relations", IsPaid: false, Relations: schema.Tables{
			&schema.Table{Name: "relation_table", IsPaid: false},
		}},
	}

	paidTables := MakeAllTablesPaid(noPaidTables)

	assert.Equal(t, 4, len(paidTables))
	assert.Equal(t, 5, len(paidTables.FlattenTables()))
	assertAllArePaid(t, paidTables)
}

func assertAllArePaid(t *testing.T, tables schema.Tables) {
	t.Helper()
	for _, table := range tables {
		assert.True(t, table.IsPaid)
		assertAllArePaid(t, table.Relations)
	}
}
