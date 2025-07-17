package logdata

import (
	"context"
	"time"
)

type Repository interface {
	GetHealthLogs(ctx context.Context, monitorID string, from, to time.Time, statusCode *int, isSuccess *bool) ([]MonitorHealthLog, error)
	// GetErrorSummary(ctx context.Context, monitorID string, from, to time.Time, statusCode *int) (map[int]int, error)
}
