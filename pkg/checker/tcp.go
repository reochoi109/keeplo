package checker

import (
	"context"
	"fmt"
	"net"
	"time"
)

type TCPChecker struct{}

func (t *TCPChecker) Check(ctx context.Context, target string) (*CheckResult, error) {
	start := time.Now()

	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", target)
	if err != nil {
		return nil, fmt.Errorf("tcp connection failed: %w", err)
	}
	elapsed := time.Since(start).Milliseconds()

	defer conn.Close()
	return &CheckResult{"up", "", int(elapsed)}, nil
}
