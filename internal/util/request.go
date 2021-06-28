package util

import (
	"github.com/gin-gonic/gin"
	"ws/configs"
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
func PublicAsset(path string) string {
	return configs.App.Url + "/public/" + path
}
