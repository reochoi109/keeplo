package handler

import (
	"errors"
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/adapter/rest/middleware"
	"keeplo/internal/adapter/rest/response"
	"keeplo/internal/domain/user"
	"keeplo/pkg/auth"
	"keeplo/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignupHandler godoc
//
//	@Summary		회원 가입
//	@Description	신규 사용자를 등록합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.SignupRequest	true	"회원가입 요청 정보"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/auth/signup [post]
func (h *Handler) SignupHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("SignupHandler - invalid request", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	u, err := h.UserService.RegisterUser(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailAlreadyExists):
			response.HandleResponse(c, http.StatusBadRequest, response.ErrorEmailAlreadyExists, nil)
		case errors.Is(err, user.ErrInvalidCredentials):
			response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		default:
			log.Error("SignupHandler - unexpected error", zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		}
		return
	}

	response.HandleResponse(c, http.StatusOK, response.SuccessUserRegistered, dto.NewUserResponse(u.ID.String(), u.Email))
}

// LoginHandler godoc
//
//	@Summary		로그인
//	@Description	이메일과 비밀번호로 로그인을 수행합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.LoginRequest	true	"로그인 요청 정보"
//	@Success		200		{object}	dto.ResponseFormat{data=dto.LoginResponse}
//	@Failure		400		{object}	dto.ResponseFormat
//	@Failure		500		{object}	dto.ResponseFormat
//	@Router			/auth/login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("LoginHandler - invalid request", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	log.Debug("LoginHandler called [Start]", zap.String("email", req.Email))
	defer log.Debug("LoginHandler [End]", zap.String("email", req.Email))

	userObj, err := h.UserService.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			response.HandleResponse(c, http.StatusBadRequest, response.ErrorUserNotFound, nil)
		case errors.Is(err, user.ErrInactiveAccount):
			response.HandleResponse(c, http.StatusUnauthorized, response.ErrorInactiveAccount, nil)
		case errors.Is(err, user.ErrInvalidCredentials):
			response.HandleResponse(c, http.StatusUnauthorized, response.ErrorInvalidCredentials, nil)
		default:
			log.Error("LoginHandler - unexpected error", zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		}
		return
	}

	token, err := auth.GenerateToken(userObj.ID.String())
	if err != nil {
		log.Error("LoginHandler - token generation failed", zap.Error(err))
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		return
	}

	log.Info("Login success", zap.String("user_id", userObj.ID.String()), zap.String("email", req.Email))
	response.HandleResponse(c, http.StatusOK, response.SuccessUserLoggedIn, dto.NewLoginResponse(token, userObj.ID.String(), req.Email))
}

// GetUserInfoHandler godoc
//
//	@Summary		내 정보 조회
//	@Description	현재 로그인한 사용자의 정보를 반환합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dto.ResponseFormat{data=dto.UserResponse}
//	@Failure		401		{object}	dto.ResponseFormat
//	@Failure		404		{object}	dto.ResponseFormat
//	@Router			/auth/me [get]
func (h *Handler) GetUserInfoHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	log.Debug("GetUserInfoHandler called [Start]", zap.String("user_id", userID))
	defer log.Debug("GetUserInfoHandler [End]", zap.String("user_id", userID))

	u, err := h.UserService.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			log.Warn("User not found", zap.String("user_id", userID), zap.Error(err))
			response.HandleResponse(c, http.StatusNotFound, response.ErrorUserNotFound, nil)
		} else {
			log.Error("GetUserInfoHandler - unexpected error", zap.String("user_id", userID), zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		}
		return
	}

	log.Info("User info fetched", zap.String("user_id", u.ID.String()), zap.String("email", u.Email))
	response.HandleResponse(c, http.StatusOK, response.SuccessUserFetched, dto.NewUserResponse(u.ID.String(), u.Email))
}

// UpdateNicknameHandler godoc
//
//	@Summary		닉네임 변경
//	@Description	사용자의 닉네임을 수정합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.UpdateNicknameRequest	true	"변경할 닉네임"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Failure		401		{object}	dto.ResponseFormat
//	@Router			/auth/me/nickname [put]
func (h *Handler) UpdateNicknameHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	var req dto.UpdateNicknameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("UpdateNicknameHandler - invalid request body", zap.Error(err), zap.String("user_id", userID))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	log.Debug("UpdateNicknameHandler called [Start]", zap.String("user_id", userID), zap.String("nickname", req.NickName))
	defer log.Debug("UpdateNicknameHandler [End]", zap.String("user_id", userID))

	err := h.UserService.UpdateNickname(ctx, userID, req.NickName)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNicknameRequired):
			log.Warn("Nickname is required", zap.String("user_id", userID), zap.Error(err))
			response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)

		default:
			log.Error("Nickname update failed", zap.String("user_id", userID), zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		}
		return
	}

	log.Info("Nickname updated", zap.String("user_id", userID), zap.String("new_nickname", req.NickName))
	response.HandleResponse(c, http.StatusOK, response.SuccessUserUpdated, nil)
}

