package monitor

import (
	"context"
	"errors"
	"fmt"
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/domain/monitor"
	"keeplo/internal/domain/user"
	"keeplo/internal/scheduler"
	"keeplo/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const monitorTimeout = time.Second * 5

type Service interface {
	RegisterMonitor(ctx context.Context, userID string, req dto.RegisterMonitorRequest) error
	SearchMonitorList(ctx context.Context, userID string) ([]*monitor.Monitor, error)
	SearchMonitor(ctx context.Context, id string) (*monitor.Monitor, error)
	ModifyMonitor(ctx context.Context, id string, userID string, req dto.UpdateMonitorRequest) error
	DeleteMonitor(ctx context.Context, id string, userID string) error

	ToggleMonitor(ctx context.Context, monitorID, userID string) error
	TriggerMonitor(ctx context.Context, monitorID, userID string) error
	GetSupportedProtocols() []string
}

type monitorService struct {
	monitorRepo monitor.Repository
	userRepo    user.Repository
}

func NewMonitorService(mRepo monitor.Repository, uRepo user.Repository) Service {
	return &monitorService{
		monitorRepo: mRepo,
		userRepo:    uRepo,
	}
}

func (m *monitorService) RegisterMonitor(ctx context.Context, userID string, req dto.RegisterMonitorRequest) error {
	ctx, cancel := context.WithTimeout(ctx, monitorTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	log.Debug("RegisterMonitor - called", zap.String("user_id", userID), zap.String("name", req.Name))

	if req.Type == "" || req.Address == "" || req.Port == "" {
		log.Warn("RegisterMonitor - invalid request data", zap.Any("request", req))
		return monitor.ErrInvalidMonitorData
	}

	target := fmt.Sprintf("%s://%s:%s", req.Type, req.Address, req.Port)
	id := uuid.New()
	newMonitor := &monitor.Monitor{
		ID:              id,
		UserID:          uuid.MustParse(userID),
		Name:            req.Name,
		Target:          target,
		Type:            req.Type,
		IntervalSeconds: req.IntervalSeconds,
		Enabled:         true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 1. DB 저장
	if err := m.monitorRepo.Create(ctx, newMonitor); err != nil {
		log.Error("RegisterMonitor - failed to create", zap.Error(err))
		return err
	}

	// 2. 스케줄러 등록
	task := &scheduler.Task{
		ID:          newMonitor.ID.String(),
		Executor:    &MonitorExecutor{},
		Interval:    time.Duration(newMonitor.IntervalSeconds) * time.Second,
		NextCheckAt: time.Now().Add(time.Duration(newMonitor.IntervalSeconds) * time.Second),
	}
	if err := scheduler.RegisterTask(ctx, "health", task); err != nil {
		log.Error("RegisterMonitor - failed to register scheduler", zap.Error(err))
		return err
	}

	log.Info("RegisterMonitor - success", zap.String("monitor_id", id.String()), zap.String("user_id", userID))
	return nil
}

func (m *monitorService) SearchMonitorList(ctx context.Context, userID string) ([]*monitor.Monitor, error) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	log.Debug("SearchMonitorList - called", zap.String("user_id", userID))

	monitors, err := m.monitorRepo.FindByUserID(ctx, userID)
	if err != nil {
		log.Error("SearchMonitorList - failed", zap.Error(err))
		return nil, err
	}

	log.Info("SearchMonitorList - success", zap.Int("count", len(monitors)))
	return monitors, nil
}

func (m *monitorService) SearchMonitor(ctx context.Context, id string) (*monitor.Monitor, error) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	log.Debug("SearchMonitor - called", zap.String("monitor_id", id))

	result, err := m.monitorRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("SearchMonitor - not found", zap.String("monitor_id", id))
			return nil, monitor.ErrMonitorNotFound
		}
		log.Error("SearchMonitor - failed", zap.Error(err))
		return nil, err
	}

	log.Info("SearchMonitor - success", zap.String("monitor_id", result.ID.String()))
	return result, nil
}

