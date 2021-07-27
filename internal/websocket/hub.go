package websocket

import (
	"sync"
	"time"
	"ws/internal/action"
	"ws/internal/event"
)

const (
	UserLogin = iota
	UserLogout
)

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
	existConn, exist := hub.GetConn(client.GetUserId())
	if exist {
		if existConn == client {
			hub.RemoveConn(client.GetUserId())
			client.close()
			hub.Call(UserLogout, client)
		}
	}
}
func (hub *BaseHub) Login(client Conn) {
	old , exist := hub.GetConn(client.GetUserId())
	timer := time.After(1 * time.Second)
	if exist {
		hub.SendAction(action.NewMoreThanOne(), old)
	}
	hub.AddConn(client)
	client.run()
	<-timer
	hub.Call(UserLogin, client)
}
func (hub *BaseHub) Ping()  {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			conns := hub.GetAllConn()
			ping := action.NewPing()
			for _, conn := range conns {
				conn.Deliver(ping)
			}
		}
	}
}
func (hub *BaseHub) Run() {
	go hub.Ping()
}

var UserHub *userHub
var ServiceHub *serviceHub

func Setup()  {
	UserHub = &userHub{
		BaseHub{
			Clients: map[int64]Conn{},
		},
	}
	UserHub.Run()
	ServiceHub = &serviceHub{
		BaseHub{
			Clients: map[int64]Conn{},
		},
	}
	ServiceHub.Run()
}

