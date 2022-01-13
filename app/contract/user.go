package contract

import "github.com/gin-gonic/gin"

type FrontendUser interface {
	GetMpOpenId() string
}

type User interface {
	GetPrimaryKey() int64
	GetUsername() string
	GetAvatarUrl() string
	Auth(c *gin.Context) bool
	GetGroupId() int64
	GetReqId() string
}
