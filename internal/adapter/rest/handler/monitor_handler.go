package handler

import (
	"errors"
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/adapter/rest/middleware"
	"keeplo/internal/adapter/rest/response"
	"keeplo/internal/domain/monitor"
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
		log.Warn("RegisterMonitorHandler - unauthorized", zap.String("ip", c.ClientIP()))
		response.HandleResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized, nil)
		return
	}

	var req dto.RegisterMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("RegisterMonitorHandler - invalid request", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	err := h.MonitorService.RegisterMonitor(ctx, userID.(string), req)
	if err != nil {
		switch {
		case errors.Is(err, monitor.ErrInvalidMonitorData):
			response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		default:
			log.Error("RegisterMonitorHandler - internal error", zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorRegisterFailed, nil)
		}
		return
	}

	log.Info("RegisterMonitorHandler - success", zap.String("user_id", userID.(string)))
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
		log.Warn("GetMonitorListHandler - unauthorized", zap.String("ip", c.ClientIP()))
		response.HandleResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized, nil)
		return
	}

	monitors, err := h.MonitorService.SearchMonitorList(ctx, userID.(string))
	if err != nil {
		log.Error("GetMonitorListHandler - fetch failed", zap.Error(err))
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorInternalServer, nil)
		return
	}

	var list []dto.MonitorResponse
	for _, m := range monitors {
		list = append(list, dto.ToMonitorResponse(m))
	}

	log.Info("GetMonitorListHandler - success", zap.String("user_id", userID.(string)))
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
		log.Warn("GetMonitorHandler - unauthorized", zap.String("ip", c.ClientIP()))
		response.HandleResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized, nil)
		return
	}

	id := c.Param("id")
	monitorObj, err := h.MonitorService.SearchMonitor(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, monitor.ErrMonitorNotFound):
			response.HandleResponse(c, http.StatusNotFound, response.ErrorMonitorNotFound, nil)
		default:
			log.Error("GetMonitorHandler - fetch failed", zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorFetchFailed, nil)
		}
		return
	}

	log.Info("GetMonitorHandler - success", zap.String("monitor_id", id), zap.String("user_id", userID.(string)))
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorFetched, dto.ToMonitorResponse(monitorObj))
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
		log.Warn("UpdateMonitorHandler - unauthorized", zap.String("ip", c.ClientIP()))
		response.HandleResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized, nil)
		return
	}

	id := c.Param("id")
	var req dto.UpdateMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("UpdateMonitorHandler - invalid request", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	err := h.MonitorService.ModifyMonitor(ctx, id, userID.(string), req)
	if err != nil {
		switch {
		case errors.Is(err, monitor.ErrMonitorNotFound):
			response.HandleResponse(c, http.StatusNotFound, response.ErrorMonitorNotFound, nil)
		case errors.Is(err, monitor.ErrPermissionDenied):
			response.HandleResponse(c, http.StatusForbidden, response.ErrorPermissionDenied, nil)
		default:
			log.Error("UpdateMonitorHandler - update failed", zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorUpdateFailed, nil)
		}
		return
	}

	log.Info("UpdateMonitorHandler - success", zap.String("monitor_id", id), zap.String("user_id", userID.(string)))
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
		log.Warn("RemoveMonitorHandler - unauthorized", zap.String("ip", c.ClientIP()))
		response.HandleResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized, nil)
		return
	}

	id := c.Param("id")
	err := h.MonitorService.DeleteMonitor(ctx, id, userID.(string))
	if err != nil {
		switch {
		case errors.Is(err, monitor.ErrMonitorNotFound):
			response.HandleResponse(c, http.StatusNotFound, response.ErrorMonitorNotFound, nil)
		case errors.Is(err, monitor.ErrPermissionDenied):
			response.HandleResponse(c, http.StatusForbidden, response.ErrorPermissionDenied, nil)
		default:
			log.Error("RemoveMonitorHandler - delete failed", zap.Error(err))
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorDeleteFailed, nil)
		}
		return
	}

	log.Info("RemoveMonitorHandler - success", zap.String("monitor_id", id), zap.String("user_id", userID.(string)))
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorDeleted, nil)
}

// ToggleMonitorHandler godoc
//
//	@Summary		모니터링 ON/OFF 전환
//	@Description	모니터링 항목을 활성화 또는 비활성화합니다.
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"모니터 ID"
//	@Success		200	{object}	dto.ResponseFormat
//	@Failure		404	{object}	dto.ResponseFormat
//	@Router			/monitor/{id}/toggle [patch]
func (h *Handler) ToggleMonitorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	monitorID := c.Param("id")

	log.Debug("ToggleMonitorHandler called", zap.String("monitor_id", monitorID))

	if err := h.MonitorService.ToggleMonitor(ctx, monitorID, userID); err != nil {
		switch err {
		case monitor.ErrMonitorNotFound:
			response.HandleResponse(c, http.StatusNotFound, response.ErrorMonitorNotFound, nil)
		case monitor.ErrPermissionDenied:
			response.HandleResponse(c, http.StatusForbidden, response.ErrorUnauthorized, nil)
		default:
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorUpdateFailed, nil)
		}
		return
	}
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorUpdated, nil)
}

// TriggerMonitorHandler godoc
//
//	@Summary		모니터링 수동 실행
//	@Description	선택한 모니터링 항목을 즉시 테스트 실행합니다.
//	@Tags			monitor
//	@Produce		json
//	@Param			id	path	string	true	"모니터 ID"
//	@Success		200	{object}	dto.ResponseFormat
//	@Failure		404	{object}	dto.ResponseFormat
//	@Failure		403	{object}	dto.ResponseFormat
//	@Failure		500	{object}	dto.ResponseFormat
//	@Router			/monitor/{id}/trigger [post]
func (h *Handler) TriggerMonitorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	userID := c.MustGet(middleware.ContextUserIDKey).(string)
	monitorID := c.Param("id")
	log.Debug("TriggerMonitorHandler called", zap.String("monitor_id", monitorID), zap.String("user_id", userID))

	if err := h.MonitorService.TriggerMonitor(ctx, monitorID, userID); err != nil {
		switch err {
		case monitor.ErrMonitorNotFound:
			response.HandleResponse(c, http.StatusNotFound, response.ErrorMonitorNotFound, nil)
		case monitor.ErrPermissionDenied:
			response.HandleResponse(c, http.StatusForbidden, response.ErrorUnauthorized, nil)
		default:
			response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorRegisterFailed, nil)
		}
		return
	}
	response.HandleResponse(c, http.StatusOK, response.Success, nil)
}

// GetSupportedProtocolsHandler godoc
//
//	@Summary		지원 프로토콜 조회
//	@Description	서버에서 지원하는 모니터링 프로토콜 목록을 반환합니다.
//	@Tags			monitor
//	@Produce		json
//	@Success		200	{object}	dto.ResponseFormat
//	@Router			/monitor/protocols [get]
func (h *Handler) GetSupportedProtocolsHandler(c *gin.Context) {
	protocols := h.MonitorService.GetSupportedProtocols()
	response.HandleResponse(c, http.StatusOK, response.Success, protocols)
}
