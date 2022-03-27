package priority_queue

import "my-little-ps/wlt_balancer/partitioner/partition"

// A Type implements heap.Interface and holds Partitions.
type Type []*partition.Type

func (pq Type) Len() int { return len(pq) }

func (pq Type) Less(i, j int) bool {
	return pq[i].Priority() > pq[j].Priority()
}

func (pq Type) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *Type) Push(x interface{}) {
	item := x.(*partition.Type)
	*pq = append(*pq, item)
}

func (pq *Type) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