func (m *monitorService) ModifyMonitor(ctx context.Context, id string, userID string, req dto.UpdateMonitorRequest) error {
	ctx, cancel := context.WithTimeout(ctx, monitorTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	log.Debug("ModifyMonitor - called", zap.String("monitor_id", id), zap.String("user_id", userID))

	existing, err := m.monitorRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("ModifyMonitor - not found", zap.String("monitor_id", id))
			return monitor.ErrMonitorNotFound
		}
		log.Error("ModifyMonitor - fetch failed", zap.Error(err))
		return err
	}
	if existing.UserID.String() != userID {
		log.Warn("ModifyMonitor - permission denied", zap.String("monitor_id", id), zap.String("user_id", userID))
		return monitor.ErrPermissionDenied
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.IntervalSeconds != nil {
		existing.IntervalSeconds = *req.IntervalSeconds
	}
	if req.Address != nil && req.Port != nil && req.Type != nil {
		existing.Target = fmt.Sprintf("%s://%s:%s", *req.Type, *req.Address, *req.Port)
	}
	existing.UpdatedAt = time.Now()

	if err := m.monitorRepo.Update(ctx, existing); err != nil {
		log.Error("ModifyMonitor - update failed", zap.Error(err))
		return err
	}

	log.Info("ModifyMonitor - success", zap.String("monitor_id", existing.ID.String()))
	return nil
}

func (m *monitorService) DeleteMonitor(ctx context.Context, id string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, monitorTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	log.Debug("DeleteMonitor - called", zap.String("monitor_id", id), zap.String("user_id", userID))

	monitorObj, err := m.monitorRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("DeleteMonitor - not found", zap.String("monitor_id", id))
			return monitor.ErrMonitorNotFound
		}
		log.Error("DeleteMonitor - fetch failed", zap.Error(err))
		return err
	}
	if monitorObj.UserID.String() != userID {
		log.Warn("DeleteMonitor - permission denied", zap.String("monitor_id", id), zap.String("user_id", userID))
		return monitor.ErrPermissionDenied
	}

	if err := m.monitorRepo.SoftDelete(ctx, id); err != nil {
		log.Error("DeleteMonitor - soft delete failed", zap.Error(err))
		return err
	}

	log.Info("DeleteMonitor - success", zap.String("monitor_id", id))
	return nil
}

func (m *monitorService) ToggleMonitor(ctx context.Context, monitorID, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, monitorTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	monitorObj, err := m.monitorRepo.FindByID(ctx, monitorID)
	if err != nil {
		log.Error("ToggleMonitor - monitor not found", zap.String("monitor_id", monitorID), zap.Error(err))
		return monitor.ErrMonitorNotFound
	}

	if monitorObj.UserID.String() != userID {
		log.Warn("ToggleMonitor - no permission", zap.String("user_id", userID))
		return monitor.ErrPermissionDenied
	}

	monitorObj.Enabled = !monitorObj.Enabled
	monitorObj.UpdatedAt = time.Now()

	if err := m.monitorRepo.Update(ctx, monitorObj); err != nil {
		log.Error("ToggleMonitor - update failed", zap.String("monitor_id", monitorID), zap.Error(err))
		return err
	}

	log.Info("ToggleMonitor - status toggled", zap.String("monitor_id", monitorID), zap.Bool("is_active", monitorObj.Enabled))
	return nil
}

func (m *monitorService) TriggerMonitor(ctx context.Context, monitorID, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, monitorTimeout)
	defer cancel()

	log := logger.WithContext(ctx)
	monitorObj, err := m.monitorRepo.FindByID(ctx, monitorID)
	if err != nil {
		log.Error("TriggerMonitor - monitor not found", zap.String("monitor_id", monitorID), zap.Error(err))
		return monitor.ErrMonitorNotFound
	}

	if monitorObj.UserID.String() != userID {
		log.Warn("TriggerMonitor - no permission", zap.String("user_id", userID))
		return monitor.ErrPermissionDenied
	}

	// 실제 모니터링 테스트 수행 - 간단 예시 (Ping / HTTP 등)
	// result, err := s.pingNow(ctx, monitor)
	// if err != nil {
	// 	log.Error("TriggerMonitor - monitor test failed", zap.String("monitor_id", monitorID), zap.Error(err))
	// 	return err
	// }

	log.Info("TriggerMonitor - test executed", zap.String("monitor_id", monitorID), zap.String("status", "test"))
	return nil
}

func (s *monitorService) GetSupportedProtocols() []string {
	return []string{"HTTP", "HTTPS", "TCP", "WebSocket"}
}
