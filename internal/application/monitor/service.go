package monitor

import (
	"context"
	"fmt"
	"keeplo/internal/domain/monitor"
	"keeplo/internal/domain/user"
	"keeplo/internal/scheduler"
	"time"
)

const monitorTimout = time.Second * 5

type MonitorService interface {
	RegisterMonitor(ctx context.Context, userID string)
	SearchMonitorList(ctx context.Context)
	SearchMonitor(ctx context.Context)
	ModifyMonitor(ctx context.Context)
	DeleteMonitor(ctx context.Context)
	CurrentMonitorInfo(ctx context.Context)
}

type monitorService struct {
	monitorRepo monitor.Repository
	userRepo    user.Repository
}

func NewMonitorService(mRepo monitor.Repository, uRepo user.Repository) MonitorService {
	return &monitorService{
		monitorRepo: mRepo,
		userRepo:    uRepo,
	}
}

func (m *monitorService) RegisterMonitor(ctx context.Context, userID string) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	// 사용자 검증

	// 예시: 모니터 엔티티 생성
	monitorID := "monitor-123"
	interval := time.Minute * 1

	// db저장

	// 스케줄러 등록
	if err := scheduler.RegisterTask(ctx, monitorID, interval); err != nil {
		fmt.Println("스케줄 등록 실패:", err)
	}

	// 완료
}

func (m *monitorService) SearchMonitorList(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	// 사용자 정보를 확인

	// 목록 조회

	// 완료
}

func (m *monitorService) SearchMonitor(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	// 사용자 정보 확인

	// 상세 정보 조회

	// 완료
}

func (m *monitorService) ModifyMonitor(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	// 사용자 정보 확인

	// 모니터 정보 확인

	// 수정 완료
}

func (m *monitorService) DeleteMonitor(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()

	// 사용자 정보 확인

	// 모니터 정보 확인

	// 삭제 완료
}

func (m *monitorService) CurrentMonitorInfo(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()
}
