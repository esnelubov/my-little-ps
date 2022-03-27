package partitioner

import (
	"container/heap"
	"my-little-ps/wlt_balancer/partitioner/partition"
	"my-little-ps/wlt_balancer/partitioner/priority_queue"
)

// KarmarkarKarp
// https://en.wikipedia.org/wiki/Largest_differencing_method
// https://stackoverflow.com/questions/68421928/how-to-reconstruct-partitions-in-the-karmarkar-karp-heuristic-multi-way-partitio
func KarmarkarKarp(numbers []int64, k int) (partitions [][]int64) {
	if k >= len(numbers) {
		partitions = make([][]int64, len(numbers))

		for i, n := range numbers {
			partitions[i] = []int64{n}
		}

		return
	}

	pq := make(priority_queue.Type, len(numbers))

	for i, n := range numbers {
		pq[i] = partition.New(n, k)
	}
	heap.Init(&pq)

	for pq.Len() > 1 {
		a := heap.Pop(&pq).(*partition.Type)
		b := heap.Pop(&pq).(*partition.Type)
		a.Merge(b)
		heap.Push(&pq, a)
	}

	partitions = make([][]int64, k)

	for i, s := range pq[0].Subsets {
		partitions[i] = s.Numbers
	}

	return
}
