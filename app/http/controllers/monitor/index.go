package monitor

import (
	"net/http"
	"ws/app/http/websocket"
	"ws/app/rpc/client"
	"ws/config"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	var adminCount int64
	var userCount int64
	isCluster := config.IsCluster()
	serverStr := ""
	if isCluster {
		adminCount = client.ConnectionAllCount("admin")
		userCount = client.ConnectionAllCount("user")
		d := client.NewDiscovery("Connection")
		services := d.GetServices()
		for _, s := range services {
			if serverStr != "" {
				serverStr += "</br>"
			}
			serverStr += s.Key
		}
	} else {
		adminCount = websocket.AdminManager.GetAllConnCount()
		userCount = websocket.UserManager.GetAllConnCount()
	}

	c.HTML(http.StatusOK, "monitor.tmpl", gin.H{
		"admin":     adminCount,
		"user":      userCount,
		"isCluster": isCluster,
		"server":    serverStr,
	})
}
