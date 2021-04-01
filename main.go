package main

import (
	"ws/config"
	"ws/db"
	"ws/modules"
	"ws/routers"
)

func init()  {
	config.Setup()
	db.Setup()
	routers.Setup()
	modules.Setup()
}
func main() {
	routers.Router.Run(config.Http.Host +":" + config.Http.Port)
	//m := flag.String("m", "m", "模式")
	//fmt.Println(*m)
	//migrate.Run()
}
