package websocket

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"time"
	"ws/app/auth"
	"ws/app/databases"
	"ws/app/models"
	"ws/app/mq"
	"ws/configs"
)
// conn管理相关方法
type ConnContainer interface {
	AddConn(conn Conn)
	GetConn(user auth.User) (Conn, bool)
	handleRepeatLogin(user auth.User, remote bool)
	GetAllConn(gid int64) []Conn
	GetOnlineTotal(gid int64) int64
	ConnExist(user auth.User) bool
	Register(connect Conn)
	Unregister(connect Conn)
	RemoveConn(user auth.User)
	IsOnline(user auth.User) bool
	GetOnlineUserIds(gid int64) []int64
}
// channel相关方法
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
	Conn *Client
	Action *Action
}
type Shard struct {
	m map[int64]Conn
	mutex sync.RWMutex
}
func (s *Shard) GetAll() []Conn  {
	s.mutex.RLock()
	defer  s.mutex.RUnlock()
	conns := make([]Conn, 0,len(s.m))
	for _, conn := range s.m {
		conns = append(conns, conn)
	}
	return conns
}
func (s *Shard) GetTotalCount() int64  {
	s.mutex.RUnlock()
	defer s.mutex.RUnlock()
	return int64(len(s.m))
}
func (s *Shard) Get(uid int64)  (conn Conn, exist bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	conn, exist = s.m[uid]
	return
}
func (s *Shard) Set(conn Conn)  {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[conn.GetUserId()] = conn
}
func (s *Shard) Remove(uid int64)  {
	s.mutex.Lock()
	defer  s.mutex.Unlock()
	delete(s.m, uid)
}

type manager struct {
	groupCount int64 // 分组数量
	groups []*Shard //
	Channel string // 当前manager channel
	ConnMessages chan *ConnMessage // 接受从conn所读取消息的chan
	onRegister ManagerHook //客户端连接成功hook
	onUnRegister ManagerHook //客户端连接断开hook
	userChannelCacheKey string // 客户端所在机器的channel cache key
	groupCacheKey string // 保存所在同类型manager群组SortSet keepalive key
	connGroupKeepAliveKey string // 对应conn的keepalive key
}


func (m *manager) getSpread(gid int64) *Shard  {
	return m.groups[gid]
}

func (m *manager) isCluster() bool {
	return configs.App.Cluster
}
// 发布消息
func (m *manager) publish(channel string, payload *mq.Payload) error {
	err := mq.Mq().Publish(channel, payload)
	return err
}
//
func (m *manager) publishToAllChannel(payload *mq.Payload)  {
	channels := m.getAllChannel()
	for _, channel := range channels {
		_ = m.publish(channel, payload)
	}
}

// 获取用户channel cache key
func (m *manager) getUserChannelKey(uid int64) string {
	return fmt.Sprintf(m.userChannelCacheKey, uid)
}

// 设置用户所在channel为当前manager
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
func (m *manager) handleRepeatLogin(user auth.User, remote bool) {
	s := m.getSpread(user.GetGroupId())
	old , exist := s.Get(user.GetPrimaryKey())
	if exist {
		m.SendAction(NewMoreThanOne(), old)
	} else if !remote {
		if m.isCluster() {
			oldChannel := m.getUserChannel(user.GetPrimaryKey())
			if oldChannel != ""  && oldChannel != m.GetSubscribeChannel() {
				_ = m.publish(oldChannel, &mq.Payload{
					Types: mq.TypeOtherLogin,
					Data: user.GetPrimaryKey(),
				})
			}
		}
	}
}
// 获取groupId对应的在线userIds
func (m *manager) GetOnlineUserIds(gid int64) []int64 {
	if m.isCluster() {
		key := fmt.Sprintf(m.connGroupKeepAliveKey, gid)
		source := time.Now().Unix() - 60
		ctx := context.Background()
		cmd := databases.Redis.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    strconv.FormatInt(source, 10),
			Max:    "+inf",
			Offset: 0,
			Count:  0,
		})
		idsStr := cmd.Val()
		ids := make([]int64, 0, len(idsStr))
		for _, s := range idsStr {
			id, ok  := strconv.ParseInt(s, 10, 64)
			if ok == nil {
				ids = append(ids, id)
			}
		}
		return ids
	} else {
		s := m.getSpread(gid)
		allConn := s.GetAll()
		ids := make([]int64, 0, len(allConn))
		for _, c := range allConn {
			ids = append(ids, c.GetUserId())
		}
		return ids
	}
}
// 获取groupId对应在线客户端数量
func (m *manager) GetOnlineTotal(gid int64) int64  {
	if m.isCluster() {
		key := fmt.Sprintf(m.connGroupKeepAliveKey, gid)
		source := time.Now().Unix() - 60
		ctx := context.Background()
		cmd := databases.Redis.ZCount(ctx, key, strconv.FormatInt(source, 10), "+inf")
		return cmd.Val()
	}
	s := m.getSpread(gid)
	return s.GetTotalCount()
}

