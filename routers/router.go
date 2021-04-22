package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)
var Router *gin.Engine

func Setup() {
	Router = gin.New()
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	Router.Static("/storage","./storage")
}
