package monitor

import (
	"net/http"
	"ws/app/rpc/rpcclient"
	"ws/config"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	var adminCount int64
	var userCount int64
	isCluster := config.IsCluster()
	if isCluster {
		adminCount = rpcclient.ConnectionAllCount("admin")
		userCount = rpcclient.ConnectionAllCount("user")
	}

	c.HTML(http.StatusOK, "monitor.tmpl", gin.H{
		"admin":     adminCount,
		"user":      userCount,
		"isCluster": isCluster,
	})
}
