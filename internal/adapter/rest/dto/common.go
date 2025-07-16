package dto

type ResponseFormat struct {
	Message      string `json:"message,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Data         any    `json:"data,omitempty"`
	ErrorCode    int    `json:"error_code,omitempty"`
}
