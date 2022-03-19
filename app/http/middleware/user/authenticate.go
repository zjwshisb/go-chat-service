package user

import (
	"github.com/gin-gonic/gin"
	"ws/app/http/requests"
	"ws/app/repositories"
)

func Authenticate(c *gin.Context) {
	token := requests.GetToken(c)
	if token != "" {
		user := repositories.UserRepo.First([]*repositories.Where{
			{
				Filed: "api_token = ?",
				Value: token,
			},
		}, []string{})
		if user != nil {
			requests.SetUser(c, user)
			return
		}
	}
	c.JSON(401, gin.H{
		"message": "Unauthorized",
	})
	c.Abort()
}
