package checker

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type HTTPChecker struct{}

func (h *HTTPChecker) Check(ctx context.Context, target string) (*CheckResult, error) {
	start := time.Now()

	client := http.Client{Timeout: defaultTimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return &CheckResult{"down", "invalid request", 0}, err
	}

	resp, err := client.Do(req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return &CheckResult{"down", err.Error(), int(elapsed)}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("HTTP status %d", resp.StatusCode)
		return &CheckResult{"down", msg, int(elapsed)}, errors.New(msg)
	}
	return &CheckResult{"up", "", int(elapsed)}, nil
}
