package service

import (
	"github.com/gin-gonic/gin"
	"ws/internal/auth"
	"ws/internal/models"
)

func Authenticate(c *gin.Context) {
	user := &models.BackendUser{}
	if user.Auth(c) {
		auth.PutBackendUser(c, user)
		c.Set("user", user)
	} else {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
