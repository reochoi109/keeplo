package logdata

import "time"

type MonitorHealthLog struct {
	ID           string    `bson:"_id"`
	MonitorID    string    `bson:"monitor_id"`
	Timestamp    time.Time `bson:"timestamp"`
	StatusCode   int       `bson:"status_code"`
	ResponseTime float64   `bson:"response_time"`
	IsSuccess    bool      `bson:"is_success"`
}
