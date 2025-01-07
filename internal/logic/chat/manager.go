package chat

import (
	"context"
	"gf-chat/api/v1"
	"gf-chat/internal/model"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
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
	getConn(customerId uint, uid uint) (iWsConn, bool)
	noticeRepeatConnect(user iChatUser, newUid string)
	getAllConn(customerId uint) []iWsConn
	getOnlineTotal(customerId uint) uint
	connExist(customerId uint, uid uint) bool
	register(ctx context.Context, conn *websocket.Conn, user iChatUser, platform string) error
	unregister(connect iWsConn)
	removeConn(user iChatUser)
	isOnline(customerId uint, uid uint) bool
	isLocalOnline(customerId uint, uid uint) bool
	getOnlineUserIds(gid uint) []uint
}

type connManager interface {
	connContainer
	run()
	ping()
	SendAction(act *v1.ChatAction, conn ...iWsConn)
	receiveMessage(cm *chatConnMessage)
	handleReceiveMessage()
	noticeRead(customerId uint, uid uint, msgIds []uint)
}

type eventArg struct {
	conn iWsConn
	msg  *model.CustomerChatMessage
}

type manager struct {
	shardCount   uint                     // 分组数量, 默认 10
	shard        []*shard                 // 分组切片
	connMessages chan *chatConnMessage    // 接受从conn所读取消息的chan
	events       map[string][]eventHandle // 事件
	pingDuration time.Duration            // default to 10 seconds
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
func (m *manager) noticeRepeatConnect(user iChatUser, newUuid string) {
	m.NoticeLocalRepeatConnect(user, newUuid)
}

func (m *manager) NoticeLocalRepeatConnect(user iChatUser, newUuid string) {
	oldConn, ok := m.getConn(user.getCustomerId(), user.getPrimaryKey())
	if ok && oldConn.getUuid() != newUuid {
		m.SendAction(action.newMoreThanOne(), oldConn)
	}
}

// GetOnlineUserIds 获取groupId对应的在线userIds
func (m *manager) getOnlineUserIds(gid uint) []uint {
	return m.GetLocalOnlineUserIds(gid)
}

func (m *manager) GetLocalOnlineUserIds(gid uint) []uint {
	s := m.getSpread(gid)
	allConn := s.getAll()
	ids := make([]uint, 0)
	for _, conn := range allConn {
		if conn.getCustomerId() == gid {
			ids = append(ids, conn.getUserId())
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
func (m *manager) getOnlineTotal(customerId uint) uint {
	return m.GetLocalOnlineTotal(customerId)
}

// IsOnline 用户是否在线
func (m *manager) isOnline(customerId uint, uid uint) bool {
	return m.isLocalOnline(customerId, uid)
}

func (m *manager) isLocalOnline(customerId uint, uid uint) bool {
	return m.connExist(customerId, uid)
}

// SendAction 给客户端发送消息
func (m *manager) SendAction(a *v1.ChatAction, clients ...iWsConn) {
	for _, c := range clients {
		c.deliver(a)
	}
}

// ConnExist 连接是否存在
func (m *manager) connExist(customerId uint, uid uint) bool {
	_, exist := m.getConn(customerId, uid)
	return exist
}

// GetConn 获取客户端
func (m *manager) getConn(customerId, uid uint) (iWsConn, bool) {
	s := m.getSpread(customerId)
	return s.get(uid)
}

// AddConn 添加客户端
func (m *manager) addConn(conn iWsConn) {
	s := m.getSpread(conn.getCustomerId())
	s.set(conn)
}

// RemoveConn 移除客户端
func (m *manager) removeConn(user iChatUser) {
	s := m.getSpread(user.getCustomerId())
	s.remove(user.getPrimaryKey())
}

// GetAllConn 获取所有客户端
func (m *manager) getAllConn(customerId uint) (conns []iWsConn) {
	s := m.getSpread(customerId)
	conns = slice.Filter(s.getAll(), func(index int, item iWsConn) bool {
		return item.getCustomerId() == customerId
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
		conns = append(conns, m.getAllConn(uint(gid))...)
	}
	return conns
}

// Unregister 客户端注销
func (m *manager) unregister(conn iWsConn) {
	ctx := gctx.New()
	existConn, exist := m.getConn(conn.getCustomerId(), conn.getUserId())
	if exist {
		if existConn == conn {
			m.removeConn(conn.getUser())
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
func (m *manager) register(ctx context.Context, conn *websocket.Conn, user iChatUser, platform string) error {
	client := newClient(conn, user, platform)
	client.manager = m
	timer := time.After(1 * time.Second)
	m.noticeRepeatConnect(client.getUser(), client.getUuid())
	m.addConn(client)
	client.run()
	<-timer
	err := m.trigger(ctx, eventRegister, eventArg{
		conn: client,
	})
	return err
}
func (m *manager) noticeRead(customerId, adminId uint, msgIds []uint) {
	conn, exist := m.getConn(customerId, adminId)
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
	duration := m.pingDuration
	if duration == 0 {
		duration = time.Second * 60
	}
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			ping := action.newPing()
			for _, s := range m.shard {
				conns := s.getAll()
				for _, conn := range conns {
					// 如果连接超过60分钟没有活动，则关闭连接
					duration := gtime.Now().Second() - conn.getLastActive().Second()
					if duration > 60*60 {
						conn.close()
					} else {
						m.SendAction(ping, conn)
					}
				}
			}
		}
	}
}

func (m *manager) run() {
	count := m.shardCount
	if count == 0 {
		count = 10
	}
	m.shard = make([]*shard, count)
	var i uint
	for i = 0; i < count; i++ {
		m.shard[i] = &shard{
			m:     make(map[uint]iWsConn),
			mutex: &sync.RWMutex{},
		}
	}
	go m.handleReceiveMessage()
	go m.ping()
}
