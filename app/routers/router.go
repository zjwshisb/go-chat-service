package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"net/http"
	"strings"
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
	if strings.ToLower(viper.GetString("App.Env")) == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	Router = gin.New()
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	Router.Static(viper.GetString("File.LocalPrefix"), viper.GetString("File.LocalPath"))
	Router.Static("/public", viper.GetString("App.PublicPath"))
	Router.GET("/", func(c *gin.Context) {
		c.JSON(200, "hello world")
	})
	registerAdmin()
	registerFrontend()
}
