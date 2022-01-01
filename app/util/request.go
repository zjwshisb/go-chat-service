package util

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	return viper.GetString("App.Url") + "/public/" + path
}

// SystemAvatar 系统头像
func SystemAvatar() string  {
	return PublicAsset("avatar.jpeg")
}