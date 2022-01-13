package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/http/requests"
	"ws/app/models"
)

func Authenticate(c *gin.Context) {
	user := &models.Admin{}
	if user.Auth(c) {
		requests.SetAdmin(c, user)
	} else {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
