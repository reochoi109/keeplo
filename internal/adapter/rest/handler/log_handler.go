package handler

import (
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/adapter/rest/response"
	"keeplo/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetMonitorHealthLogHandler godoc
//
//	@Summary		모니터 헬스 로그 조회
//	@Description	특정 모니터의 헬스 체크 이력을 조회합니다.
//	@Tags			log
//	@Produce		json
//	@Param			monitor_id	path	string	true	"모니터 ID"
//	@Param			limit		query	int		false	"페이지당 항목 수 (기본 50)"
//	@Param			offset		query	int		false	"시작 위치 오프셋 (기본 0)"
//	@Success		200			{object}	dto.ResponseFormat
//	@Failure		400			{object}	dto.ResponseFormat
//	@Failure		401			{object}	dto.ResponseFormat
//	@Failure		500			{object}	dto.ResponseFormat
//	@Router			/log/health/{monitor_id} [get]
func (h *Handler) GetMonitorHealthLogHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	monitorID := c.Param("monitor_id")
	if monitorID == "" {
		log.Warn("GetMonitorHealthLogHandler - missing monitor ID")
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorBadRequest, nil)
		return
	}

	var req dto.MonitorHealthLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Warn("GetMonitorHealthLogHandler - invalid query", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	// 기본값 설정
	if req.Limit == 0 {
		req.Limit = 50
	}

	log.Debug("GetMonitorHealthLogHandler",
		zap.String("monitor_id", monitorID),
		zap.Int("limit", req.Limit),
		zap.Int("offset", req.Offset),
	)

	// 목업 데이터
	mock := []dto.MonitorHealthLogResponse{}
	for i := req.Offset; i < req.Offset+req.Limit; i++ {
		mock = append(mock, dto.MonitorHealthLogResponse{
			Timestamp:      time.Now().Add(-time.Duration(i) * time.Minute),
			StatusCode:     200,
			ResponseTimeMs: 120 + i,
			IsSuccess:      true,
			Message:        "OK",
		})
	}

	result := gin.H{
		"logs":   mock,
		"limit":  req.Limit,
		"offset": req.Offset,
		"total":  999, // 예시 값
	}
	response.HandleResponse(c, http.StatusOK, response.Success, result)
}

// GetHealthErrorSummaryHandler godoc
//
//	@Summary		헬스 체크 실패 요약
//	@Description	특정 모니터의 실패 응답 요약을 조회합니다.
//	@Tags			log
//	@Produce		json
//	@Param			monitor_id		path	string	true	"모니터 ID"
//	@Param			from			query	string	false	"조회 시작 시각 (RFC3339 형식, 예: 2025-07-17T00:00:00Z)"
//	@Param			to				query	string	false	"조회 종료 시각 (RFC3339 형식, 예: 2025-07-17T23:59:59Z)"
//	@Param			status_code	query	int		false	"특정 상태코드 필터 (예: 500)"
//	@Success		200	{object}	dto.ResponseFormat
//	@Failure		400	{object}	dto.ResponseFormat
//	@Failure		401	{object}	dto.ResponseFormat
//	@Failure		500	{object}	dto.ResponseFormat
//	@Router			/log/health/{monitor_id}/errors [get]
func (h *Handler) GetHealthErrorSummaryHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	monitorID := c.Param("monitor_id")
	log.Debug("GetHealthErrorSummaryHandler", zap.String("monitor_id", monitorID))

	var req dto.ErrorSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Warn("GetHealthErrorSummaryHandler - invalid query", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	now := time.Now()
	if req.To == nil {
		req.To = &now
	}
	if req.From == nil {
		from := now.Add(-24 * time.Hour)
		req.From = &from
	}

	log.Debug("ErrorSummaryRequest parsed",
		zap.Time("from", *req.From),
		zap.Time("to", *req.To),
		zap.Any("status_code", req.StatusCode),
	)

	// 전체 에러 중 일부 필터링 가능
	all := []dto.ErrorCountResponse{
		{StatusCode: 500, Count: 7},
		{StatusCode: 504, Count: 4},
		{StatusCode: 408, Count: 3},
	}

	filtered := all
	if req.StatusCode != nil {
		filtered = []dto.ErrorCountResponse{}
		for _, e := range all {
			if e.StatusCode == *req.StatusCode {
				filtered = append(filtered, e)
			}
		}
	}

	// 실패 건수 재계산
	total := 0
	for _, e := range filtered {
		total += e.Count
	}

	summary := dto.MonitorErrorSummaryResponse{
		TotalFailures:   total,
		LastFailureTime: now.Add(-10 * time.Minute),
		ErrorCounts:     filtered,
	}
	response.HandleResponse(c, http.StatusOK, response.Success, summary)
}

