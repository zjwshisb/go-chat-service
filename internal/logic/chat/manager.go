package chat

import (
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"sync"
	"time"

	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

type MessageHandle interface {
	handleReceiveMessage()
	handleMessage(cm *chatConnMessage)
	handleOffline(msg *model.CustomerChatMessage)
	DeliveryMessage(msg *model.CustomerChatMessage)
}

type ManagerHook = func(conn iWsConn)

type Shard struct {
	m     map[uint]iWsConn
	mutex sync.RWMutex
}

func (s *Shard) getAll() []iWsConn {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	conns := make([]iWsConn, 0, len(s.m))
	for _, conn := range s.m {
		conns = append(conns, conn)
	}
	return conns
}

func (s *Shard) getTotalCount() uint {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return uint(len(s.m))
}

func (s *Shard) get(uid uint) (conn iWsConn, exist bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	conn, exist = s.m[uid]
	return
}
func (s *Shard) set(conn iWsConn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[conn.GetUserId()] = conn
}
func (s *Shard) remove(uid uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.m, uid)
}

type manager struct {
	ShardCount   uint                  // 分组数量
	shard        []*Shard              // 分组切片
	ConnMessages chan *chatConnMessage // 接受从conn所读取消息的chan
	OnRegister   ManagerHook           //conn连接成功hook
	OnUnRegister ManagerHook           //conn连接断开hook
	Types        string                //类型
}

func (m *manager) GetTypes() string {
	return m.Types
}

func (m *manager) getMod(customerId uint) uint {
	return customerId % m.ShardCount
}

func (m *manager) getSpread(customerId uint) *Shard {
	return m.shard[m.getMod(customerId)]
}

// ReceiveMessage 接受消息
func (m *manager) ReceiveMessage(cm *chatConnMessage) {
	m.ConnMessages <- cm
}

// NoticeRepeatConnect 重复链接
func (m *manager) NoticeRepeatConnect(user IChatUser, newUuid string) {
	m.NoticeLocalRepeatConnect(user, newUuid)
}

func (m *manager) NoticeLocalRepeatConnect(user IChatUser, newUuid string) {
	oldConn, ok := m.GetConn(user.GetCustomerId(), user.GetPrimaryKey())
	if ok && oldConn.GetUuid() != newUuid {
		m.SendAction(service.Action().NewMoreThanOne(), oldConn)
	}
}

// GetOnlineUserIds 获取groupId对应的在线userIds
func (m *manager) GetOnlineUserIds(gid uint) []uint {
	return m.GetLocalOnlineUserIds(gid)
}

func (m *manager) GetLocalOnlineUserIds(gid uint) []uint {
	s := m.getSpread(gid)
	allConn := s.getAll()
	ids := make([]uint, 0)
	for _, conn := range allConn {
		if conn.GetCustomerId() == gid {
			ids = append(ids, conn.GetUserId())
		}
	}
	return ids
}

// GetLocalOnlineTotal 获取本地groupId对应在线客户端数量
func (m *manager) GetLocalOnlineTotal(customerId uint) uint {
	s := m.getSpread(customerId)
	return s.getTotalCount()
}

// GetOnlineTotal 获取groupId对应在线客户端数量
func (m *manager) GetOnlineTotal(customerId uint) uint {
	return m.GetLocalOnlineTotal(customerId)
}

// IsOnline 用户是否在线
func (m *manager) IsOnline(customerId uint, uid uint) bool {
	return m.IsLocalOnline(customerId, uid)
}

func (m *manager) IsLocalOnline(customerId uint, uid uint) bool {
	return m.ConnExist(customerId, uid)
}

// SendAction 给客户端发送消息
func (m *manager) SendAction(a *api.ChatAction, clients ...iWsConn) {
	for _, c := range clients {
		c.Deliver(a)
	}
}

// ConnExist 连接是否存在
func (m *manager) ConnExist(customerId uint, uid uint) bool {
	_, exist := m.GetConn(customerId, uid)
	return exist
}

// GetConn 获取客户端
func (m *manager) GetConn(customerId, uid uint) (client iWsConn, ok bool) {
	s := m.getSpread(customerId)
	client, ok = s.get(uid)
	return
}

// AddConn 添加客户端
func (m *manager) AddConn(conn iWsConn) {
	s := m.getSpread(conn.GetCustomerId())
	s.set(conn)
}

// RemoveConn 移除客户端
func (m *manager) RemoveConn(user IChatUser) {
	s := m.getSpread(user.GetCustomerId())
	s.remove(user.GetPrimaryKey())
}

// GetAllConn 获取所有客户端
func (m *manager) GetAllConn(customerId uint) (conns []iWsConn) {
	s := m.getSpread(customerId)
	return s.getAll()
}

func (m *manager) GetAllConnCount() uint {
	var count uint
	for gid := range m.shard {
		count += m.GetLocalOnlineTotal(uint(gid))
	}
	return count
}

func (m *manager) GetTotalConn() []iWsConn {
	conns := make([]iWsConn, 0)
	for gid := range m.shard {
		conns = append(conns, m.GetAllConn(uint(gid))...)
	}
	return conns
}

// Unregister 客户端注销
func (m *manager) Unregister(conn iWsConn) {
	s := m.getSpread(conn.GetCustomerId())
	existConn, exist := s.get(conn.GetUserId())
	if exist {
		if existConn == conn {
			m.RemoveConn(conn.GetUser())
			if m.OnUnRegister != nil {
				m.OnUnRegister(conn)
			}
		}
	}
}

// Register 客户端注册
// 先处理是否重复连接
// 集群模式下，如果不在本机则投递一个消息
func (m *manager) Register(conn *websocket.Conn, user IChatUser, platform string) {
	client := &client{
		Conn:        conn,
		CloseSignal: make(chan interface{}),
		Send:        make(chan *api.ChatAction, 100),
		Once:        sync.Once{},
		Manager:     m,
		User:        user,
		Uuid:        guid.S(),
		Created:     time.Now().Unix(),
		Limiter:     rate.NewLimiter(5, 10),
		Platform:    platform,
	}
	timer := time.After(1 * time.Second)
	m.NoticeRepeatConnect(client.GetUser(), client.GetUuid())
	m.AddConn(client)
	client.Run()
	<-timer
	if m.OnRegister != nil {
		m.OnRegister(client)
	}
}
func (m *manager) NoticeRead(customerId, adminId uint, msgIds []uint) {
	conn, exist := m.GetConn(customerId, adminId)
	if exist {
		action := service.Action().NewReadAction(msgIds)
		m.SendAction(action, conn)
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
			ping := service.Action().NewPing()
			for _, s := range m.shard {
				conns := s.getAll()
				m.SendAction(ping, conns...)
			}
		}
	}
}

func (m *manager) Run() {
	m.shard = make([]*Shard, m.ShardCount)
	var i uint
	for i = 0; i < m.ShardCount; i++ {
		m.shard[i] = &Shard{
			m:     make(map[uint]iWsConn),
			mutex: sync.RWMutex{},
		}
	}
	go m.Ping()
}

// Destroy
// 释放相关资源
func (m *manager) Destroy() {
}
