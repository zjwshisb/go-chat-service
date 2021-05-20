package internal

import (
	"ws/internal/databases"
	"ws/internal/log"
	"ws/internal/routers"
	"ws/internal/websocket"
)

func Setup() {
	databases.Setup()
	log.Setup()
	websocket.Setup()
	routers.Setup()
}