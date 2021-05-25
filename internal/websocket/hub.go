package websocket

import (
	"sync"
	"ws/internal/action"
	"ws/internal/event"
)

const (
	UserLogin = iota
	UserLogout
)

type Manager interface {
	SendAction(act *action.Action, conn ...Conn)
	AddConn(connect Conn)
	RemoveConn(key int64)
	GetConn(key int64) (Conn, bool)
	GetAllConn() []Conn
	Login(connect Conn)
	Logout(connect Conn)
}
type BaseHub struct {
	Clients map[int64]Conn
	lock    sync.RWMutex
	event.BaseEvent
}
func (hub *BaseHub) SendAction(a  *action.Action, clients ...Conn) {
	for _,c := range clients {
		c.Deliver(a)
	}
}
func (hub *BaseHub) GetConn(uid int64) (client Conn,ok bool) {
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	client, ok = hub.Clients[uid]
	return
}
func (hub *BaseHub) AddConn(client Conn) {
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
func (hub *BaseHub) GetAllConn() (s []Conn){
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	r := make([]Conn, 0)
	for _, c := range hub.Clients {
		r = append(r, c)
	}
	return r
}
func (hub *BaseHub) Logout(client Conn) {
	hub.RemoveConn(client.GetUserId())
	client.close()
	go hub.Call(UserLogout, client)
}

func (hub *BaseHub) Login(client Conn) {
	hub.AddConn(client)
	client.run()
	go hub.Call(UserLogin, client)
}

var UserHub *userHub
var ServiceHub *serviceHub

func Setup()  {
	UserHub = &userHub{
		BaseHub{
			Clients: map[int64]Conn{},
		},
	}
	ServiceHub = &serviceHub{
		BaseHub{
			Clients: map[int64]Conn{},
		},
	}
	ServiceHub.Setup()
}

