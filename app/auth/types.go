package auth

import "github.com/gin-gonic/gin"

type User interface {
	GetPrimaryKey() int64
	GetUsername() string
	GetAvatarUrl() string
	GetMpOpenId() string
	Auth(c *gin.Context) bool
}
