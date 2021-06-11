package service

import (
	"github.com/gin-gonic/gin"
	"ws/internal/models"
)

func getUser(c *gin.Context) *models.ServiceUser {
	ui, _ := c.Get("user")
	user := ui.(*models.ServiceUser)
	return user
}
