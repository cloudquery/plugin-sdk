package premium

import (
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/assert"
	"testing"
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
	}

	paidTables := MakeAllTablesPaid(noPaidTables)

	assert.Equal(t, 3, len(paidTables))
	for _, table := range paidTables {
		assert.True(t, table.IsPaid)
	}
}
