package monitor

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	Name            string    `gorm:"not null"`
	Target          string    `gorm:"not null"`
	Type            string    `gorm:"not null"` // HTTP, WS, TCP
	IntervalSeconds int       `gorm:"not null;default:60"`
	Enabled         bool      `gorm:"default:true"`
	LastCheckedAt   *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
