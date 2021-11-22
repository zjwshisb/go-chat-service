package auth

import "github.com/gin-gonic/gin"

type PrimaryUser interface {
	GetPrimaryKey() int64
}

type User interface {
	PrimaryUser
	GetUsername() string
	GetAvatarUrl() string
	GetMpOpenId() string
	Auth(c *gin.Context) bool
	GetReqId() string
}
