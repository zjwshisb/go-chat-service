package middleware

import (
	"github.com/gin-gonic/gin"
	"ws/models"
)

func Authenticate(c *gin.Context) {
	user := &models.ServiceUser{}
	user.Auth(c)
	if user.ID != 0 {
		c.Set("user", user)
	} else {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