// 给客户端发送消息
func (m *manager) SendAction(a  *Action, clients ...Conn) {
	for _,c := range clients {
		c.Deliver(a)
	}
}
// 用户是否在线
func (m *manager) IsOnline(user auth.User) bool  {
	if m.isCluster() {
		ctx := context.Background()
		key := fmt.Sprintf(m.connGroupKeepAliveKey, user.GetGroupId())
		lastTime := time.Now().Unix() - 60
		cmd := databases.Redis.ZScore(ctx, key, strconv.FormatInt(user.GetPrimaryKey(), 10))
		source := cmd.Val()
		return source > float64(lastTime)
	}else {
		return m.ConnExist(user)
	}
}
// 用户客户端是否存在
func (m *manager) ConnExist(user auth.User) bool {
	_, exist := m.GetConn(user)
	return exist
}
// 获取客户端
func (m *manager) GetConn(user auth.User) (client Conn,ok bool) {
	s := m.getSpread(user.GetGroupId())
	client, ok = s.Get(user.GetPrimaryKey())
	return
}
// 添加客户端
func (m *manager) AddConn(conn Conn) {
	s := m.getSpread(conn.GetGroupId())
	s.Set(conn)
}
// 移除客户端
func (m *manager) RemoveConn(user auth.User) {
	s := m.getSpread(user.GetGroupId())
	s.Remove(user.GetPrimaryKey())
}
// 获取所有客户端
func (m *manager) GetAllConn(groupId int64) (conns []Conn){
	s := m.getSpread(groupId)
	return s.GetAll()
}
// 客户端注销
func (m *manager) Unregister(conn Conn) {
	s := m.getSpread(conn.GetGroupId())
	existConn, exist := s.Get(conn.GetUserId())
	if exist {
		if existConn == conn {
			m.removeUserChannel(conn.GetUserId())
			m.RemoveConn(conn.GetUser())
			if m.onUnRegister != nil {
				m.onUnRegister(conn)
			}
		}
	}
}
// 客户端注册
// 先处理是否重复连接
// 集群模式下，如果不在本机则投递一个消息
func (m *manager) Register(conn Conn) {
	timer := time.After(1 * time.Second)
	m.handleRepeatLogin(conn.GetUser(), false)
	m.AddConn(conn)
	m.setUserChannel(conn.GetUserId())
	conn.run()
	<-timer
	if m.onRegister != nil {
		m.onRegister(conn)
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
			ping := NewPing()
			for _, s := range m.groups {
				conns := s.GetAll()
				m.SendAction(ping, conns...)
			}
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
	m.groups = make([]*Shard, m.groupCount, m.groupCount)
	var i int64
	for i= 0; i< m.groupCount; i++ {
		m.groups[i] = &Shard{
			m:     make(map[int64]Conn),
			mutex: sync.RWMutex{},
		}
	}
	go m.Ping()
	if m.isCluster() {
		go m.registerChannel()
	}
}

// 释放相关资源
func (m *manager) Destroy()  {
	if m.isCluster() {
		m.unRegisterChannel()
		for _, s:= range m.groups {
			conns := s.GetAll()
			for _, conn := range conns {
				m.removeUserChannel(conn.GetUserId())
			}
		}
	}
}
