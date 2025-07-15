package handler

import (
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/adapter/rest/middleware"
	"keeplo/internal/adapter/rest/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 모니터링 등록을 제외해서 다른 사용자가 나를 위장해서 접근하지 않도록 막는것이 관건.

// RegisterMonitorHandler godoc
//
//	@Summary		모니터링 추가
//	@Description	모니터링 추가 요청
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			monitor	body		dto.RegisterMonitorRequest	true	"신규 모니터링"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/monitor [post]
func (h *Handler) RegisterMonitorHandler(c *gin.Context) {
	var req dto.RegisterMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		response.HandleResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized, nil)
		return
	}

	if err := h.MonitorService.RegisterMonitor(c.Request.Context(), userID.(string), req); err != nil {
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorRegisterFailed, nil)
		return
	}
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorRegistered, nil)
}

// GetMonitorListHandler godoc
//
//	@Summary		모니터링 목록
//	@Description	사용자가 등록한 모니터링 목록 조회
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/monitor [get]
func (h *Handler) GetMonitorListHandler(c *gin.Context) {
	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		response.AbortWithResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized)
		return
	}

	monitors, err := h.MonitorService.SearchMonitorList(c.Request.Context(), userID.(string))
	if err != nil {
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
//	@Summary		상세 모니터링 정보
//	@Description	사용자가 등록한 모니터링 상세 정보 조회
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			id 		path 		string 			true 	"모니터링 고유 번호"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/monitor/{id} [get]
func (h *Handler) GetMonitorHandler(c *gin.Context) {
	_, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		response.AbortWithResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized)
		return
	}
	id := c.Param("id")

	m, err := h.MonitorService.SearchMonitor(c.Request.Context(), id)
	if err != nil {
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorFetchFailed, nil)
		return
	}
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorFetched, dto.ToMonitorResponse(m))
}

// UpdateMonitorHandler godoc
//
//	@Summary		모니터링 상세 정보 업데이트
//	@Description	모니터링 상세 정보 업데이트 요청
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			id 		path 		string 				      true "모니터링 고유 번호"
//	@Param			monitor body 		dto.UpdateMonitorRequest  true "모니터링 업데이트 정보"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/monitor/{id} [put]
func (h *Handler) UpdateMonitorHandler(c *gin.Context) {
	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		response.AbortWithResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized)
		return
	}
	id := c.Param("id")
	var req dto.UpdateMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	err := h.MonitorService.ModifyMonitor(c.Request.Context(), id, userID.(string), req)
	if err != nil {
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorUpdateFailed, nil)
		return
	}
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorUpdated, nil)
}

// RemoveMonitorHandler godoc
//
//	@Summary		모니터링 삭제
//	@Description	모니터링 삭제 요청
//	@Tags			monitor
//	@Accept			json
//	@Produce		json
//	@Param			id 		path 		string 				      true "모니터링 고유 번호"
//	@Success		200		{object}	dto.ResponseFormat
//	@Failure		400		{object}	dto.ResponseFormat
//	@Router			/monitor/{id} [delete]
func (h *Handler) RemoveMonitorHandler(c *gin.Context) {
	_, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		response.AbortWithResponse(c, http.StatusUnauthorized, response.ErrorUnauthorized)
		return
	}

	id := c.Param("id")
	if err := h.MonitorService.DeleteMonitor(c.Request.Context(), id); err != nil {
		response.HandleResponse(c, http.StatusInternalServerError, response.ErrorMonitorDeleteFailed, nil)
		return
	}
	response.HandleResponse(c, http.StatusOK, response.SuccessMonitorDeleted, nil)
}
