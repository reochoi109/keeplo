package scheduler

import (
	"context"
	"time"
)

type Executor interface {
	Execute(ctx context.Context) error
}

type Task struct {
	ID          string
	Executor    Executor
	NextCheckAt time.Time
	Interval    time.Duration
	Index       int
}

type TaskQueue interface {
	Push(task *Task) error
	Pop(ctx context.Context) (*Task, error)
	UpdateTask(ID string, newTime time.Time) error
	RemoveTask(ID string)
	Length() int
	Close()
}
