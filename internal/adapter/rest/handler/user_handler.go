package handler

import (
	"keeplo/internal/adapter/rest/response"

	"github.com/gin-gonic/gin"
)

// SignupHandler godoc
//
//	@Summary		회원 가입
//	@Description	회원 가입 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.SignupRequest	true	"가입 사용자 정보"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/api/v1/auth/signup [post]
func SignupHandler(c *gin.Context) {
	response.HandleResponse(c, 200, 200, gin.H{"message ": "SignupHandler"})
}

// LoginHandler godoc
//
//	@Summary		로그인
//	@Description	로그인 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.LoginRequest	true	"로그인 사용자 정보"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/api/v1/auth/signup [post]
func LoginHandler(c *gin.Context) {
	response.HandleResponse(c, 200, 200, gin.H{"message ": "LoginHandler"})
}

// GetUserInfoHandler godoc
//
//	@Summary		사용자 정보
//	@Description	사용자 정보 조회
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"사용자 고유 번호"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/api/v1/auth/me/{id} [get]
func GetUserInfoHandler(c *gin.Context) {
	response.HandleResponse(c, 200, 200, gin.H{"message ": "GetUserInfoHandler"})
}

// UpdateUserInfoHandler godoc
//
//	@Summary		사용자 정보 수정
//	@Description	사용자 정보 수정 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		dto.UpdateUserInfoRequest	true	"사용자 수정 정보"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/api/v1/auth/me/{id} [put]
func UpdateUserInfoHandler(c *gin.Context) {
	response.HandleResponse(c, 200, 200, gin.H{"message ": "UpdateUserInfoHandler"})
}

// LogoutHandler godoc
//
//	@Summary		로그아웃
//	@Description	로그아웃 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"사용자 고유번호"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/api/v1/auth/me/{id}/logout [delete]
func LogoutHandler(c *gin.Context) {
	response.HandleResponse(c, 200, 200, gin.H{"message ": "LogoutHandler"})
}

// ReSignHandler godoc
//
//	@Summary		회원탈퇴
//	@Description	회원탈퇴 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"사용자 고유번호"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/api/v1/auth/me/{id}/resign [delete]
func ReSignHandler(c *gin.Context) {
	response.HandleResponse(c, 204, 200, gin.H{"message ": "ReSignHandler"})
}

// DuplicateEmail godoc
//
//	@Summary		이메일 중복 체크
//	@Description	이메일 중복 체크 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	body		dto.DuplicateEmailRequest	true	"이메일"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/api/v1/paste/{id} [get]
func DuplicateEmail(c *gin.Context) {
	response.HandleResponse(c, 200, 200, gin.H{"message ": "DuplicateEmail"})
}

// CheckPassword godoc
//
//	@Summary		비밀번호 확인
//	@Description	비밀번호 확인 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			paste	body		dto.CheckPasswordRequest	true	"Paste update Content"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/api/v1/paste/{id} [post]
func CheckPassword(c *gin.Context) {
	response.HandleResponse(c, 200, 200, gin.H{"message ": "CheckPassword"})
}
