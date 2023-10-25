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
