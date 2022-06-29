package websocket

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"sync"
	"time"
	"ws/app/contract"
	"ws/app/databases"
	"ws/app/models"
	"ws/app/mq"
	rpcClient "ws/app/rpc/client"
)

// ConnContainer 管理相关方法
type ConnContainer interface {
	AddConn(conn Conn)
	GetConn(user contract.User) (Conn, bool)
	NoticeRepeatConnect(user contract.User, oldUuid string)
	GetAllConn(gid int64) []Conn
	GetOnlineTotal(gid int64) int64
	ConnExist(user contract.User) bool
	Register(connect Conn)
	Unregister(connect Conn)
	RemoveConn(user contract.User)
	IsOnline(user contract.User) bool
	IsLocalOnline(user contract.User) bool
	GetOnlineUserIds(gid int64) []int64
	Do(c func(), f func())
}

type ServiceManager interface {
	getUserServerKey(uid int64) string
	setUserServer(uid int64)
	removeUserServer(uid int64)
	getUserServer(uid int64) string
	getServer() string
	isCluster() bool
}

type ConnManager interface {
	ConnContainer
	ServiceManager
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
	shard        []*Shard          // 分组切片
	ipAddr       string            // 当前manager channel name
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

// 发布消息
func (m *manager) publish(channel string, payload *mq.Payload) error {
	err := mq.Publish(channel, payload)
	return err
}

func (m *manager) getServer() string {
	return m.ipAddr
}

// 获取用户server cache key
func (m *manager) getUserServerKey(uid int64) string {
	return fmt.Sprintf("%s:%d:service", m.GetTypes(), uid)
}

// 设置用户所在server为当前server
func (m *manager) setUserServer(uid int64) {
	m.Do(func() {
		ctx := context.Background()
		key := m.getUserServerKey(uid)
		databases.Redis.Set(ctx, key, m.getServer(), time.Hour*24*2)
	}, nil)
}

// 移除用户所在channel
func (m *manager) removeUserServer(uid int64) {
	m.Do(func() {
		ctx := context.Background()
		key := m.getUserServerKey(uid)
		databases.Redis.Del(ctx, key)
	}, nil)
}

// 获取用户channel
func (m *manager) getUserServer(uid int64) string {
	ctx := context.Background()
	key := m.getUserServerKey(uid)
	cmd := databases.Redis.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return ""
	}
	return cmd.Val()
}

// ReceiveMessage 接受消息
func (m *manager) ReceiveMessage(cm *ConnMessage) {
	m.ConnMessages <- cm
}

// NoticeRepeatConnect 重复链接
func (m *manager) NoticeRepeatConnect(user contract.User, newUuid string) {
	m.Do(func() {
		oldServer := m.getUserServer(user.GetPrimaryKey())
		if oldServer != "" {
			rpcClient.NoticeRepeatConnect(user.GetPrimaryKey(), m.GetTypes(), newUuid, oldServer)
		}
	}, func() {
		m.NoticeLocalRepeatConnect(user, newUuid)
	})
}

func (m *manager) NoticeLocalRepeatConnect(user contract.User, newUuid string) {
	oldConn, ok := m.GetConn(user)
	if ok && oldConn.GetUuid() != newUuid {
		m.SendAction(NewMoreThanOne(), oldConn)
	}
}

// GetOnlineUserIds 获取groupId对应的在线userIds
func (m *manager) GetOnlineUserIds(gid int64) []int64 {
	if m.isCluster() {
		return rpcClient.ConnectionIds(gid, m.types)
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
		return rpcClient.ConnectionTotal(gid, m.types)
	}
	return m.GetLocalOnlineTotal(gid)
}

// IsOnline 用户是否在线
func (m *manager) IsOnline(user contract.User) bool {
	if m.isCluster() {
		server := m.getUserServer(user.GetPrimaryKey())
		if server != "" {
			return rpcClient.ConnectionOnline(user.GetPrimaryKey(), m.GetTypes(), server)
		}
		return false
	} else {
		return m.IsLocalOnline(user)
	}
}

func (m *manager) IsLocalOnline(user contract.User) bool {
	return m.ConnExist(user)
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
	m.NoticeRepeatConnect(conn.GetUser(), conn.GetUuid())
	m.AddConn(conn)
	m.setUserServer(conn.GetUserId())
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
}

// Destroy
// 释放相关资源
func (m *manager) Destroy() {
	m.Do(func() {
	}, nil)
}
