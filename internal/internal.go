package internal

import (
	"ws/internal/databases"
	"ws/internal/log"
	"ws/internal/websocket"
)

func Setup() {
	databases.Setup()
	log.Setup()
	websocket.Setup()
}