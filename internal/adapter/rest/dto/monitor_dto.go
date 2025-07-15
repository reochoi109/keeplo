package dto

import "keeplo/internal/domain/monitor"

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

// Response --------------------------------------

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
