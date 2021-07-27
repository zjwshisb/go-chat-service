package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, "404 not found")
}

func RespSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"success": true,
		"code": 0,
	})
}
func RespPagination(c *gin.Context, p interface{}) {
	c.JSON(http.StatusOK, p)

}
func RespValidateFail(c *gin.Context, msg interface{}) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": msg,
		"success": false,
	})
}
func RespFail(c *gin.Context, msg interface{} , code int) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"code":code,
		"message": msg,
	})
}
func RespError(c *gin.Context, msg interface{}){
	c.JSON(500, gin.H{
		"message": msg,
	})
}
