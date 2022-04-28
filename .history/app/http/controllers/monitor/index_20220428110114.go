package monitor

import (
	"net/http"
	"ws/app/rpc/rpcclient"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "monitor.tmpl", gin.H{
		"admin": rpcclient.ConnectionAllCount("admin"),
	})
}
