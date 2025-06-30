package checker

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func RequestWebSocket(opt Option) Response {
	dialer := websocket.Dialer{
		HandshakeTimeout: time.Duration(opt.TimeoutSeconds) * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: !opt.VerifySSL},
	}

	reqHeader := http.Header{}
	for k, v := range opt.Headers {
		reqHeader.Set(k, v)
	}

	start := time.Now()
	duration := time.Since(start).Milliseconds()
	conn, _, err := dialer.Dial(opt.URL, reqHeader) // header 안들어가?
	if err != nil {
		return Response{"DOWN", 0, duration, err}
	}

	defer conn.Close()
	return Response{"DOWN", 0, duration, err}
}
