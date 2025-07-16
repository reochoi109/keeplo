package scheduler

import (
	"context"
	"time"
)

type Executor interface {
	Execute(ctx context.Context, playload any) error
}

type Task struct {
	ID          string
	Executor    Executor
	NextCheckAt time.Time
	Interval    time.Duration
	Index       int
	Payload     any
}

type TaskQueue interface {
	Push(task *Task) error
	Pop(ctx context.Context) (*Task, error)
	UpdateTask(ID string, newTime time.Time) error
	RemoveTask(ID string)
	Length() int
	Close()
}
