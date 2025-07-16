package middleware

import (
	"context"
	"keeplo/pkg/idgen"
	"keeplo/pkg/logger"

	"github.com/gin-gonic/gin"
)

func UseTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := idgen.GenerateTraceID()
		ctx := context.WithValue(c.Request.Context(), logger.ContextTraceID, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
