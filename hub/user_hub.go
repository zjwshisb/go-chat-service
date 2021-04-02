package hub

import (
	"errors"
	"sync"
)
type userHub struct {
	Clients map[int64]*UClient
	Lock sync.RWMutex
	Waiting []*UClient
	WaitingLock sync.RWMutex
}

func (hub *userHub) getClient(id int64) (client *UClient,err error) {
	hub.Lock.RLock()
	defer hub.Lock.RUnlock()
	client, ok := hub.Clients[id]
	if !ok {
		err = errors.New("client not exists")
	}
	return
}
func (hub *userHub) Logout(client *UClient) {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()
	delete(hub.Clients, client.UserId)
}

func (hub *userHub) Login(client *UClient) {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()
	hub.Clients[client.UserId] = client
	go client.Run()
}