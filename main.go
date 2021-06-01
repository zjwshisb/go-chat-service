package main

import (
	"ws/configs"
	"ws/internal"
	"ws/internal/routers"
)

func main() {
	internal.Setup()
	routers.Router.Run(configs.Http.Host +":" + configs.Http.Port)
}
