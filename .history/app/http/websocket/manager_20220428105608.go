package websocket

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
	"ws/app/contract"
	"ws/app/databases"
	"ws/app/models"
	"ws/app/mq"
	"ws/app/rpc/rpcclient"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// ConnContainer 管理相关方法
type ConnContainer interface {
	AddConn(conn Conn)
	GetConn(user contract.User) (Conn, bool)
	publishMoreThanOne(user contract.User)
	GetAllConn(gid int64) []Conn
	GetUserUuid(user contract.User) string
	SetUserUuid(user contract.User, uuid string)
	GetOnlineTotal(gid int64) int64
	ConnExist(user contract.User) bool
	Register(connect Conn)
	Unregister(connect Conn)
	RemoveConn(user contract.User)
	IsOnline(user contract.User) bool
	GetOnlineUserIds(gid int64) []int64
	Do(c func(), f func())
}

// ChannelManager channel相关方法
// 集群模式下当读发消息的conn不在同一台机器时
// 通过订阅/发布进行消息通信
type ChannelManager interface {
	GetSubscribeChannel() string
	publish(channel string, payload *mq.Payload) error
	publishToAllChannel(payload *mq.Payload)
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
	GetTypes() string
}

type MessageHandle interface {
	handleReceiveMessage()
	handleMessage(cm *ConnMessage)
	handleRemoteMessage()
	handleOffline(msg *models.Message)
	DeliveryMessage(msg *models.Message, remote bool)
}

type ManagerHook = func(conn Conn)

type ConnMessage struct {
	Conn   *Client
	Action *Action
}
type Shard struct {
	m     map[int64]Conn
	mutex sync.RWMutex
}

func (s *Shard) GetAll() []Conn {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	conns := make([]Conn, 0, len(s.m))
	for _, conn := range s.m {
		conns = append(conns, conn)
	}
	return conns
}

func (s *Shard) GetTotalCount() int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return int64(len(s.m))
}

func (s *Shard) Get(uid int64) (conn Conn, exist bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	conn, exist = s.m[uid]
	return
}
func (s *Shard) Set(conn Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[conn.GetUserId()] = conn
}
func (s *Shard) Remove(uid int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.m, uid)
}

type manager struct {
	shardCount   int64             // 分组数量
	shard        []*Shard          //
	Channel      string            // 当前manager channel name
	ConnMessages chan *ConnMessage // 接受从conn所读取消息的chan
	onRegister   ManagerHook       //conn连接成功hook
	onUnRegister ManagerHook       //conn连接断开hook
	types        string            //类型
}

func (m *manager) GetTypes() string {
	return m.types
}
func (m *manager) Do(clusterFunc func(), single func()) {
	if m.isCluster() {
		if clusterFunc != nil {
			clusterFunc()
		}
	} else {
		if single != nil {
			single()
		}
	}
}

func (m *manager) getMod(gid int64) int64 {
	return gid % m.shardCount
}

func (m *manager) getSpread(gid int64) *Shard {
	return m.shard[m.getMod(gid)]
}

func (m *manager) isCluster() bool {
	return viper.GetBool("App.Cluster")
}

func (m *manager) UserUuidKey(uid int64) string {
	return fmt.Sprintf(m.types+":%d:uuid", uid)
}

// GetUserUuid 获取用户的当前连接uuid
func (m *manager) GetUserUuid(user contract.User) string {
	cmd := databases.Redis.Get(context.Background(), m.UserUuidKey(user.GetPrimaryKey()))
	return cmd.Val()
}

// SetUserUuid 设置用户的当前连接的uuid
func (m *manager) SetUserUuid(user contract.User, uuid string) {
	databases.Redis.Set(context.Background(), m.UserUuidKey(user.GetPrimaryKey()), uuid, time.Hour*24)
}

// 发布消息
func (m *manager) publish(channel string, payload *mq.Payload) error {
	err := mq.Publish(channel, payload)
	return err
}

//
func (m *manager) publishToAllChannel(payload *mq.Payload) {
	channels := m.getAllChannel()
	fmt.Println(channels)
	for _, channel := range channels {
		_ = m.publish(channel, payload)
	}
}

