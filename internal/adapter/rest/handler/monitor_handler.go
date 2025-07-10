package handler

import "github.com/gin-gonic/gin"

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
//	@Router			/api/v1/monitor [post]
func RegisterMonitorHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "RegisterMonitorHandler"})
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
//	@Router			/api/v1/monitor [get]
func GetMonitorListHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GetMonitorListHandler"})
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
//	@Router			/api/v1/monitor/{id} [get]
func GetMonitorHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GetMonitorHandler"})
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
//	@Router			/api/v1/monitor/{id} [put]
func UpdateMonitorHandler(c *gin.Context) {
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
//	@Router			/api/v1/monitor/{id} [delete]
func RemoveMonitorHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "RemoveMonitorHandler"})
}
