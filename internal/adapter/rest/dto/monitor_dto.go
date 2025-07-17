package dto

import (
	"keeplo/internal/domain/monitor"
	"time"
)

// Request --------------------------------------

type RegisterMonitorRequest struct {
	Name            string `json:"name" binding:"required"`
	Address         string `json:"address" binding:"required"` // 도메인 or IP
	Port            string `json:"port" binding:"required"`    // 포트 번호
	Type            string `json:"type" binding:"required,oneof=http https websocket tcp"`
	IntervalSeconds int    `json:"interval_seconds" binding:"required,min=10"`
}

type UpdateMonitorRequest struct {
	Name            *string `json:"name,omitempty"`
	Address         *string `json:"address,omitempty"`
	Port            *string `json:"port,omitempty"`
	Type            *string `json:"type,omitempty"`
	IntervalSeconds *int    `json:"interval_seconds,omitempty"`
}

type TimeSeriesRequest struct {
	From       *time.Time `form:"from" time_format:"2006-01-02T15:04:05Z07:00"`
	To         *time.Time `form:"to" time_format:"2006-01-02T15:04:05Z07:00"`
	Interval   string     `form:"interval"`
	StatusCode *int       `form:"status_code"`
	IsSuccess  *bool      `form:"is_success"`
}

type ErrorSummaryRequest struct {
	From       *time.Time `form:"from" time_format:"2006-01-02T15:04:05Z07:00"`
	To         *time.Time `form:"to" time_format:"2006-01-02T15:04:05Z07:00"`
	StatusCode *int       `form:"status_code"` // nullable
}

type MonitorHealthLogRequest struct {
	Limit  int `form:"limit" json:"limit" binding:"omitempty,min=1"`
	Offset int `form:"offset" json:"offset" binding:"omitempty,min=0"`
}

// Response --------------------------------------

type TimeSeriesPointResponse struct {
	Time            time.Time `json:"time"`
	AvgResponseTime float64   `json:"avg_response_time"`
}
type MonitorErrorSummaryResponse struct {
	TotalFailures   int                  `json:"total_failures"`
	LastFailureTime time.Time            `json:"last_failure_time"`
	ErrorCounts     []ErrorCountResponse `json:"error_counts"`
}

type ErrorCountResponse struct {
	StatusCode int `json:"status_code"`
	Count      int `json:"count"`
}

type MonitorHealthLogResponse struct {
	Timestamp      time.Time `json:"timestamp"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int       `json:"response_time_ms"`
	IsSuccess      bool      `json:"is_success"`
	Message        string    `json:"message"`
}

type MonitorResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Target          string `json:"target"`
	Type            string `json:"type"`
	IntervalSeconds int    `json:"interval_seconds"`
	Enabled         bool   `json:"enabled"`
	LastCheckedAt   string `json:"last_checked_at,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func ToMonitorResponse(m *monitor.Monitor) MonitorResponse {
	return MonitorResponse{
		ID:              m.ID.String(),
		Name:            m.Name,
		Target:          m.Target,
		Type:            m.Type,
		IntervalSeconds: m.IntervalSeconds,
		Enabled:         m.Enabled,
		LastCheckedAt:   m.LastCheckedAt.String(),
		CreatedAt:       m.CreatedAt.String(),
		UpdatedAt:       m.UpdatedAt.String(),
	}
}
