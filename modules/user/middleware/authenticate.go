package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ws/models"
)

func Authenticate(c *gin.Context) {
	user := &models.User{}
	user.Auth(c)
	fmt.Println(user)
	if user.ID != 0 {
		c.Set("user", user)
	} else {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
