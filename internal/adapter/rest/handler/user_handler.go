package handler

import (
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/adapter/rest/middleware"
	"keeplo/internal/adapter/rest/response"
	"keeplo/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
//	@Router			/auth/signup [post]
func (h *Handler) SignupHandler(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleResponse(c, http.StatusBadRequest, 40001, gin.H{
			"error": "입력값이 유효하지 않습니다.",
		})
		return
	}

	user, err := h.UserService.RegisterUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.HandleResponse(c, http.StatusBadRequest, 40002, gin.H{
			"error": err.Error(),
		})
		return
	}

	response.HandleResponse(c, http.StatusOK, 20000, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

// LoginHandler godoc
//
//	@Summary		로그인
//	@Description	로그인 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.LoginRequest	true	"로그인 사용자 정보"
//	@Success		200		{object}	dto.ResponseFormat{data=dto.LoginResponse}
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/auth/login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleResponse(c, 400, 400, nil)
		return
	}

	user, err := h.UserService.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.HandleResponse(c, 400, 400, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		response.HandleResponse(c, 500, 500, nil)
		return
	}

	response.HandleResponse(c, 200, 200, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
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
//	@Router			/auth/me [get]
func (h *Handler) GetUserInfoHandler(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserIDKey).(string)

	user, err := h.UserService.FindByID(c.Request.Context(), userID)
	if err != nil {
		response.HandleResponse(c, 404, 404, gin.H{"error": "user not found"})
		return
	}
	response.HandleResponse(c, 200, 200, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

// UpdateNicknameHandler godoc
//
//	@Summary		닉네임 변경
//	@Description	사용자의 닉네임을 변경합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.UpdateNicknameRequest	true	"닉네임 변경 요청"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/auth/me/nickname [put]
func (h *Handler) UpdateNicknameHandler(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserIDKey).(string)

	var req dto.UpdateNicknameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleResponse(c, 400, 400, gin.H{"error": "잘못된 요청 형식입니다."})
		return
	}

	err := h.UserService.UpdateNickname(c.Request.Context(), userID, req.NickName)
	if err != nil {
		response.HandleResponse(c, 400, 400, gin.H{"error": err.Error()})
		return
	}

	response.HandleResponse(c, 200, 200, gin.H{"message": "닉네임이 변경되었습니다."})
}

// UpdatePasswordHandler godoc
//
//	@Summary		비밀번호 변경
//	@Description	사용자의 비밀번호를 변경합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.UpdatePasswordRequest	true	"비밀번호 변경 요청"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/auth/me/password [put]
func (h *Handler) UpdatePasswordHandler(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserIDKey).(string)

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleResponse(c, 400, 400, gin.H{"error": "잘못된 요청 형식입니다."})
		return
	}

	err := h.UserService.UpdatePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		response.HandleResponse(c, 400, 400, gin.H{"error": err.Error()})
		return
	}

	response.HandleResponse(c, 200, 200, gin.H{"message": "비밀번호가 변경되었습니다."})
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
//	@Router			/auth/me/logout [delete]
func (h *Handler) LogoutHandler(c *gin.Context) {
	// JWT 기반에서는 클라이언트에서 토큰 삭제 또는 무효화
	response.HandleResponse(c, 200, 200, gin.H{
		"message": "로그아웃 되었습니다. 클라이언트에서 토큰을 삭제해주세요.",
	})
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
//	@Router			/auth/me/resign [delete]
func (h *Handler) ReSignHandler(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserIDKey).(string)

	err := h.UserService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		response.HandleResponse(c, 400, 400, gin.H{"error": err.Error()})
		return
	}

	response.HandleResponse(c, 200, 200, gin.H{
		"message": "회원 탈퇴가 완료되었습니다.",
	})
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
//	@Router			/auth/duplicate [get]
func (h *Handler) DuplicateEmail(c *gin.Context) {
	var req dto.DuplicateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleResponse(c, 400, 400, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	exists, err := h.UserService.CheckDuplicateEmail(c.Request.Context(), req.Email)
	if err != nil {
		response.HandleResponse(c, 500, 500, gin.H{"error": err.Error()})
		return
	}

	response.HandleResponse(c, 200, 200, gin.H{"is_duplicate": exists})
}

// CheckPassword godoc
//
//	@Summary		비밀번호 확인 요청
//	@Description	비밀번호 확인 요청
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			password	body	dto.CheckPasswordRequest	true	"비밀번호"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/auth/password [post]
func (h *Handler) CheckPassword(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	var req dto.CheckPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleResponse(c, 400, 400, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	user, err := h.UserService.FindByID(c.Request.Context(), userID)
	if err != nil {
		response.HandleResponse(c, 404, 404, gin.H{"error": "사용자 정보가 없습니다."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.HandleResponse(c, 401, 401, gin.H{"error": "비밀번호가 일치하지 않습니다."})
		return
	}

	response.HandleResponse(c, 200, 200, gin.H{"message": "비밀번호가 확인되었습니다."})
}
