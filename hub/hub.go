package hub

import (
	"sync"
	"ws/action"
)

type WebsocketManager interface {
	SendAction(act *action.Action, conn ...WebsocketConn)
	AddConn(connect WebsocketConn)
	RemoveConn(key int64)
	GetConn(key int64) (WebsocketConn, bool)
	GetAllConn() []WebsocketConn
	Login(connect WebsocketConn)
	Logout(connect WebsocketConn)
}
type BaseHub struct {
	Clients map[int64]WebsocketConn
	lock    sync.RWMutex
}
func (hub *BaseHub) SendAction(a  *action.Action, clients ...WebsocketConn) {
	for _,c := range clients {
		c.Deliver(a)
	}
}
func (hub *BaseHub) GetConn(uid int64) (client WebsocketConn,ok bool) {
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	client, ok = hub.Clients[uid]
	return
}
func (hub *BaseHub) AddConn(client WebsocketConn) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	hub.Clients[client.GetUserId()] = client
	return
}
func (hub *BaseHub) RemoveConn(uid int64) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	delete(hub.Clients, uid)
}
func (hub *BaseHub) GetAllConn() (s []WebsocketConn){
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	r := make([]WebsocketConn, 0)
	for _, c := range hub.Clients {
		r = append(r, c.(WebsocketConn))
	}
	return r
}
func (hub *BaseHub) Logout(client WebsocketConn) {
	hub.RemoveConn(client.GetUserId())
	client.close()
}

func (hub *BaseHub) Login(client WebsocketConn) {
	hub.AddConn(client)
	client.run()
}

var UserHub *userHub
var ServiceHub *serviceHub

func Setup()  {
	UserHub = &userHub{
		BaseHub{
			Clients: map[int64]WebsocketConn{},
		},
	}
	ServiceHub = &serviceHub{
		BaseHub{
			Clients: map[int64]WebsocketConn{},
		},
	}
}

