package checker

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func RequestHTTP(opt Option) Response {
	client := &http.Client{
		Timeout: time.Duration(opt.TimeoutSeconds) * time.Second,
	}

	req, err := http.NewRequest(opt.Method, opt.URL, strings.NewReader(opt.Body))
	if err != nil {
		return Response{"DOWN", 0, 0, err}
	}

	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	switch opt.AuthType {
	case "basic":
		parts := strings.SplitN(opt.AuthValue, ":", 2)
		if len(parts) == 2 {
			req.SetBasicAuth(parts[0], parts[1])
		}
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+opt.AuthValue)
	case "apikey":
		req.Header.Set("X-API-Key", opt.AuthValue)
	}

	if !opt.VerifySSL {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// 요청 실행
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return Response{"DOWN", 0, duration, err}
	}
	defer resp.Body.Close()

	status := "DOWN"
	for _, expect := range opt.ExpectStatuses {
		if resp.StatusCode == expect {
			status = "UP"
			break
		}
	}

	return Response{status, resp.StatusCode, duration, nil}
}

func Run() {
	opt := Option{
		Name:   "Test Site",
		URL:    "https://example.com/api/ping",
		Method: "GET",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		TimeoutSeconds: 5,
		ExpectStatuses: []int{200, 201},
		AuthType:       "bearer",
		AuthValue:      "your-token",
		VerifySSL:      true,
	}

	result := RequestHTTP(opt)

	fmt.Printf("[%s] %s → %s (%dms) Code: %d\n",
		opt.Name, opt.URL, result.Status, result.ResponseTimeMs, result.ResponseCode)
}
