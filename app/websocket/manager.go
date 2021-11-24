package websocket

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"time"
	"ws/app/databases"
	"ws/app/models"
	"ws/app/mq"
	"ws/configs"
)
// conn管理相关方法
type ConnContainer interface {
	AddConn(connect Conn)
	GetConn(key int64) (Conn, bool)
	GetAllConn() []Conn
	GetTotal() int
	ConnExist(uid int64) bool
	Register(connect Conn)
	Unregister(connect Conn)
	RemoveConn(key int64)
}
// channel相关方法
// 集群模式下当读发消息的conn不在同一台机器时
// 通过订阅/发布进行消息通信
type ChannelManager interface {
	GetSubscribeChannel() string
	publish(channel string, payload *mq.Payload) error
	getUserChannelKey(uid int64) string
	setUserChannel(uid int64)
	removeUserChannel(uid int64)
	getUserChannel(uid int64) string
	getAllChannel() []string
	registerChannel()
	unRegisterChannel()
	isCluster() bool
}

type ConnManager interface {
	ConnContainer
	ChannelManager
	Run()
	Destroy()
	Ping()
	SendAction(act *Action, conn ...Conn)
	ReceiveMessage(cm *ConnMessage)
}

type MessageHandle interface {
	handleReceiveMessage()
	handleMessage(cm *ConnMessage)
	handleRemoteMessage()
	handleOffline(msg *models.Message)
	DeliveryMessage(msg *models.Message)
}

type ManagerHook = func(conn Conn)

type ConnMessage struct {
	Conn *Client
	Action *Action
}

type manager struct {
	Clients map[int64]Conn
	lock    sync.RWMutex // clients的读写锁
	Channel string // 当前manager channel
	ConnMessages chan *ConnMessage // 接受从conn所读取消息的chan
	onRegister ManagerHook //客户端连接成功hook
	onUnRegister ManagerHook //客户端连接断开hook
	userChannelCacheKey string // 客户端channel cache key
	groupCacheKey string // manager群组cache key
}

func (m *manager) isCluster() bool {
	return configs.App.Cluster
}
// 发布消息
func (m *manager) publish(channel string, payload *mq.Payload) error {
	err := mq.Mq().Publish(channel, payload)
	return err
}

// 获取用户channel cache key
func (m *manager) getUserChannelKey(uid int64) string {
	return fmt.Sprintf(m.userChannelCacheKey, uid)
}

// 设置用户所在channel
// 默认有效期24小时，用于程序意外退出后的清理
func (m *manager) setUserChannel(uid int64)  {
	if m.isCluster() {
		ctx := context.Background()
		key := m.getUserChannelKey(uid)
		databases.Redis.Set(ctx, key, m.GetSubscribeChannel(), time.Hour * 24)
	}
}
// 移除用户所在channel
func (m *manager) removeUserChannel(uid int64)  {
	if m.isCluster() {
		ctx := context.Background()
		databases.Redis.Del(ctx, m.getUserChannelKey(uid))
	}
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
			if m.onUnRegister != nil {
				m.onUnRegister(client)
			}
		}
	}
}
// 客户端注册
// 如果是重复打开页面，关闭之前连接
// 集群模式下，如果不在本机则投递一个消息
func (m *manager) Register(client Conn) {
	old , exist := m.GetConn(client.GetUserId())
	timer := time.After(1 * time.Second)
	if exist {
		m.SendAction(NewMoreThanOne(), old)
	} else {
		if m.isCluster() {
			oldChannel := m.getUserChannel(client.GetUserId())
			if oldChannel != "" {
				m.publish(oldChannel, &mq.Payload{
					Types: mq.TypeAdminLogin,
					Data: client.GetUserId(),
				})
			}
		}
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
// 通过心跳发送""字符串，如果发送失败，则调用conn的close方法执行清理
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
// 同事清理掉意外退出的机器的channel
func (m *manager) getAllChannel() []string  {
	ctx := context.Background()
	now := time.Now().Unix()
	fz := now -  (60 * 60 * 1 + 10)
	cmd := databases.Redis.ZRangeByScore(ctx, m.groupCacheKey, &redis.ZRangeBy{
		Min:   strconv.FormatInt(fz, 10),
		Max:    "+inf",
		Offset: 0,
		Count:  0,
	})
	// 清理失效的channel
	databases.Redis.ZRemRangeByScore(ctx , m.groupCacheKey, "-inf", strconv.FormatInt(fz, 10))
	return cmd.Val()
}
// 集群模式下
// 注册频道
// 心跳更新最后时间，用于程序意外退出后的清理
func (m *manager) registerChannel()  {
	fn := func() {
		ctx := context.Background()
		databases.Redis.ZAdd(ctx, m.groupCacheKey, &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: m.Channel,
		})
	}
	fn()
	go func() {
		tinker := time.NewTicker(time.Minute)
		for {
			<-tinker.C
			fn()
		}
	}()

}
// 移除频道
func (m *manager) unRegisterChannel()  {
	if m.isCluster() {
		ctx := context.Background()
		databases.Redis.ZRem(ctx, m.groupCacheKey, m.Channel)
	}
}

func (m *manager) Run() {
	go m.Ping()
	if m.isCluster() {
		go m.registerChannel()
	}
}

// 释放相关资源
func (m *manager) Destroy()  {
	if m.isCluster() {
		m.unRegisterChannel()
		conns := m.GetAllConn()
		for _, conn := range conns {
			m.removeUserChannel(conn.GetUserId())
		}
	}
}