// UpdatePasswordHandler godoc
//
//	@Summary		비밀번호 변경
//	@Description	기존 비밀번호를 새로운 비밀번호로 변경합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.UpdatePasswordRequest	true	"비밀번호 변경 정보"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Failure		401		{object}	dto.ResponseFormat
//	@Router			/auth/me/password [put]
func (h *Handler) UpdatePasswordHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("UpdatePasswordHandler - invalid request body", zap.Error(err), zap.String("user_id", userID))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	log.Debug("UpdatePasswordHandler called [Start]", zap.String("user_id", userID))
	defer log.Debug("UpdatePasswordHandler [End]", zap.String("user_id", userID))

	err := h.UserService.UpdatePassword(ctx, userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrPasswordMismatch):
			log.Warn("Current password mismatch", zap.String("user_id", userID))
			response.HandleResponse(c, http.StatusUnauthorized, response.ErrorPasswordMismatch, nil)

		case errors.Is(err, user.ErrUserNotFound):
			log.Warn("User not found", zap.String("user_id", userID))
			response.HandleResponse(c, http.StatusNotFound, response.ErrorUserNotFound, nil)

		default:
			log.Error("Password update failed", zap.String("user_id", userID), zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		}
		return
	}

	log.Info("Password updated successfully", zap.String("user_id", userID))
	response.HandleResponse(c, http.StatusOK, response.SuccessPasswordChanged, nil)
}

// LogoutHandler godoc
//
//	@Summary		로그아웃
//	@Description	JWT 기반 로그아웃. 클라이언트 토큰 삭제 필요.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dto.ResponseFormat
//	@Router			/auth/me/logout [delete]
func (h *Handler) LogoutHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	log.Debug("LogoutHandler - user logged out", zap.String("user_id", userID))

	response.HandleResponse(c, http.StatusOK, response.SuccessLoggedOut, nil)
}

// ReSignHandler godoc
//
//	@Summary		회원 탈퇴
//	@Description	현재 로그인한 사용자를 탈퇴 처리합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.ResponseFormat
//	@Failure		401	{object}	dto.ResponseFormat
//	@Failure		404	{object}	dto.ResponseFormat
//	@Failure		500	{object}	dto.ResponseFormat
//	@Router			/auth/me/resign [delete]
func (h *Handler) ReSignHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	log.Debug("ReSignHandler called [Start]", zap.String("user_id", userID))
	defer log.Debug("ReSignHandler [End]", zap.String("user_id", userID))

	err := h.UserService.DeleteUser(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			log.Warn("User not found", zap.String("user_id", userID))
			response.HandleResponse(c, http.StatusNotFound, response.ErrorUserNotFound, nil)

		case errors.Is(err, user.ErrAlreadyDeleted):
			log.Warn("User already deleted", zap.String("user_id", userID))
			response.HandleResponse(c, http.StatusBadRequest, response.ErrorBadRequest, nil)

		default:
			log.Error("User deletion failed", zap.String("user_id", userID), zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		}
		return
	}

	log.Info("User resigned successfully", zap.String("user_id", userID))
	response.HandleResponse(c, http.StatusOK, response.SuccessUserResigned, nil)
}

// DuplicateEmail godoc
//
//	@Summary		이메일 중복 확인
//	@Description	입력한 이메일이 이미 사용 중인지 확인합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	body	dto.DuplicateEmailRequest	true	"중복 확인할 이메일"
//	@Success		200		{object}	dto.ResponseFormat{data=bool}
//	@Failure		400		{object}	dto.ResponseFormat
//	@Failure		500		{object}	dto.ResponseFormat
//	@Router			/auth/duplicate [get]
func (h *Handler) DuplicateEmail(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	var req dto.DuplicateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("DuplicateEmail - invalid request body", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	log.Debug("DuplicateEmail called [Start]", zap.String("email", req.Email))
	defer log.Debug("DuplicateEmail [End]", zap.String("email", req.Email))

	exists, err := h.UserService.CheckDuplicateEmail(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmailFormat):
			log.Warn("Invalid email format", zap.String("email", req.Email))
			response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)

		default:
			log.Error("Email check failed", zap.String("email", req.Email), zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorDatabase, nil)
		}
		return
	}

	log.Info("Duplicate email check success", zap.String("email", req.Email), zap.Bool("exists", exists))
	response.HandleResponse(c, http.StatusOK, response.SuccessDuplicateChecked, dto.NewDuplicateEmailResponse(exists))
}

// CheckPassword godoc
//
//	@Summary		비밀번호 확인
//	@Description	입력한 비밀번호가 현재 비밀번호와 일치하는지 확인합니다.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			password	body	dto.CheckPasswordRequest	true	"확인할 비밀번호"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Failure		401		{object}	dto.ResponseFormat
//	@Failure		404		{object}	dto.ResponseFormat
//	@Router			/auth/password [post]
func (h *Handler) CheckPassword(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	var req dto.CheckPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("CheckPassword - invalid request body", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	log.Debug("CheckPassword called [Start]", zap.String("user_id", userID))
	defer log.Debug("CheckPassword [End]", zap.String("user_id", userID))

	err := h.UserService.CheckPassword(ctx, userID, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			log.Warn("User not found", zap.String("user_id", userID))
			response.HandleResponse(c, http.StatusNotFound, response.ErrorUserNotFound, nil)

		case errors.Is(err, user.ErrPasswordMismatch):
			log.Warn("Password mismatch", zap.String("user_id", userID))
			response.HandleResponse(c, http.StatusUnauthorized, response.ErrorPasswordMismatch, nil)

		default:
			log.Error("CheckPassword failed", zap.String("user_id", userID), zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorDatabase, nil)
		}
		return
	}

	log.Info("Password verified", zap.String("user_id", userID))
	response.HandleResponse(c, http.StatusOK, response.SuccessPasswordVerified, nil)
}
