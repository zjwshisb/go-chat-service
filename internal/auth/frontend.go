package auth

import (
	"github.com/gin-gonic/gin"
	"ws/internal/models"
)

func UserGuard() User {
	return &models.User{}
}
func SetUser(c *gin.Context, user User)  {
	c.Set("user", user)
}
func GetUser(c *gin.Context) User {
	ui, _ := c.Get("user")
	user := ui.(User)
	return user
}

type User interface {
	GetPrimaryKey() int64
	GetUsername() string
	GetAvatarUrl() string
	Auth(c *gin.Context) bool
}
