package scheduler

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"
)

type InMemoryQueue struct {
	queue   PriorityQueue
	mu      sync.Mutex
	cond    *sync.Cond
	taskMap map[string]*Task // 모니터 ID 기준 Task 관리
	closed  bool
}

func NewInMemoryQueue() *InMemoryQueue {
	q := &InMemoryQueue{
		queue:   make(PriorityQueue, 0),
		taskMap: make(map[string]*Task),
	}
	q.cond = sync.NewCond(&q.mu)
	heap.Init(&q.queue)
	return q
}

func (q *InMemoryQueue) Push(task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.taskMap[task.ID]; exists {
		return fmt.Errorf("task already exists")
	}

	heap.Push(&q.queue, task)
	q.taskMap[task.ID] = task
	q.cond.Signal()
	return nil
}

func (q *InMemoryQueue) Pop(ctx context.Context) (*Task, error) {
	q.mu.Lock()

	for {
		if q.closed {
			q.mu.Unlock()
			return nil, fmt.Errorf("queue closed")
		}

		if q.queue.Len() == 0 {
			q.mu.Unlock()

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				q.mu.Lock()
				q.cond.Wait()
				continue
			}
		}
		task := q.queue.Peek()
		now := time.Now()

		// 아직 실행 시간이 안 됨 → unlock 후 wait
		if task.NextCheckAt.After(now) {
			wait := task.NextCheckAt.Sub(now)
			q.mu.Unlock()

			timer := time.NewTimer(wait)
			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return nil, ctx.Err()

			case <-timer.C:
				q.mu.Lock()
				continue
			}
		}
		item := heap.Pop(&q.queue).(*Task)
		delete(q.taskMap, item.ID)
		q.mu.Unlock()
		return item, nil
	}
}

func (q *InMemoryQueue) UpdateTask(monitorID string, newTime time.Time) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	task, ok := q.taskMap[monitorID]
	if !ok {
		return fmt.Errorf("task not found")
	}

	task.NextCheckAt = newTime
	heap.Fix(&q.queue, task.Index)
	q.cond.Signal()
	return nil
}

func (q *InMemoryQueue) RemoveTask(monitorID string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	task, ok := q.taskMap[monitorID]
	if !ok {
		return
	}

	heap.Remove(&q.queue, task.Index)
	delete(q.taskMap, monitorID)
}

func (q *InMemoryQueue) Length() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.queue.Len()
}

func (q *InMemoryQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
}
