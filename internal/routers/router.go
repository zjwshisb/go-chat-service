package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/configs"
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
	Router = gin.New()
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	Router.Static(configs.File.LocalPrefix, configs.File.LocalPath)
	Router.Static("/public", "./public")
	registerBackend()
	registerFrontend()
}
