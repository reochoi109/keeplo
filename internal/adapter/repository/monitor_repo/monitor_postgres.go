package monitor_repo

import (
	"context"
	"keeplo/internal/domain/monitor"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MonitorGorm struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	Name            string    `gorm:"not null"`
	Target          string    `gorm:"not null"`
	Type            string    `gorm:"not null"`
	IntervalSeconds int       `gorm:"not null;default:60"`
	Enabled         bool      `gorm:"default:true"`
	LastCheckedAt   *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	IsDeleted       bool
}

type GormMonitorRepo struct {
	db *gorm.DB
}

func NewGormMonitorRepo(db *gorm.DB) monitor.Repository {
	return &GormMonitorRepo{db: db}
}

func (r *GormMonitorRepo) Create(ctx context.Context, m *monitor.Monitor) error {
	return r.db.WithContext(ctx).Create(toGorm(m)).Error
}

func (r *GormMonitorRepo) FindByUserID(ctx context.Context, userID string) ([]*monitor.Monitor, error) {
	var results []MonitorGorm
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_deleted = false", userID).
		Find(&results).Error; err != nil {
		return nil, err
	}

	var list []*monitor.Monitor
	for _, g := range results {
		list = append(list, toEntity(&g))
	}
	return list, nil
}

func (r *GormMonitorRepo) FindByID(ctx context.Context, id string) (*monitor.Monitor, error) {
	var g MonitorGorm
	if err := r.db.WithContext(ctx).
		Where("id = ? AND is_deleted = false", id).
		First(&g).Error; err != nil {
		return nil, err
	}
	return toEntity(&g), nil
}

func (r *GormMonitorRepo) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&MonitorGorm{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

func (r *GormMonitorRepo) HardDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&MonitorGorm{}).Error
}

func (r *GormMonitorRepo) Update(ctx context.Context, m *monitor.Monitor) error {
	return r.db.WithContext(ctx).
		Model(&MonitorGorm{}).
		Where("id = ?", m.ID).
		Updates(MonitorGorm{
			Name:            m.Name,
			Type:            m.Type,
			Target:          m.Target,
			IntervalSeconds: m.IntervalSeconds,
			UpdatedAt:       m.UpdatedAt,
		}).Error
}

func (r *GormMonitorRepo) WithTx(ctx context.Context, fn func(txRepo monitor.Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &GormMonitorRepo{db: tx}
		return fn(txRepo)
	})
}

func toEntity(m *MonitorGorm) *monitor.Monitor {
	return &monitor.Monitor{
		ID:              m.ID,
		UserID:          m.UserID,
		Name:            m.Name,
		Target:          m.Target,
		Type:            m.Type,
		IntervalSeconds: m.IntervalSeconds,
		Enabled:         m.Enabled,
		LastCheckedAt:   m.LastCheckedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func toGorm(m *monitor.Monitor) *MonitorGorm {
	return &MonitorGorm{
		ID:              m.ID,
		UserID:          m.UserID,
		Name:            m.Name,
		Target:          m.Target,
		Type:            m.Type,
		IntervalSeconds: m.IntervalSeconds,
		Enabled:         m.Enabled,
		LastCheckedAt:   m.LastCheckedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		IsDeleted:       false,
	}
}
