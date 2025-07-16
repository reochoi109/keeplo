package monitor

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Name            string
	Target          string
	Type            string
	IntervalSeconds int
	Enabled         bool
	LastCheckedAt   *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type HealthLog struct {
	ID         string    // 로그 ID (Mongo ObjectID 문자열 또는 UUID 등)
	MonitorID  string    // string으로 처리하여 양쪽 DB에서 모두 사용 가능
	Status     string    // "up" | "down"
	Message    string    // 실패 시 메시지
	ResponseMs int       // 응답 시간 (ms)
	Timestamp  time.Time // 체크된 시각
}
