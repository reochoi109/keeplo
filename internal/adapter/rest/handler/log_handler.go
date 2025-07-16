package handler

import "github.com/gin-gonic/gin"

func GetMonitorHealthLogHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message ": "GetMonitorHealthLogHandler"})
}

func GetMonitorStatusHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message ": "GetMonitorStatusHandler"})
}
