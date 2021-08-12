package app

import (
	_ "ws/app/databases"
	_ "ws/app/http/requests"
	_ "ws/app/log"
	_ "ws/app/routers"
	_ "ws/app/routers"
	"ws/app/websocket"
)


func Setup() {
	websocket.Setup()
}
