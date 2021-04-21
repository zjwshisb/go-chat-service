package hub

import (
	"sync"
	"ws/util"
)
type UserClientMap struct {
	Clients map[int64]*UClient // 客户端map
	lock sync.RWMutex
}

func (m *UserClientMap) GetClient(uid int64) (client *UClient,ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	client, ok = m.Clients[uid]
	return
}
func (m *UserClientMap) AddClient(client *UClient) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Clients[client.User.ID] = client
	return
}
func (m *UserClientMap) RemoveClient(uid int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.Clients, uid)
}
func (m *UserClientMap) GetAllClient() (s []*UClient){
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make([]*UClient, 0)
	for _, c := range m.Clients {
		r = append(r, c)
	}
	return r
}

type userHub struct {
	AcceptedClient *UserClientMap
	WaitingClient *UserClientMap
	util.Hook
}

func (hub *userHub) Change2waiting(client *UClient)  {
	hub.AcceptedClient.RemoveClient(client.User.ID)
	hub.WaitingClient.AddClient(client)
}

func (hub *userHub) Change2accept(client *UClient) {
	hub.WaitingClient.RemoveClient(client.User.ID)
	hub.AcceptedClient.AddClient(client)
}

func (hub *userHub) Logout(client *UClient) {
	defer func() {
		hub.TriggerHook(userLogout, client)
		Hub.Server.TriggerHook(userLogout, client)
	}()
	hub.AcceptedClient.RemoveClient(client.User.ID)
	hub.WaitingClient.RemoveClient(client.User.ID)
}

func (hub *userHub) Login(client *UClient) {
	if old, ok := hub.AcceptedClient.GetClient(client.User.ID); ok{
		old.close()
	}
	if old, ok := hub.WaitingClient.GetClient(client.User.ID); ok {
		old.close()
	}
	if client.ServerId > 0 {
		hub.AcceptedClient.AddClient(client)
	} else {
		hub.WaitingClient.AddClient(client)
	}
	client.Run()
	Hub.Server.TriggerHook(userLogin, client)
}
