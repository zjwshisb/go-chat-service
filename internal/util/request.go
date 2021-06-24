package util

import (
	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) (token string) {
	token  = ""
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 {
		token = bearerToken[7:]
	}
	if token == "" {
		if queryToken, ok := c.GetQuery("token"); ok {
			token = queryToken
		}
	}
	return token
}

