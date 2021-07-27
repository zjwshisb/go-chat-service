package auth

import (
	"github.com/gin-gonic/gin"
	"ws/internal/models"
)


func UserGuard() User {
	return &models.User{}
}

func SetUser(c *gin.Context, user User)  {
	c.Set("frontend", user)
}

func GetUser(c *gin.Context) User {
	ui, _ := c.Get("frontend")
	user := ui.(User)
	return user
}
