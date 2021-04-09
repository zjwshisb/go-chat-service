package hub

import (
	"sync"
	"ws/util"
)
type userHub struct {
	Clients map[int64]*UClient // 已接入的客户端
	Lock sync.RWMutex
	Waiting map[int64]*UClient  //等待接入的客户端
	util.Hook
}

func (hub *userHub) getClient(id int64) (client *UClient,ok bool) {
	hub.Lock.RLock()
	defer hub.Lock.RUnlock()
	client, ok = hub.Clients[id]
	return
}
func (hub *userHub) getWaitClient(id int64) (client *UClient,ok bool) {
	hub.Lock.RLock()
	defer hub.Lock.RUnlock()
	client, ok = hub.Waiting[id]
	return
}

func (hub *userHub) Logout(client *UClient) {
	hub.Lock.Lock()
	defer func() {
		hub.Lock.Unlock()
		hub.TriggerHook(userLogout, client)
		Hub.Server.TriggerHook(userLogout, client)
	}()
	delete(hub.Clients, client.User.ID)
	delete(hub.Waiting, client.User.ID)
}

func (hub *userHub) Login(client *UClient) {
	if old, ok := hub.getClient(client.User.ID); ok{
		old.close()
	}
	if old, ok := hub.getWaitClient(client.User.ID); ok {
		old.close()
	}
	hub.Lock.Lock()
	defer hub.Lock.Unlock()
	if client.ServerId > 0 {

		hub.Clients[client.User.ID] = client

	} else {
		hub.Waiting[client.User.ID] = client
	}
	client.Run()
}