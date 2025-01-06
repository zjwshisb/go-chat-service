package chat

import (
	"context"
	"gf-chat/api/v1"
	"gf-chat/internal/model"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const eventRegister = "register"
const eventUnRegister = "unregister"
const eventMessage = "message"

// 返回err会停止后续事件的执行
type eventHandle = func(ctx context.Context, arg eventArg) error

type connContainer interface {
	addConn(conn iWsConn)
	GetConn(customerId uint, uid uint) (iWsConn, bool)
	NoticeRepeatConnect(user IChatUser, newUid string)
	GetAllConn(customerId uint) []iWsConn
	GetOnlineTotal(customerId uint) uint
	ConnExist(customerId uint, uid uint) bool
	Register(ctx context.Context, conn *websocket.Conn, user IChatUser, platform string) error
	Unregister(connect iWsConn)
	removeConn(user IChatUser)
	IsOnline(customerId uint, uid uint) bool
	IsLocalOnline(customerId uint, uid uint) bool
	GetOnlineUserIds(gid uint) []uint
}

type connManager interface {
	connContainer
	run()
	ping()
	SendAction(act *v1.ChatAction, conn ...iWsConn)
	receiveMessage(cm *chatConnMessage)
	handleReceiveMessage()
	GetTypes() string
	NoticeRead(customerId uint, uid uint, msgIds []uint)
}

type eventArg struct {
	conn iWsConn
	msg  *model.CustomerChatMessage
}

type manager struct {
	shardCount   uint                     // 分组数量
	shard        []*shard                 // 分组切片
	connMessages chan *chatConnMessage    // 接受从conn所读取消息的chan
	events       map[string][]eventHandle // 事件
	types        string                   //类型
}

// 注册事件
func (m *manager) on(name string, handle eventHandle) {
	if m.events == nil {
		m.events = make(map[string][]eventHandle)
	}
	m.events[name] = append(m.events[name], handle)
}

// 触发事件
func (m *manager) trigger(ctx context.Context, name string, arg eventArg) error {
	handlers, exist := m.events[name]
	if exist {
		for _, handle := range handlers {
			err := handle(ctx, arg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *manager) GetTypes() string {
	return m.types
}

func (m *manager) getMod(customerId uint) uint {
	return customerId % m.shardCount
}

func (m *manager) getSpread(customerId uint) *shard {
	return m.shard[m.getMod(customerId)]
}

// ReceiveMessage 接受消息
func (m *manager) receiveMessage(cm *chatConnMessage) {
	m.connMessages <- cm
}

// 从conn接受消息并处理
func (m *manager) handleReceiveMessage() {
	for {
		payload := <-m.connMessages
		go func() {
			ctx := gctx.New()
			err := m.trigger(ctx, eventMessage, eventArg{
				conn: payload.Conn,
				msg:  payload.Msg,
			})
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
			}
		}()
	}
}

// NoticeRepeatConnect 重复链接
func (m *manager) NoticeRepeatConnect(user IChatUser, newUuid string) {
	m.NoticeLocalRepeatConnect(user, newUuid)
}

func (m *manager) NoticeLocalRepeatConnect(user IChatUser, newUuid string) {
	oldConn, ok := m.GetConn(user.GetCustomerId(), user.GetPrimaryKey())
	if ok && oldConn.GetUuid() != newUuid {
		m.SendAction(action.newMoreThanOne(), oldConn)
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
func (m *manager) SendAction(a *v1.ChatAction, clients ...iWsConn) {
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
func (m *manager) GetConn(customerId, uid uint) (iWsConn, bool) {
	s := m.getSpread(customerId)
	return s.get(uid)
}

// AddConn 添加客户端
func (m *manager) addConn(conn iWsConn) {
	s := m.getSpread(conn.GetCustomerId())
	s.set(conn)
}

// RemoveConn 移除客户端
func (m *manager) removeConn(user IChatUser) {
	s := m.getSpread(user.GetCustomerId())
	s.remove(user.GetPrimaryKey())
}

// GetAllConn 获取所有客户端
func (m *manager) GetAllConn(customerId uint) (conns []iWsConn) {
	s := m.getSpread(customerId)
	conns = slice.Filter(s.getAll(), func(index int, item iWsConn) bool {
		return item.GetCustomerId() == customerId
	})
	return
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
	ctx := gctx.New()
	existConn, exist := m.GetConn(conn.GetCustomerId(), conn.GetUserId())
	if exist {
		if existConn == conn {
			m.removeConn(conn.GetUser())
			err := m.trigger(ctx, eventUnRegister, eventArg{
				conn: conn,
			})
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
			}
		}
	}
}

// Register 客户端注册
// 先处理是否重复连接
// 集群模式下，如果不在本机则投递一个消息
func (m *manager) Register(ctx context.Context, conn *websocket.Conn, user IChatUser, platform string) error {
	client := newClient(conn, user, platform)
	client.manager = m
	timer := time.After(1 * time.Second)
	m.NoticeRepeatConnect(client.GetUser(), client.GetUuid())
	m.addConn(client)
	client.Run()
	<-timer
	err := m.trigger(ctx, eventRegister, eventArg{
		conn: client,
	})
	return err
}
func (m *manager) NoticeRead(customerId, adminId uint, msgIds []uint) {
	conn, exist := m.GetConn(customerId, adminId)
	if exist {
		act := action.newReadAction(msgIds)
		m.SendAction(act, conn)
	}
}

// Ping
// 给所有客户端发送心跳
// 客户端因意外断开链接，服务器没有关闭事件，无法得知连接已关闭
// 通过心跳发送""字符串，如果发送失败，则调用conn的close方法执行清理
func (m *manager) ping() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			ping := action.newPing()
			for _, s := range m.shard {
				conns := s.getAll()
				m.SendAction(ping, conns...)
			}
		}
	}
}

func (m *manager) run() {
	m.shard = make([]*shard, m.shardCount)
	var i uint
	for i = 0; i < m.shardCount; i++ {
		m.shard[i] = &shard{
			m:     make(map[uint]iWsConn),
			mutex: &sync.RWMutex{},
		}
	}
	go m.handleReceiveMessage()
	go m.ping()
}
