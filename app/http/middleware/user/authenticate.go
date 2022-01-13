package user

import (
	"github.com/gin-gonic/gin"
	"ws/app/http/requests"
	"ws/app/models"
)

func Authenticate(c *gin.Context) {
	user := &models.User{}
	if user.Auth(c) {
		requests.SetUser(c, user)
	} else {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