// GetResponseTimeChartHandler godoc
//
//	@Summary		응답 시간 시계열 데이터
//	@Description	특정 모니터의 응답 시간 차트를 조회합니다.
//	@Tags			log
//	@Produce		json
//	@Param			monitor_id		path	string	true	"모니터 ID"
//	@Param			from			query	string	false	"조회 시작 시각 (RFC3339 형식, 예: 2025-07-17T00:00:00Z)"
//	@Param			to				query	string	false	"조회 종료 시각 (RFC3339 형식, 예: 2025-07-17T23:59:59Z)"
//	@Param			interval		query	string	false	"집계 간격 (예: 1m, 5m, 1h)"
//	@Param			status_code	query	int		false	"특정 응답 코드 (예: 200, 500)"
//	@Param			is_success		query	boolean	false	"성공 여부 필터 (true/false)"
//	@Success		200	{object}	dto.ResponseFormat
//	@Failure		400	{object}	dto.ResponseFormat
//	@Failure		401	{object}	dto.ResponseFormat
//	@Failure		500	{object}	dto.ResponseFormat
//	@Router			/log/health/{monitor_id}/timeseries [get]
func (h *Handler) GetResponseTimeChartHandler(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.WithContext(ctx)

	monitorID := c.Param("monitor_id")
	if monitorID == "" {
		log.Warn("GetResponseTimeChartHandler - missing monitor ID")
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorBadRequest, nil)
		return
	}

	var req dto.TimeSeriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Warn("GetResponseTimeChartHandler - invalid query", zap.Error(err))
		response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
		return
	}

	now := time.Now()
	if req.To == nil {
		req.To = &now
	}
	if req.From == nil {
		from := now.Add(-1 * time.Hour)
		req.From = &from
	}
	interval := 5 * time.Minute
	if req.Interval != "" {
		dur, err := dto.ParseInterval(req.Interval)
		if err != nil {
			log.Warn("GetResponseTimeChartHandler - invalid interval", zap.Error(err))
			response.HandleResponse(c, http.StatusBadRequest, response.ErrorValidationFailed, nil)
			return
		}
		interval = dur
	}

	log.Debug("Parsed TimeSeriesRequest",
		zap.String("monitor_id", monitorID),
		zap.Time("from", *req.From),
		zap.Time("to", *req.To),
		zap.Duration("interval", interval),
		zap.Any("status_code", req.StatusCode),
		zap.Any("is_success", req.IsSuccess),
	)

	var points []dto.TimeSeriesPointResponse
	for t := *req.From; t.Before(*req.To); t = t.Add(interval) {
		points = append(points, dto.TimeSeriesPointResponse{
			Time:            t,
			AvgResponseTime: 120 + float64(t.Unix()%30),
		})
	}

	response.HandleResponse(c, http.StatusOK, response.Success, points)
}

//
// ------------------------------
// [4] 알림 이력
// ------------------------------

// type NotificationLog struct {
// 	Timestamp time.Time `json:"timestamp"`
// 	Type      string    `json:"type"`   // email, slack
// 	Status    string    `json:"status"` // sent, failed
// 	Message   string    `json:"message"`
// }

// // GetNotificationLogsHandler godoc
// //
// //	@Summary		알림 전송 이력
// //	@Description	특정 모니터의 알림 발송 이력을 조회합니다.
// //	@Tags			log
// //	@Produce		json
// //	@Param			monitor_id	path	string	true	"모니터 ID"
// //	@Success		200		{object}	dto.ResponseFormat
// //	@Failure		400		{object}	dto.ResponseFormat
// //	@Failure		401		{object}	dto.ResponseFormat
// //	@Failure		500		{object}	dto.ResponseFormat
// //	@Router			/log/notifications/{monitor_id} [get]
// func (h *Handler) GetNotificationLogsHandler(c *gin.Context) {
// 	monitorID := c.Param("monitor_id")
// 	logger.WithContext(c.Request.Context()).Debug("GetNotificationLogsHandler", zap.String("monitor_id", monitorID))

// 	var logs []NotificationLog
// 	for i := 0; i < 5; i++ {
// 		logs = append(logs, NotificationLog{
// 			Timestamp: time.Now().Add(-time.Duration(i*15) * time.Minute),
// 			Type:      "email",
// 			Status:    "sent",
// 			Message:   "서버가 응답하지 않습니다",
// 		})
// 	}
// 	response.HandleResponse(c, http.StatusOK, response.Success, logs)
// }
