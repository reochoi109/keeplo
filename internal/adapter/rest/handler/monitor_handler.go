package handler

import "github.com/gin-gonic/gin"

func RegisterMonitorHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "RegisterMonitorHandler"})
}

func GetMonitorListHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GetMonitorListHandler"})
}

func GetMonitorHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GetMonitorHandler"})
}

func UpdateMonitorHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "UpdateMonitorHandler"})
}

func RemoveMonitorHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "RemoveMonitorHandler"})
}
