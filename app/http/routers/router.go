package routers

import (
	"embed"
	"html/template"
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

//go:embed templates/monitor.tmpl
var f embed.FS

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
	templ := template.Must(template.New("").ParseFS(f, "templates/*.tmpl"))
	Router.SetHTMLTemplate(templ)
	if config.GetEnv() == "local" {
		Router.GET("/monitor", monitor.Index)
	}
	registerAdmin()
	registerFrontend()
}
