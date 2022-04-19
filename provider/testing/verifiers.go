package testing

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"github.com/cloudquery/faker/v3/support/slice"
	"github.com/georgysavva/scany/pgxscan"
)

type Row map[string]interface{}

// VerifyRowPredicateInTable is a base verifier accepting single row verifier for specific table from schema
func VerifyRowPredicateInTable(tableName string, rowVerifier func(*testing.T, Row)) Verifier {
	var verifier Verifier
	verifier = func(t *testing.T, table *schema.Table, conn pgxscan.Querier, shouldSkipIgnoreInTest bool) {
		if tableName == table.Name {
			rows := getRows(t, conn, table, shouldSkipIgnoreInTest)
			for _, row := range rows {
				rowVerifier(t, row)
			}
		}
		for _, r := range table.Relations {
			verifier(t, r, conn, shouldSkipIgnoreInTest)
		}
	}
	return verifier
}

func getRows(t *testing.T, conn pgxscan.Querier, table *schema.Table, shouldSkipIgnoreInTest bool) []Row {
	if shouldSkipIgnoreInTest && table.IgnoreInTests {
		t.Skipf("table %s marked as IgnoreInTest. Skipping...", table.Name)
	}

	var rows []Row
	err := pgxscan.Get(
		context.Background(),
		conn,
		&rows,
		fmt.Sprintf("select json_agg(%[1]s) from %[1]s", table.Name),
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range table.Columns {
		if shouldSkipIgnoreInTest && c.IgnoreInTests {
			for _, row := range rows {
				delete(row, c.Name)
			}
		}
	}

	return rows
}

// VerifyNoEmptyColumnsExcept verifies that for each row in table its columns are not empty except passed
func VerifyNoEmptyColumnsExcept(tableName string, except ...string) Verifier {
	return VerifyRowPredicateInTable(tableName, func(t *testing.T, row Row) {
		for k, v := range row {
			if !slice.Contains(except, k) && v == nil {
				t.Fatal("VerifyNoEmptyColumnsExcept failed: illegal row found")
			}
		}
	})
}

// VerifyAtMostOneOf verifies that for each row in table at most one column from oneof is not empty
func VerifyAtMostOneOf(tableName string, oneof ...string) Verifier {
	return VerifyRowPredicateInTable(tableName, func(t *testing.T, row Row) {
		cnt := 0
		for _, k := range oneof {
			v, ok := row[k]
			if !ok {
				t.Fatalf("VerifyAtMostOneOf failed: column %s doesn't exist", k)
			}
			if v != nil {
				cnt++
			}
		}

		if cnt > 1 {
			t.Fatal("VerifyAtMostOneOf failed: illegal row found")
		}
	})
}

// VerifyAtLeastOneRow verifies that main table from schema has at least one row
func VerifyAtLeastOneRow() Verifier {
	return func(t *testing.T, table *schema.Table, conn pgxscan.Querier, _ bool) {
		rows, err := conn.Query(context.Background(), fmt.Sprintf("select * from %s;", table.Name))
		if err != nil {
			t.Fatal(err)
		}
		if !rows.Next() {
			t.Fatal("VerifyAtLeastOneRow failed: table is empty")
		}

		rows.Close()
	}
}