// 获取用户channel cache key
func (m *manager) getUserChannelKey(uid int64) string {
	return fmt.Sprintf("%s:%d:channel", m.types, uid)
}

// 设置用户所在channel为当前manager
func (m *manager) setUserChannel(uid int64) {
	m.Do(func() {
		ctx := context.Background()
		key := m.getUserChannelKey(uid)
		databases.Redis.Set(ctx, key, m.GetSubscribeChannel(), time.Hour*24*2)
	}, nil)
}

// 移除用户所在channel
func (m *manager) removeUserChannel(uid int64) {
	m.Do(func() {
		ctx := context.Background()
		databases.Redis.Del(ctx, m.getUserChannelKey(uid))
	}, nil)
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

// GetSubscribeChannel 获取当前manager channel name
func (m *manager) GetSubscribeChannel() string {
	return m.Channel
}

// ReceiveMessage 接受消息
func (m *manager) ReceiveMessage(cm *ConnMessage) {
	m.ConnMessages <- cm
}

// websocket 重复链接
func (m *manager) publishMoreThanOne(user contract.User) {
	m.Do(func() {
		oldChannel := m.getUserChannel(user.GetPrimaryKey())
		if oldChannel != "" {
			_ = m.publish(oldChannel, &mq.Payload{
				Types: mq.TypeMoreThanOne,
				Data:  user.GetPrimaryKey(),
			})
		}
	}, func() {
		m.noticeMoreThanOne(user)
	})
}

func (m *manager) noticeMoreThanOne(user contract.User) {
	oldConn, ok := m.GetConn(user)
	if ok && oldConn.GetUuid() != m.GetUserUuid(user) {
		m.SendAction(NewMoreThanOne(), oldConn)
	}
}

// GetOnlineUserIds 获取groupId对应的在线userIds
func (m *manager) GetOnlineUserIds(gid int64) []int64 {
	if m.isCluster() {
		return rpcclient.ConnectionIds(gid, m.types)
	} else {
		return m.GetLocalOnlineUserIds(gid)
	}
}
func (m *manager) GetLocalOnlineUserIds(gid int64) []int64 {
	s := m.getSpread(gid)
	allConn := s.GetAll()
	ids := make([]int64, 0)
	for _, conn := range allConn {
		if conn.GetGroupId() == gid {
			ids = append(ids, conn.GetUserId())
		}
	}
	return ids
}

// GetLocalOnlineTotal 获取本地groupId对应在线客户端数量
func (m *manager) GetLocalOnlineTotal(gid int64) int64 {
	s := m.getSpread(gid)
	return s.GetTotalCount()
}

// GetOnlineTotal 获取groupId对应在线客户端数量
func (m *manager) GetOnlineTotal(gid int64) int64 {
	if m.isCluster() {
		return rpcclient.ConnectionTotal(gid, m.types)
	}
	return m.GetLocalOnlineTotal(gid)
}

// IsOnline 用户是否在线
func (m *manager) IsOnline(user contract.User) bool {
	if m.isCluster() {
		return rpcclient.ConnectionExist(user.GetPrimaryKey())
	} else {
		return m.ConnExist(user)
	}
}

// SendAction 给客户端发送消息
func (m *manager) SendAction(a *Action, clients ...Conn) {
	for _, c := range clients {
		c.Deliver(a)
	}
}

// ConnExist 连接是否存在
func (m *manager) ConnExist(user contract.User) bool {
	_, exist := m.GetConn(user)
	return exist
}

// GetConn 获取客户端
func (m *manager) GetConn(user contract.User) (client Conn, ok bool) {
	s := m.getSpread(user.GetGroupId())
	client, ok = s.Get(user.GetPrimaryKey())
	return
}

// AddConn 添加客户端
func (m *manager) AddConn(conn Conn) {
	s := m.getSpread(conn.GetGroupId())
	s.Set(conn)
}

// RemoveConn 移除客户端
func (m *manager) RemoveConn(user contract.User) {
	s := m.getSpread(user.GetGroupId())
	s.Remove(user.GetPrimaryKey())
}

// GetAllConn 获取所有客户端
func (m *manager) GetAllConn(groupId int64) (conns []Conn) {
	s := m.getSpread(groupId)
	return s.GetAll()
}

func (m *manager) GetAllConnCount() int64 {
	var count int64
	for gid := range m.shard {
		count += m.GetOnlineTotal(int64(gid))
	}
	return count
}

func (m *manager) GetTotalConn() []Conn {
	conns := make([]Conn, 0)
	for gid := range m.shard {
		conns = append(conns, m.GetAllConn(int64(gid))...)
	}
	return conns
}

// Unregister 客户端注销
func (m *manager) Unregister(conn Conn) {
	s := m.getSpread(conn.GetGroupId())
	existConn, exist := s.Get(conn.GetUserId())
	if exist {
		if existConn == conn {
			m.RemoveConn(conn.GetUser())
			if m.onUnRegister != nil {
				m.onUnRegister(conn)
			}
		}
	}
}

// Register 客户端注册
// 先处理是否重复连接
// 集群模式下，如果不在本机则投递一个消息
func (m *manager) Register(conn Conn) {
	timer := time.After(1 * time.Second)
	m.publishMoreThanOne(conn.GetUser())
	m.AddConn(conn)
	m.SetUserUuid(conn.GetUser(), conn.GetUuid())
	m.setUserChannel(conn.GetUserId())
	conn.run()
	<-timer
	if m.onRegister != nil {
		m.onRegister(conn)
	}
}

// Ping
// 给所有客户端发送心跳
// 客户端因意外断开链接，服务器没有关闭事件，无法得知连接已关闭
// 通过心跳发送""字符串，如果发送失败，则调用conn的close方法执行清理
func (m *manager) Ping() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			ping := NewPing()
			for _, s := range m.shard {
				conns := s.GetAll()
				m.SendAction(ping, conns...)
			}
		}
	}
}
func (m *manager) getClusterKey() string {
	return fmt.Sprintf("%s:cluster:group", m.types)
}

