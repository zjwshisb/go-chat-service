package responses

import (
	"net/http"
	"ws/app/repositories"

	"github.com/gin-gonic/gin"
)

func RespSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"success": true,
		"code":    0,
	})
}
func RespPagination[T any](c *gin.Context, p *repositories.Pagination[T]) {
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
