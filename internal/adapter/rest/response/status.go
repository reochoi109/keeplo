package response

type StatusCode int

const (
	//  Success Codes (1xxx)
	Success StatusCode = 1000

	// --- Monitor Success (1100~)
	SuccessMonitorRegistered StatusCode = 1101
	SuccessMonitorListed     StatusCode = 1102
	SuccessMonitorDeleted    StatusCode = 1103
	SuccessMonitorUpdated    StatusCode = 1104
	SuccessMonitorFetched    StatusCode = 1105

	// --- Auth Success (1200~)
	SuccessUserRegistered   StatusCode = 1201
	SuccessUserLoggedIn     StatusCode = 1202
	SuccessUserFetched      StatusCode = 1203
	SuccessUserUpdated      StatusCode = 1204
	SuccessUserResigned     StatusCode = 1205
	SuccessPasswordChanged  StatusCode = 1206
	SuccessDuplicateChecked StatusCode = 1207
	SuccessPasswordVerified StatusCode = 1208
	SuccessLoggedOut        StatusCode = 1209

	//  Client Error Codes (4xxx)
	ErrorBadRequest       StatusCode = 4000
	ErrorValidationFailed StatusCode = 4001

	// --- Monitor Errors (4100~)
	ErrorMonitorNotFound      StatusCode = 4101
	ErrorMonitorLimitExceeded StatusCode = 4102
	ErrorInvalidProtocol      StatusCode = 4103
	ErrorMonitorUpdateFailed  StatusCode = 4104
	ErrorMonitorDeleteFailed  StatusCode = 4105
	ErrorMonitorFetchFailed   StatusCode = 4106

	// --- Auth Errors (4200~)
	ErrorUserNotFound       StatusCode = 4201
	ErrorEmailAlreadyExists StatusCode = 4202
	ErrorPasswordMismatch   StatusCode = 4203

	// Auth & Rate Limit (4400~)
	ErrorUnauthorized      StatusCode = 4400
	ErrorRateLimitExceeded StatusCode = 4403

	// Server Error Codes (5xxx)
	ErrorInternalServer StatusCode = 5000
	ErrorDatabase       StatusCode = 5001

	// --- Monitor Failures (5100~)
	ErrorMonitorRegisterFailed StatusCode = 5101
)

var messageMap = map[StatusCode]string{
	// Success
	Success:                  "요청이 성공적으로 처리되었습니다.",
	SuccessMonitorRegistered: "모니터링 항목이 성공적으로 등록되었습니다.",
	SuccessMonitorListed:     "모니터링 목록 조회 성공.",
	SuccessMonitorDeleted:    "모니터링 항목이 성공적으로 삭제되었습니다.",
	SuccessMonitorUpdated:    "모니터링 항목이 성공적으로 수정되었습니다.",
	SuccessMonitorFetched:    "모니터링 상세 정보 조회 성공.",
	SuccessUserRegistered:    "회원가입이 완료되었습니다.",
	SuccessUserLoggedIn:      "로그인 성공.",
	SuccessUserFetched:       "사용자 정보 조회 성공.",
	SuccessUserUpdated:       "사용자 정보가 성공적으로 수정되었습니다.",
	SuccessUserResigned:      "회원 탈퇴가 완료되었습니다.",
	SuccessPasswordChanged:   "비밀번호가 성공적으로 변경되었습니다.",
	SuccessDuplicateChecked:  "이메일 중복 확인 완료.",
	SuccessPasswordVerified:  "비밀번호가 확인되었습니다.",
	SuccessLoggedOut:         "로그아웃 되었습니다. 클라이언트에서 토큰을 삭제해주세요.",

	// Client Errors
	ErrorBadRequest:           "잘못된 요청입니다.",
	ErrorValidationFailed:     "입력값 유효성 검사에 실패했습니다.",
	ErrorMonitorNotFound:      "해당 모니터링 항목을 찾을 수 없습니다.",
	ErrorMonitorLimitExceeded: "모니터링 등록 가능 수를 초과했습니다.",
	ErrorInvalidProtocol:      "지원하지 않는 프로토콜입니다.",
	ErrorMonitorUpdateFailed:  "모니터링 수정에 실패했습니다.",
	ErrorMonitorDeleteFailed:  "모니터링 삭제에 실패했습니다.",
	ErrorMonitorFetchFailed:   "모니터링 상세 조회에 실패했습니다.",
	ErrorUserNotFound:         "해당 사용자를 찾을 수 없습니다.",
	ErrorEmailAlreadyExists:   "이미 사용 중인 이메일입니다.",
	ErrorPasswordMismatch:     "비밀번호가 일치하지 않습니다.",

	// Auth / Rate Limit
	ErrorUnauthorized:      "인증이 필요합니다.",
	ErrorRateLimitExceeded: "요청이 너무 많습니다. 잠시 후 다시 시도해주세요.",

	// Server Errors
	ErrorInternalServer:        "서버 내부 오류가 발생했습니다.",
	ErrorDatabase:              "데이터베이스 오류가 발생했습니다.",
	ErrorMonitorRegisterFailed: "모니터링 등록 중 오류가 발생했습니다.",
}

func GetMessage(code StatusCode) string {
	if msg, ok := messageMap[code]; ok {
		return msg
	}
	return "정의되지 않은 상태 코드입니다."
}
