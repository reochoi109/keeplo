package checker

import (
	"net"
	"time"
)

func RequestTCP(opt Option) Response {
	start := time.Now()
	duration := time.Since(start).Milliseconds()
	conn, err := net.DialTimeout("tcp", opt.URL, time.Duration(opt.TimeoutSeconds)*time.Second)
	if err != nil {
		return Response{"DOWN", 0, duration, err}
	}
	defer conn.Close()
	return Response{"UP", 0, duration, nil}
}
