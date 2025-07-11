package scheduler

import (
	"context"
	"fmt"
	"time"
)

type Scheduler struct {
	queue TaskQueue
}

func NewScheduler(queue TaskQueue) *Scheduler {
	return &Scheduler{
		queue: queue,
	}
}

func (s *Scheduler) RegisterTask(ctx context.Context, id string, interval time.Duration) error {
	task := &Task{
		ID:          id,
		NextCheckAt: time.Now().Add(interval),
		Interval:    interval,
	}

	if err := s.queue.Push(task); err != nil {
		return s.queue.UpdateTask(id, task.NextCheckAt)
	}
	return nil
}

func (s *Scheduler) Start(ctx context.Context) {
	for {
		task, err := s.queue.Pop(ctx)
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("service closing...")
				return
			}
			fmt.Println("queue pop error:", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		go s.handleTask(ctx, task)
	}
}

func (s *Scheduler) handleTask(ctx context.Context, task *Task) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in task:", r)
		}
	}()

	// 실제 헬스체크 로직은 이곳에서 처리
	fmt.Printf("[CHECK] ID: %s | Time: %s\n", task.ID, time.Now().Format("15:04:05"))

	// 다음 체크 시간 계산
	next := time.Now().Add(task.Interval)
	// ing....

	// 재등록 (Push → 중복 시 Update)
	err := s.queue.Push(&Task{
		ID:          task.ID,
		NextCheckAt: next,
		Interval:    task.Interval,
	})
	if err != nil {
		if updateErr := s.queue.UpdateTask(task.ID, next); updateErr != nil {
			fmt.Printf("[ERROR] Failed to update task %s: %v\n", task.ID, updateErr)
		}
	}
}
