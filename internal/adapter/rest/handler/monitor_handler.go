package handler

import (
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/adapter/rest/middleware"
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
		c.JSON(http.StatusBadRequest, gin.H{"요청형식 오류": err})
		return
	}

	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"인증 실패": nil})
		return
	}

	if err := h.MonitorService.RegisterMonitor(c.Request.Context(), userID.(string), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"모니터 등록 실패": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"등록 완료": nil})
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
		c.JSON(http.StatusUnauthorized, gin.H{"인증 실패": nil})
		return
	}
	monitors, err := h.MonitorService.SearchMonitorList(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"모니터 조회 실패": err})
		return
	}
	var list []dto.MonitorResponse
	for _, m := range monitors {
		list = append(list, dto.ToMonitorResponse(m))
	}
	c.JSON(http.StatusOK, list)
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
		c.JSON(http.StatusUnauthorized, gin.H{"인증 실패": nil})
		return
	}
	id := c.Param("id")

	m, err := h.MonitorService.SearchMonitor(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"monitor search fail": err})
		return
	}
	res := dto.ToMonitorResponse(m)
	c.JSON(200, gin.H{"message": res})
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
		c.JSON(http.StatusUnauthorized, gin.H{"인증 실패 ": nil})
		return
	}
	id := c.Param("id")

	var req dto.UpdateMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	err := h.MonitorService.ModifyMonitor(c.Request.Context(), id, userID.(string), req)
	if err != nil {
		return
	}

	c.JSON(200, gin.H{"message": "UpdateMonitorHandler"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"인증실패": nil})
		return
	}
	id := c.Param("id")

	if err := h.MonitorService.DeleteMonitor(c.Request.Context(), id); err != nil {
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
