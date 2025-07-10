package monitor

import "time"

type URLMonitor struct {
	ID          string
	UserID      string
	Name        string
	Protocol    string // 예: "http"
	Address     string
	Interval    int // 초 단위
	LastChecked time.Time
}

type CheckHistory struct {
	MonitorID   string
	Status      string
	ErrorDetail string
	Timestamp   time.Time
	ResponseMs  int64
}
