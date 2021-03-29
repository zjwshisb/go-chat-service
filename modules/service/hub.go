package service

import (
	"log"
	"sync"
)

type hub struct {
	Clients map[int64]*Client
	Login chan *Client
	Logout chan *Client
	Lock sync.Mutex
}
func (hub *hub) run()  {
	for {
		select {
		case client := <- hub.Login:
			hub.Lock.Lock()
			hub.Clients[client.UserId] = client
			log.Print("login")
			hub.Lock.Unlock()
		case client := <- hub.Logout:
			hub.Lock.Lock()
			delete(hub.Clients, client.UserId)
			log.Print("logout")
			hub.Lock.Unlock()
		}
	}
}

