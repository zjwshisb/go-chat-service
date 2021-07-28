package auth

import (
	"github.com/gin-gonic/gin"
	"ws/app/models"
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
