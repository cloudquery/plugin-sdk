package queue

import (
	"container/heap"
	"sync"

	"github.com/cloudquery/plugin-sdk/v4/schema"
)

const (
	maxPriority    = 10
	defaultMaxSize = 50000
)

type TableClientPair struct {
	Table  *schema.Table
	Client schema.ClientMeta
	Parent *schema.Resource
}

type Item struct {
	value    TableClientPair
	priority int
	index    int
}

type priorityQueue []*Item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *priorityQueue) Update(item *Item, value TableClientPair, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func newQueue(tableClients []TableClientPair) *priorityQueue {
	pq := make(priorityQueue, len(tableClients))
	for i, tc := range tableClients {
		pq[i] = &Item{
			value:    tc,
			priority: maxPriority,
			index:    i,
		}
	}
	heap.Init(&pq)
	return &pq
}

type ConcurrentQueue struct {
	queue                     *priorityQueue
	lock                      sync.Mutex
	itemCountByTableAndClient map[string]int
}

func tableClientID(tableClient TableClientPair) string {
	return tableClient.Client.ID() + ":" + tableClient.Table.Name
}

func NewConcurrentQueue(tableClients []TableClientPair) *ConcurrentQueue {
	cq := ConcurrentQueue{
		queue:                     newQueue(tableClients),
		itemCountByTableAndClient: make(map[string]int),
	}
	for _, tc := range tableClients {
		if _, ok := cq.itemCountByTableAndClient[tc.Client.ID()+":"+tc.Table.Name]; ok {
			cq.itemCountByTableAndClient[tableClientID(tc)]++
		} else {
			cq.itemCountByTableAndClient[tableClientID(tc)] = 1
		}
	}
	return &cq
}

func (cq *ConcurrentQueue) getPriority(tableClient TableClientPair) int {
	priority := maxPriority - cq.itemCountByTableAndClient[tableClientID(tableClient)]
	return priority
}

func (cq *ConcurrentQueue) Push(tableClient TableClientPair) {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	if _, ok := cq.itemCountByTableAndClient[tableClientID(tableClient)]; !ok {
		cq.itemCountByTableAndClient[tableClientID(tableClient)] = 0
	}
	heap.Push(cq.queue, &Item{
		value:    tableClient,
		priority: cq.getPriority(tableClient),
	})
	cq.itemCountByTableAndClient[tableClientID(tableClient)]++
}

func (cq *ConcurrentQueue) Pop() *TableClientPair {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	if cq.queue.Len() == 0 {
		return nil
	}
	tableClient := heap.Pop(cq.queue).(*Item).value
	cq.itemCountByTableAndClient[tableClient.Client.ID()+":"+tableClient.Table.Name]--
	return &tableClient
}
