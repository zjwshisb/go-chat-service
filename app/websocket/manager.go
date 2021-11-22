package websocket

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
	"ws/app/databases"
	"ws/app/models"
)

type ConnManager interface {
	SendAction(act *Action, conn ...Conn)
	AddConn(connect Conn)
	RemoveConn(key int64)
	GetConn(key int64) (Conn, bool)
	GetAllConn() []Conn
	Register(connect Conn)
	Unregister(connect Conn)
	Ping()
	Run()
	ReceiveMessage(cm *ConnMessage)
}

type MessageHandle interface {
	handleReceiveMessage()
	handleMessage(cm *ConnMessage)
	handleRemoteMessage()
	deliveryMessage(msg *models.Message, from Conn, session *models.ChatSession)
}

type ManagerHook = func(conn Conn)

type ConnMessage struct {
	Conn *Client
	Action *Action
}

type manager struct {
	Clients map[int64]Conn
	lock    sync.RWMutex
	Channel string // 当前manager channel
	ConnMessages chan *ConnMessage // 接受client消息的chan
	onRegister ManagerHook //客户端连接成功hook
	onUnRegister ManagerHook //客户端连接断开hook
	userChannelCacheKey string // 客户端channel cache key
	groupCacheKey string // manager群组cache key
}

// 发布消息
func (m *manager) publish(channel string, payload *payload) error {
	ctx := context.Background()
	cmd := databases.Redis.Publish(ctx, channel, payload)
	return cmd.Err()
}

// 获取用户channel cache key
func (m *manager) getUserChannelKey(uid int64) string {
	return fmt.Sprintf(m.userChannelCacheKey, uid)
}

// 设置用户channel
func (m *manager) setUserChannel(uid int64)  {
	ctx := context.Background()
	key := m.getUserChannelKey(uid)
	databases.Redis.Set(ctx, key, m.GetSubscribeChannel(), 0)
}
// 移除用户channel
func (m *manager) removeUserChannel(uid int64)  {
	ctx := context.Background()
	databases.Redis.Del(ctx, m.getUserChannelKey(uid))
}
// 获取用户channel
func (m *manager) getUserChannel(uid int64) string {
	ctx := context.Background()
	key := m.getUserChannelKey(uid)
	cmd := databases.Redis.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return ""
	}
	return cmd.Val()
}
// 获取当前manager channel
func (m *manager) GetSubscribeChannel() string {
	return m.Channel
}
// 接受消息
func (m *manager) ReceiveMessage(cm *ConnMessage)  {
	m.ConnMessages <- cm
}

// 获取当前客户端数量
func (m *manager) GetTotal() int  {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.GetAllConn())
}

// 给客户端发送消息
func (m *manager) SendAction(a  *Action, clients ...Conn) {
	for _,c := range clients {
		c.Deliver(a)
	}
}
// 客户端是否存在
func (m *manager) ConnExist(uid int64) bool {
	_, exist := m.GetConn(uid)
	return exist
}
// 获取客户端
func (m *manager) GetConn(uid int64) (client Conn,ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	client, ok = m.Clients[uid]
	return
}
// 添加客户端
func (m *manager) AddConn(client Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Clients[client.GetUserId()] = client
}
// 移除客户端
func (m *manager) RemoveConn(uid int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.Clients, uid)
}
// 获取所有客户端
func (m *manager) GetAllConn() (s []Conn){
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make([]Conn, 0, len(m.Clients))
	for _, c := range m.Clients {
		r = append(r, c)
	}
	return r
}
// 客户端注销
func (m *manager) Unregister(client Conn) {
	existConn, exist := m.GetConn(client.GetUserId())
	if exist {
		if existConn == client {
			m.removeUserChannel(client.GetUserId())
			m.RemoveConn(client.GetUserId())
			client.close()
			if m.onUnRegister != nil {
				m.onUnRegister(client)
			}
		}
	}
}
// 客户端注册
// 如果是打开多个tab，关闭之前的连接
func (m *manager) Register(client Conn) {
	old , exist := m.GetConn(client.GetUserId())
	timer := time.After(1 * time.Second)
	if exist {
		m.SendAction(NewMoreThanOne(), old)
	}
	m.AddConn(client)
	m.setUserChannel(client.GetUserId())
	client.run()
	<-timer
	if m.onRegister != nil {
		m.onRegister(client)
	}
}

// 给所有客户端发送心跳
// 客户端因意外断开链接，服务器没有关闭事件，无法得知连接已关闭
// 通过心跳发送""字符串，如果发送失败，则调用conn的close方法回收
func (m *manager) Ping()  {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			conns := m.GetAllConn()
			ping := NewPing()
			m.SendAction(ping, conns...)
		}
	}
}
// 获取同类型的所有channel
func (m *manager) getAllChannel() []string  {
	ctx := context.Background()
	cmd := databases.Redis.SMembers(ctx, m.groupCacheKey)
	return cmd.Val()
}
// 注册频道
func (m *manager) registerChannel()  {
	ctx := context.Background()
	databases.Redis.SAdd(ctx, m.groupCacheKey, m.Channel)
}
// 移除频道
func (m *manager) unRegisterChannel()  {
	ctx := context.Background()
	databases.Redis.SRem(ctx, m.groupCacheKey, m.Channel)
}

func (m *manager) Run() {
	go m.Ping()
	m.registerChannel()

}
func (m *manager) destroy()  {
	m.unRegisterChannel()
}
