package checker

type Option struct {
	Name           string
	Type           string // "http", "websocket", "tcp"
	URL            string
	Method         string
	AuthType       string
	AuthValue      string
	Body           string
	Headers        map[string]string
	TimeoutSeconds int
	ExpectStatuses []int
	VerifySSL      bool
}



type Response struct {
	Status         string // "UP" or "DOWN"
	ResponseCode   int
	ResponseTimeMs int64
	Error          error
}
