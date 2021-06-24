package backend

import (
	"github.com/gin-gonic/gin"
	"ws/internal/auth"
	"ws/internal/models"
)

func Authenticate(c *gin.Context) {
	user := &models.BackendUser{}
	if user.Auth(c) {
		auth.SetBackendUser(c, user)
	} else {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
