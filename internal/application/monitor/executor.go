package monitor

import (
	"context"
	"errors"
	"fmt"
	"keeplo/internal/domain/monitor"
	"keeplo/pkg/checker"
	"keeplo/pkg/logger"

	"go.uber.org/zap"
)

type MonitorExecutor struct{}

func (e *MonitorExecutor) Execute(ctx context.Context, payload any) error {
	log := logger.WithContext(ctx)

	m, ok := payload.(*monitor.Monitor)
	if !ok {
		return errors.New("invalid payload type: expected *Monitor")
	}

	var c checker.Checker
	switch m.Type {
	case "HTTP", "HTTPS":
		c = &checker.HTTPChecker{}
	case "TCP":
		c = &checker.TCPChecker{}
	case "WebSocket":
		c = &checker.WSChecker{}
	default:
		return fmt.Errorf("unsupported protocol: %s", m.Type)
	}

	result, err := c.Check(ctx, m.Target)
	if err != nil {
		log.Warn("MonitorExecutor - check failed", zap.String("monitor_id", m.ID.String()), zap.Error(err))
	} else {
		log.Info("MonitorExecutor - check passed",
			zap.String("monitor_id", m.ID.String()),
			zap.String("status", result.Status),
			zap.Int("ms", result.ResponseMs),
		)
	}

	// TODO: 로깅, 알림 전송 등
	return err
}
