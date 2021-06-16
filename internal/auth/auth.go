package auth

import (
	"github.com/gin-gonic/gin"
	"ws/internal/models"
)

func GetBackendUser(c *gin.Context) *models.BackendUser {
	ui, _ := c.Get("user")
	user := ui.(*models.BackendUser)
	return user
}
func PutBackendUser(c *gin.Context, user *models.BackendUser) {
	c.Set("user", user)
}

