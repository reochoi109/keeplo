package checker

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const defaultTimeout = 5 * time.Second

type CheckResult struct {
	Status     string // "up" or "down"
	Message    string // 실패 이유 (또는 비어있음)
	ResponseMs int    // 응답 시간 (ms)
}

type Checker interface {
	Check(ctx context.Context, target string) (*CheckResult, error)
}

func RunHealthCheck(ctx context.Context, proto string, target string) (*CheckResult, error) {
	var c Checker

	switch strings.ToLower(proto) {
	case "http", "https":
		c = &HTTPChecker{}
	case "tcp":
		c = &TCPChecker{}
	case "ws", "websocket":
		c = &WSChecker{}
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", proto)
	}

	return c.Check(ctx, target)
}
