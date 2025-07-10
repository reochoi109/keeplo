package handler

import (
	"github.com/gin-gonic/gin"
)

func SingupHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message ": "SingupHandler"})
}

func LoginHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message ": "LoginHandler"})
}

func GetUserHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message ": "GetUserHandler"})
}

func LogoutHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "LogoutHandler"})
}
