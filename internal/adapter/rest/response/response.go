package response

import (
	"keeplo/internal/adapter/rest/dto"

	"github.com/gin-gonic/gin"
)

func HandleResponse(c *gin.Context, statusCode int, code StatusCode, data any) {
	msg := GetMessage(code)

	if statusCode >= 200 && statusCode < 300 {
		c.JSON(statusCode, dto.ResponseFormat{
			Message: msg,
			Data:    data,
		})
	} else {
		c.JSON(statusCode, dto.ResponseFormat{
			ErrorCode:    int(code),
			ErrorMessage: msg,
		})
	}
}

func AbortWithResponse(c *gin.Context, statusCode int, code StatusCode) {
	msg := GetMessage(code)

	c.AbortWithStatusJSON(statusCode, dto.ResponseFormat{
		ErrorCode:    int(code),
		ErrorMessage: msg,
	})
}
