package websocket

import (
	"sync"
	"time"
)

const (
	UserLogin = iota
	UserLogout
)
type BaseHub struct {
	Clients map[int64]Conn
	lock    sync.RWMutex
	baseEvent
}
// 获取当前客户端数量
func (hub *BaseHub) GetTotal() int  {
	return len(hub.GetAllConn())
}
// 给客户端发送消息
func (hub *BaseHub) SendAction(a  *Action, clients ...Conn) {
	for _,c := range clients {
		c.Deliver(a)
	}
}
// 获取客户端
func (hub *BaseHub) GetConn(uid int64) (client Conn,ok bool) {
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	client, ok = hub.Clients[uid]
	return
}
// 添加客户端
func (hub *BaseHub) AddConn(client Conn) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	hub.Clients[client.GetUserId()] = client
}
// 移除客户端
func (hub *BaseHub) RemoveConn(uid int64) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	delete(hub.Clients, uid)
}
// 获取所有客户端
func (hub *BaseHub) GetAllConn() (s []Conn){
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	r := make([]Conn, 0)
	for _, c := range hub.Clients {
		r = append(r, c)
	}
	return r
}
// 客户端登出
func (hub *BaseHub) Logout(client Conn) {
	existConn, exist := hub.GetConn(client.GetUserId())
	if exist {
		if existConn == client {
			hub.RemoveConn(client.GetUserId())
			client.close()
		}
	}
	hub.Call(UserLogout, client)
}
// 客户端登入
func (hub *BaseHub) Login(client Conn) {
	old , exist := hub.GetConn(client.GetUserId())
	timer := time.After(1 * time.Second)
	if exist { // 如果是打开多个tab，关闭之前的连接
		hub.SendAction(NewMoreThanOne(), old)
	}
	hub.AddConn(client)
	client.run()
	<-timer
	hub.Call(UserLogin, client)
}
// 给所有客户端发送心跳
// 客户端因意外断开链接，服务器没有关闭事件，无法得知连接已关闭
// 通过心跳发送""字符串，如果发送失败，则调用conn的close方法回收
func (hub *BaseHub) Ping()  {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			conns := hub.GetAllConn()
			ping := NewPing()
			hub.SendAction(ping, conns...)
		}
	}
}
func (hub *BaseHub) Run() {
	go hub.Ping()
}

var UserHub *userHub
var AdminHub *adminHub

func Setup()  {
	UserHub = &userHub{
		BaseHub{
			Clients: map[int64]Conn{},
		},
	}
	UserHub.Run()
	AdminHub = &adminHub{
		BaseHub{
			Clients: map[int64]Conn{},
		},
	}
	AdminHub.Run()
}

