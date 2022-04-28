package routers

import (
	"net/http"
	"strings"
	"ws/app/http/controllers/monitor"
	"ws/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Router *gin.Engine

var (
	upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Setup() {
	if strings.ToLower(config.GetEnv()) == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	Router = gin.New()

	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	Router.Static("/assets", config.GetStoragePath()+"/assets")
	Router.GET("/", func(c *gin.Context) {
		c.JSON(200, "hello world")
	})
	if config.GetEnv() == "local" {
		Router.LoadHTMLGlob("templates/*")
		Router.GET("/monitor", monitor.Index)
	}
	registerAdmin()
	registerFrontend()
}
