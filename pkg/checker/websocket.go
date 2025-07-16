package checker

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type WSChecker struct{}

func (w *WSChecker) Check(ctx context.Context, target string) (*CheckResult, error) {
	start := time.Now()

	dialer := websocket.Dialer{
		HandshakeTimeout: defaultTimeout,
	}

	done := make(chan error, 1)
	go func() {
		conn, _, err := dialer.Dial(target, nil)
		if err != nil {
			done <- fmt.Errorf("websocket dial failed: %w", err)
			return
		}
		conn.Close()
		done <- nil
	}()

	select {
	case <-ctx.Done():
		elapsed := time.Since(start).Milliseconds()
		return &CheckResult{
			Status:     "down",
			Message:    "websocket check timeout",
			ResponseMs: int(elapsed),
		}, ctx.Err()

	case err := <-done:
		elapsed := time.Since(start).Milliseconds()

		if err != nil {
			return &CheckResult{
				Status:     "down",
				Message:    err.Error(),
				ResponseMs: int(elapsed),
			}, err
		}

		return &CheckResult{
			Status:     "up",
			Message:    "",
			ResponseMs: int(elapsed),
		}, nil
	}
}
