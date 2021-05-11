package main

import (
	"ws/configs"
	"ws/internal/databases"
	"ws/internal/routers"
	hub "ws/internal/websocket"
)

func init()  {
	configs.Setup()
	databases.Setup()
	routers.Setup()
	hub.Setup()
}
func main() {
	routers.Router.Run(configs.Http.Host +":" + configs.Http.Port)
	//m := flag.String("m", "m", "模式")
	//fmt.Println(*m)
	//migrate.Seed()
}
