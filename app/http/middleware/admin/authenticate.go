package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/http/requests"
	"ws/app/repositories"
)

func Authenticate(c *gin.Context) {
	token := requests.GetToken(c)
	uid, err := requests.ParseToken(token)
	if err == nil {
		admin := repositories.AdminRepo.First([]*repositories.Where{
			{
				Filed: "id = ?",
				Value: uid,
			},
		}, []string{})
		if admin != nil {
			requests.SetAdmin(c, admin)
			return
		}
	}
	c.JSON(401, gin.H{
		"message": "Unauthorized",
	})
	c.Abort()
}
