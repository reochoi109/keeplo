package scheduler

import (
	"context"
	"errors"
	"keeplo/pkg/logger"
	"sync"
	"time"

	"go.uber.org/zap"
)

var ErrQueueNotFound = errors.New("not found queue")

type scheduler struct {
	queues map[string]TaskQueue
	lock   sync.RWMutex
}

var Scheduler *scheduler

func NewScheduler() {
	Scheduler = &scheduler{
		queues: make(map[string]TaskQueue),
	}
}

func AddQueue(name string, queue TaskQueue) {
	Scheduler.lock.Lock()
	defer Scheduler.lock.Unlock()

	log := logger.Log
	if _, exists := Scheduler.queues[name]; exists {
		log.Warn("Queue already exists", zap.String("queue", name))
		return
	}
	Scheduler.queues[name] = queue
	log.Info("Queue added to scheduler", zap.String("queue", name))
	go startQueueWorker(name, queue)
}

// 큐에 Task 등록
func RegisterTask(ctx context.Context, queueName string, task *Task) error {
	Scheduler.lock.RLock()
	queue, ok := Scheduler.queues[queueName]
	Scheduler.lock.RUnlock()
	log := logger.Log
	if !ok {
		log.Error("Queue not found", zap.String("queue", queueName))
		return ErrQueueNotFound
	}

	if err := queue.Push(task); err != nil {
		log.Warn("Task push failed, trying update", zap.String("task_id", task.ID), zap.Error(err))
		return queue.UpdateTask(task.ID, task.NextCheckAt)
	}

	log.Debug("Task registered successfully", zap.String("task_id", task.ID), zap.String("queue", queueName))
	return nil
}

// 업데이트
func UpdateTask(ctx context.Context, queueName string, updated *Task) error {
	log := logger.WithContext(ctx)
	Scheduler.lock.Lock()
	queue, ok := Scheduler.queues[queueName]
	Scheduler.lock.Unlock()

	if !ok {
		log.Error("Queue not found", zap.String("queue", queueName))
		return ErrQueueNotFound
	}

	// update task는 NextCheckAt만 갱신하므로 다른 필드는 remove -> Push 방식으로 처리
	queue.RemoveTask(updated.ID)
	if err := queue.Push(updated); err != nil {
		log.Error("Failed to update task", zap.String("task_id", updated.ID), zap.Error(err))
		return err
	}

	log.Info("Task updated", zap.String("task_id", updated.ID), zap.String("queue", queueName))
	return nil
}

// Task 제거
func RemoveTask(queueName, taskID string) {
	Scheduler.lock.RLock()
	queue, ok := Scheduler.queues[queueName]
	Scheduler.lock.RUnlock()
	log := logger.Log
	if !ok {
		log.Warn("Queue not found when removing task", zap.String("queue", queueName))
		return
	}

	queue.RemoveTask(taskID)
	log.Info("Task removed from queue", zap.String("task_id", taskID), zap.String("queue", queueName))
}

// 각 큐별 고루틴 루프
func startQueueWorker(queueName string, queue TaskQueue) {
	ctx := context.Background()
	log := logger.Log
	for {
		task, err := queue.Pop(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Info("Queue shutting down", zap.String("queue", queueName))
				return
			}
			log.Error("Queue pop error", zap.String("queue", queueName), zap.Error(err))
			time.Sleep(500 * time.Millisecond)
			continue
		}
		go handleTask(queueName, task)
	}
}

func handleTask(queueName string, task *Task) {
	log := logger.Log

	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered from panic in task", zap.Any("recover", r), zap.String("task_id", task.ID))
		}
	}()

	if task.Executor != nil {
		if err := task.Executor.Execute(context.Background(), task.Payload); err != nil {
			log.Error("Task execution failed", zap.String("task_id", task.ID), zap.Error(err))
		}
	}

	next := time.Now().Add(task.Interval)
	err := RegisterTask(context.Background(), queueName, &Task{
		ID:          task.ID,
		Executor:    task.Executor,
		Payload:     task.Payload,
		Interval:    task.Interval,
		NextCheckAt: next,
	})
	if err != nil {
		log.Error("Failed to reschedule task", zap.String("task_id", task.ID), zap.Error(err))
		return
	}
}
