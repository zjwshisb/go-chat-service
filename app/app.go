package app

import (
	"ws/app/databases"
	_ "ws/app/http/requests"
	"ws/app/log"
	"ws/app/routers"
	"ws/app/websocket"
)


func Setup() {
	databases.Setup()
	log.Setup()
	websocket.Setup()
	routers.Setup()
}
