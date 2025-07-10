package response

import (
	"keeplo/internal/adapter/rest/dto"

	"github.com/gin-gonic/gin"
)

func HandleResponse(c *gin.Context, statusCode int, code int, data any) {
	if statusCode >= 200 && statusCode < 300 {
		c.JSON(statusCode, dto.ResponseFormat{
			Message: "",
			Data:    data,
		})
	} else {
		c.JSON(statusCode, dto.ResponseFormat{
			ErrorMessage: "",
			ErrorCode:    int(code),
		})
	}
}

func AbortWithResponse(c *gin.Context, statusCode int, code int) {
	c.AbortWithStatusJSON(statusCode, dto.ResponseFormat{
		ErrorMessage: "",
		ErrorCode:    int(code),
	})
}
