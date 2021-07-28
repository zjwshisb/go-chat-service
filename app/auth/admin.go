package auth

import (
	"github.com/gin-gonic/gin"
	"ws/app/models"
)

func GetAdmin(c *gin.Context) *models.Admin {
	ui, _ := c.Get("admin")
	user := ui.(*models.Admin)
	return user
}

func SetAdmin(c *gin.Context, user *models.Admin) {
	c.Set("admin", user)
}

