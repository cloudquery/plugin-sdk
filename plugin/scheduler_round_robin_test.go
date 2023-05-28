package plugin

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v3/schema"
)

func TestRoundRobinInterleave(t *testing.T) {
	table1 := &schema.Table{Name: "test_table"}
	table2 := &schema.Table{Name: "test_table2"}
	client1 := &testExecutionClient{}
	client2 := &testExecutionClient{}
	client3 := &testExecutionClient{}
	cases := []struct {
		name                  string
		tables                schema.Tables
		preInitialisedClients [][]schema.ClientMeta
		want                  []tableClient
	}{
		{
			name:                  "single table",
			tables:                schema.Tables{table1},
			preInitialisedClients: [][]schema.ClientMeta{{client1}},
			want:                  []tableClient{{table: table1, client: client1}},
		},
		{
			name:                  "two tables with different clients",
			tables:                schema.Tables{table1, table2},
			preInitialisedClients: [][]schema.ClientMeta{{client1}, {client1, client2}},
			want: []tableClient{
				{table: table1, client: client1},
				{table: table2, client: client1},
				{table: table2, client: client2},
			},
		},
		{
			name:                  "two tables with different clients",
			tables:                schema.Tables{table1, table2},
			preInitialisedClients: [][]schema.ClientMeta{{client1, client3}, {client1, client2}},
			want: []tableClient{
				{table: table1, client: client1},
				{table: table2, client: client1},
				{table: table1, client: client3},
				{table: table2, client: client2},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := roundRobinInterleave(tc.tables, tc.preInitialisedClients)
			if len(got) != len(tc.want) {
				t.Fatalf("got %d tableClients, want %d", len(got), len(tc.want))
			}
			for i := range got {
				if got[i].table != tc.want[i].table {
					t.Errorf("got table %v, want %v", got[i].table, tc.want[i].table)
				}
				if got[i].client != tc.want[i].client {
					t.Errorf("got client %v, want %v", got[i].client, tc.want[i].client)
				}
			}
		})
	}
}