// 获取同类型的所有channel
func (m *manager) getAllChannel() []string {
	ctx := context.Background()
	now := time.Now().Unix()
	fz := now - (60 * 2)
	cmd := databases.Redis.ZRangeByScore(ctx, m.getClusterKey(), &redis.ZRangeBy{
		Min:    strconv.FormatInt(fz, 10),
		Max:    "+inf",
		Offset: 0,
		Count:  0,
	})
	return cmd.Val()
}

// 集群模式下
// 注册频道
// 心跳更新最后时间，用于程序意外退出后的清理
func (m *manager) registerChannel() {
	fn := func() {
		ctx := context.Background()
		databases.Redis.ZAdd(ctx, m.getClusterKey(), &redis.Z{
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

// ClearInactiveChannel 清理无效的channel
func (m *manager) ClearInactiveChannel() {
	m.Do(func() {
		ctx := context.Background()
		now := time.Now().Unix()
		fz := now - (60 * 2)
		databases.Redis.ZRemRangeByScore(ctx, m.getClusterKey(), "-inf", strconv.FormatInt(fz, 10))
	}, nil)
}

// 移除频道
func (m *manager) unRegisterChannel() {
	m.Do(func() {
		ctx := context.Background()
		databases.Redis.ZRem(ctx, m.getClusterKey(), m.Channel)
	}, nil)
}

func (m *manager) Run() {
	m.shard = make([]*Shard, m.shardCount, m.shardCount)
	var i int64
	for i = 0; i < m.shardCount; i++ {
		m.shard[i] = &Shard{
			m:     make(map[int64]Conn),
			mutex: sync.RWMutex{},
		}
	}
	go m.Ping()
	m.Do(func() {
		go m.registerChannel()
	}, nil)
}

// Destroy
// 释放相关资源
func (m *manager) Destroy() {
	m.Do(func() {
		m.unRegisterChannel()
	}, nil)
}
