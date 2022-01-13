package requests

import (
	"github.com/gin-gonic/gin"
	"ws/app/contract"
	"ws/app/models"
)

func GetAdmin(c *gin.Context) contract.User {
	ui, _ := c.Get("admin")
	user := ui.(contract.User)
	return user
}

func SetAdmin(c *gin.Context, user *models.Admin) {
	c.Set("admin", user)
}

func SetUser(c *gin.Context, user contract.User)  {
	c.Set("frontend", user)
}

func GetUser(c *gin.Context) contract.User {
	ui, _ := c.Get("frontend")
	user := ui.(contract.User)
	return user
}
