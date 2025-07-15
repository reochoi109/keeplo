package monitor

import (
	"context"
	"fmt"
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/domain/monitor"
	"keeplo/internal/domain/user"
	"keeplo/internal/scheduler"
	"time"

	"github.com/google/uuid"
)

const monitorTimout = time.Second * 5

type Service interface {
	RegisterMonitor(ctx context.Context, userID string, req dto.RegisterMonitorRequest) error
	SearchMonitorList(ctx context.Context, userID string) ([]*monitor.Monitor, error)
	SearchMonitor(ctx context.Context, id string) (*monitor.Monitor, error)
	ModifyMonitor(ctx context.Context, id string, userID string, req dto.UpdateMonitorRequest) error
	DeleteMonitor(ctx context.Context, id string) error
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
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()
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

	if err := m.monitorRepo.Create(ctx, newMonitor); err != nil {
		return err
	}

	// 스케줄러 등록
	interval := time.Duration(req.IntervalSeconds) * time.Second
	return scheduler.RegisterTask(ctx, newMonitor.ID.String(), interval)
}

func (m *monitorService) SearchMonitorList(ctx context.Context, userID string) ([]*monitor.Monitor, error) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	return m.monitorRepo.FindByUserID(ctx, userID)
}

func (m *monitorService) SearchMonitor(ctx context.Context, id string) (*monitor.Monitor, error) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	return m.monitorRepo.FindByID(ctx, id)
}

func (m *monitorService) ModifyMonitor(ctx context.Context, id string, userID string, req dto.UpdateMonitorRequest) error {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	existing, err := m.monitorRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("모니터 조회 실패: %w", err)
	}
	if existing.UserID.String() != userID {
		return fmt.Errorf("수정 권한이 없습니다")
	}

	// 수정 필드만 반영
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
		return fmt.Errorf("업데이트 실패: %w", err)
	}

	// 스케줄러도 갱신 (필요 시만)
	// if req.IntervalSeconds != nil {
	// 	interval := time.Duration(*req.IntervalSeconds) * time.Second
	// 	if err := scheduler.UpdateTask(ctx, existing.ID.String(), interval); err != nil {
	// 		return fmt.Errorf("스케줄 갱신 실패: %w", err)
	// 	}
	// }

	return nil
}

func (m *monitorService) DeleteMonitor(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	return m.monitorRepo.SoftDelete(ctx, id)
}
