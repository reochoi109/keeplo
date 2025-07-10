package monitor

import (
	"context"
	"keeplo/internal/domain/monitor"
	"time"
)

const monitorTimout = time.Second * 5

type MonitorService interface{}

type monitorService struct {
	repo monitor.Repository
}

func (m *monitorService) RegisterMonitor(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, monitorTimout)
	defer cancel()
	// 스케줄러 등록

	// db저장

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
