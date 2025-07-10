package scheduler

import (
	"container/heap"
	"sync"
)

type InMemoryQueue struct {
	queue PriorityQueue
	mu    sync.Mutex
	cond  *sync.Cond
}

func NewInMemoryQueue() *InMemoryQueue {
	q := &InMemoryQueue{
		queue: make(PriorityQueue, 0),
	}
	q.cond = sync.NewCond(&q.mu)
	heap.Init(&q.queue)
	return q
}
