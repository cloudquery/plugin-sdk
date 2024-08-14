package queue

import (
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

type testClient struct {
	id string
}

func (tc *testClient) ID() string {
	return tc.id
}

func TestNewConcurrentQueue(t *testing.T) {
	tableClients := []TableClientPair{
		{
			Table:  &schema.Table{Name: "table-1"},
			Client: &testClient{id: "client-1"},
		},
		{
			Table:  &schema.Table{Name: "table-2"},
			Client: &testClient{id: "client-2"},
		},
	}
	queue := NewConcurrentQueue(tableClients)
	require.Equal(t, len(tableClients), queue.queue.Len())
}

func TestConcurrentQueuePushPop(t *testing.T) {
	tableClients := []TableClientPair{
		{
			Table:  &schema.Table{Name: "table-1"},
			Client: &testClient{id: "client-1"},
		},
		{
			Table:  &schema.Table{Name: "table-2"},
			Client: &testClient{id: "client-2"},
		},
	}
	queue := NewConcurrentQueue(tableClients)

	// Simulate table 3 as a child of table 1
	table3Count := 5
	for i := 0; i < table3Count; i++ {
		queue.Push(TableClientPair{
			Table:  &schema.Table{Name: "table-3"},
			Client: &testClient{id: "client-1"},
		})
	}

	// Simulate table 4 as a child of table 2
	table4Count := 10
	for i := 0; i < table4Count; i++ {
		queue.Push(TableClientPair{
			Table:  &schema.Table{Name: "table-4"},
			Client: &testClient{id: "client-2"},
		})
	}

	gotClients := make([]TableClientPair, 0)
	for {
		if queue.queue.Len() == 0 {
			break
		}
		gotClients = append(gotClients, *queue.Pop())
	}
	require.Equal(t, len(tableClients)+table3Count+table4Count, len(gotClients))

	require.Equal(t, "table-1", gotClients[0].Table.Name)
	require.Equal(t, "table-2", gotClients[1].Table.Name)

	// Priority declines as there are more items in the queue
	// We expect to get 5 of table-3 and table-4, then the rest of table-4
	require.Equal(t, "table-4", gotClients[2].Table.Name)
	require.Equal(t, "table-3", gotClients[3].Table.Name)
	require.Equal(t, "table-3", gotClients[4].Table.Name)
	require.Equal(t, "table-4", gotClients[5].Table.Name)
	require.Equal(t, "table-3", gotClients[6].Table.Name)
	require.Equal(t, "table-4", gotClients[7].Table.Name)
	require.Equal(t, "table-4", gotClients[8].Table.Name)
	require.Equal(t, "table-3", gotClients[9].Table.Name)
	require.Equal(t, "table-4", gotClients[10].Table.Name)
	require.Equal(t, "table-3", gotClients[11].Table.Name)
	require.Equal(t, "table-4", gotClients[12].Table.Name)
	require.Equal(t, "table-4", gotClients[13].Table.Name)
	require.Equal(t, "table-4", gotClients[14].Table.Name)
	require.Equal(t, "table-4", gotClients[15].Table.Name)
	require.Equal(t, "table-4", gotClients[16].Table.Name)
}
