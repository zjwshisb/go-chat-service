package user

import (
	"github.com/gin-gonic/gin"
	"ws/internal/auth"
)

func Authenticate(c *gin.Context) {
	user := auth.UserGuard()
	if user.Auth(c) {
		auth.SetUser(c, user)
	} else {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
