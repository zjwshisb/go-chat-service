package auth

import (
	"github.com/gin-gonic/gin"
	"ws/internal/models"
)

type User interface {
	GetPrimaryKey() int64
	GetUsername() string
	GetAvatarUrl() string
	GetMpOpenId() string
	Auth(c *gin.Context) bool
}

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
