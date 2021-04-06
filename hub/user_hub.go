package hub

import (
	"sync"
	"ws/util"
)
type userHub struct {
	Clients map[int64]*UClient // 已接入的客户端
	Lock sync.RWMutex
	Waiting map[int64]*UClient  //等待接入的客户端
	WaitingLock sync.RWMutex
	util.Hook
}

func (hub *userHub) getClient(id int64) (client *UClient,ok bool) {
	hub.Lock.RLock()
	defer hub.Lock.RUnlock()
	client, ok = hub.Clients[id]
	return
}
func (hub *userHub) Logout(client *UClient) {
	hub.Lock.Lock()
	defer func() {
		hub.Lock.Unlock()
		hub.TriggerHook(userLogout, client)
		Hub.Server.TriggerHook(userLogout, client)
	}()
	delete(hub.Clients, client.UserId)
	delete(hub.Waiting, client.UserId)
}

func (hub *userHub) Login(client *UClient) {
	defer hub.TriggerHook(userLogout)
	defer Hub.Server.TriggerHook(userLogin, client)
	if client.ServerId > 0 {
		hub.Lock.Lock()
		hub.Clients[client.UserId] = client
		defer hub.Lock.Unlock()
	} else {
		hub.WaitingLock.Lock()
		hub.Waiting[client.UserId] = client
		defer hub.WaitingLock.Unlock()
	}
	go client.Run()
}