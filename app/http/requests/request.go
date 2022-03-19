package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"net/http"
)

func GetToken(c *gin.Context) (token string) {
	token = ""
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

func CreateToken(uid string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
	})
	token, err := at.SignedString([]byte(viper.GetString("App.Secret")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(token string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("App.Secret")), nil
	})
	if err != nil {
		return "", err
	}
	return claim.Claims.(jwt.MapClaims)["uid"].(string), nil
}

func RespSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"success": true,
		"code":    0,
	})
}
func RespPagination(c *gin.Context, p interface{}) {
	c.JSON(http.StatusOK, p)
}
func RespValidateFail(c *gin.Context, msg interface{}) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": msg,
	})
}
func RespFail(c *gin.Context, msg interface{}, code int) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"code":    code,
		"message": msg,
	})
}
func RespError(c *gin.Context, msg interface{}) {
	c.JSON(500, gin.H{
		"message": msg,
	})
}
func RespNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, "404 not found")
}
