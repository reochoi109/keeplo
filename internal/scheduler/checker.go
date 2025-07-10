package scheduler

import (
	"context"
	"time"
)

// task

// in memory

// scheduler

// monitoring check function

// register function

type Task struct {
	MonitorID   string
	NextCheckAt time.Time
}

type TaskQueue interface {
	Push(task *Task) error
	Pop(ctx context.Context) (*Task, error)
	Length() int
}
