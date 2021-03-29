package http

import (
	"github.com/gin-gonic/gin"
	"ws/models"
)

func Me(c *gin.Context) {
	ui, _ := c.Get("user")
	user := ui.(*models.ServerUser)
	c.JSON(200, gin.H{
		"username": user.Username,
		"id": user.ID,
	})
}
