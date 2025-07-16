package handler

import (
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/adapter/rest/middleware"
	"keeplo/internal/adapter/rest/response"
	"keeplo/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 모니터링 등록을 제외해서 다른 사용자가 나를 위장해서 접근하지 않도록 막는것이 관건.

// RegisterMonitorHandler godoc
//
//	@Summary		모니터링 추가
//	@Description	모니터링 항목을 신규 등록합니다.
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			monitor	body		dto.RegisterMonitorRequest	true	"신규 모니터링 등록 요청"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Failure		401		{object}	dto.ResponseFormat
//	@Failure		500		{object}	dto.ResponseFormat
//	@Router			/monitor [post]
func (h *Handler) RegisterMonitorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		log.Error("authorized : not found user info", zap.String("ip", c.Request.Host))
		response.HandleResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized, nil)
		return
	}
	log.Debug("RegisterMonitorHandler called [Start]", zap.String("id", userID.(string)))
	defer log.Debug("RegisterMonitorHandler [End]", zap.String("id", userID.(string)))

	var req dto.RegisterMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("input error : invalid data", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	if err := h.MonitorService.RegisterMonitor(c.Request.Context(), userID.(string), req); err != nil {
		log.Error("internal error", zap.Error(err))
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorRegisterFailed, nil)
		return
	}
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorRegistered, nil)
}

// GetMonitorListHandler godoc
//
//	@Summary		모니터링 목록 조회
//	@Description	사용자가 등록한 모든 모니터링 항목을 조회합니다.
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dto.ResponseFormat{data=[]dto.MonitorResponse}
//	@Failure		401		{object}	dto.ResponseFormat
//	@Failure		500		{object}	dto.ResponseFormat
//	@Router			/monitor [get]
func (h *Handler) GetMonitorListHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		log.Error("unauthorized access", zap.String("ip", c.Request.Host))
		response.AbortWithResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized)
		return
	}

	log.Debug("GetMonitorListHandler called [Start]", zap.String("user_id", userID.(string)))
	defer log.Debug("GetMonitorListHandler [End]", zap.String("user_id", userID.(string)))

	monitors, err := h.MonitorService.SearchMonitorList(ctx, userID.(string))
	if err != nil {
		log.Error("failed to fetch monitor list", zap.Error(err))
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		return
	}

	var list []dto.MonitorResponse
	for _, m := range monitors {
		list = append(list, dto.ToMonitorResponse(m))
	}
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorListed, list)
}

// GetMonitorHandler godoc
//
//	@Summary		모니터링 상세 조회
//	@Description	특정 모니터링 항목의 상세 정보를 조회합니다.
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"모니터링 고유 ID"
//	@Success		200	{object}	dto.ResponseFormat{data=dto.MonitorResponse}
//	@Failure		400	{object}	dto.ResponseFormat
//	@Failure		401	{object}	dto.ResponseFormat
//	@Failure		500	{object}	dto.ResponseFormat
//	@Router			/monitor/{id} [get]
func (h *Handler) GetMonitorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		log.Error("unauthorized access", zap.String("ip", c.Request.Host))
		response.AbortWithResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized)
		return
	}

	id := c.Param("id")
	log.Debug("GetMonitorHandler called [Start]", zap.String("id", id), zap.String("user_id", userID.(string)))
	defer log.Debug("GetMonitorHandler [End]", zap.String("id", id))

	m, err := h.MonitorService.SearchMonitor(ctx, id)
	if err != nil {
		log.Error("failed to fetch monitor detail", zap.String("monitor_id", id), zap.Error(err))
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorFetchFailed, nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorFetched, dto.ToMonitorResponse(m))
}

// UpdateMonitorHandler godoc
//
//	@Summary		모니터링 수정
//	@Description	기존 모니터링 항목의 정보를 수정합니다.
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"모니터링 고유 ID"
//	@Param			monitor	body	dto.UpdateMonitorRequest	true	"수정할 정보"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Failure		401		{object}	dto.ResponseFormat
//	@Failure		500		{object}	dto.ResponseFormat
//	@Router			/monitor/{id} [put]
func (h *Handler) UpdateMonitorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		log.Error("unauthorized access", zap.String("ip", c.Request.Host))
		response.AbortWithResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized)
		return
	}

	id := c.Param("id")
	log.Debug("UpdateMonitorHandler called [Start]", zap.String("id", id), zap.String("user_id", userID.(string)))
	defer log.Debug("UpdateMonitorHandler [End]", zap.String("id", id))

	var req dto.UpdateMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("input error : invalid update data", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	err := h.MonitorService.ModifyMonitor(ctx, id, userID.(string), req)
	if err != nil {
		log.Error("update error", zap.Error(err))
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorUpdateFailed, nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorUpdated, nil)
}

// RemoveMonitorHandler godoc
//
//	@Summary		모니터링 삭제
//	@Description	특정 모니터링 항목을 삭제합니다.
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"모니터링 고유 ID"
//	@Success		200	{object}	dto.ResponseFormat
//	@Failure		401	{object}	dto.ResponseFormat
//	@Failure		500	{object}	dto.ResponseFormat
//	@Router			/monitor/{id} [delete]
func (h *Handler) RemoveMonitorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		log.Error("unauthorized access", zap.String("ip", c.Request.Host))
		response.AbortWithResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized)
		return
	}

	id := c.Param("id")
	log.Debug("RemoveMonitorHandler called [Start]", zap.String("id", id), zap.String("user_id", userID.(string)))
	defer log.Debug("RemoveMonitorHandler [End]", zap.String("id", id))

	if err := h.MonitorService.DeleteMonitor(ctx, id); err != nil {
		log.Error("delete error", zap.Error(err))
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorDeleteFailed, nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorDeleted, nil)
}
