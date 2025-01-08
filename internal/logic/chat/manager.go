package chat

import (
	"context"
	"fmt"
	"gf-chat/api"
	v1 "gf-chat/api/chat/v1"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const eventRegister = "register"
const eventUnRegister = "unregister"
const eventMessage = "message"

const userServer = "%s:user:%d:server"

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
	getOnlineUserIds(ctx context.Context, gid uint, forceLocal ...bool) ([]uint, error)
	setUserServer(ctx context.Context, uid uint, server string) error
	getUserServer(ctx context.Context, uid uint) (string, error)
	getConnInfo(ctx context.Context, customerId, uid uint, forceLocal ...bool) (bool, string)
}

type connManager interface {
	connContainer
	run()
	ping()
	SendAction(act *api.ChatAction, conn ...iWsConn)
	receiveMessage(cm *chatConnMessage)
	handleReceiveMessage()
	noticeRead(ctx context.Context, customerId uint, uid uint, msgIds []uint, forceLocal ...bool) error
}

type eventArg struct {
	conn iWsConn
	msg  *model.CustomerChatMessage
}

func newManager(shareCount uint, msgCount int, pingDuration time.Duration, cluster bool, types string) *manager {
	return &manager{
		shardCount:   shareCount,
		connMessages: make(chan *chatConnMessage, msgCount),
		events:       nil,
		pingDuration: pingDuration,
		cluster:      cluster,
		types:        types,
	}
}

type manager struct {
	shardCount   uint                     // 分组数量, 默认 10
	shard        []*shard                 // 分组切片
	connMessages chan *chatConnMessage    // 接受从conn所读取消息的chan
	events       map[string][]eventHandle // 事件
	pingDuration time.Duration            // default to 10 seconds
	cluster      bool
	types        string
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

func (m *manager) userServerKey(id uint) string {
	return fmt.Sprintf(userServer, m.types, id)
}

func (m *manager) setUserServer(ctx context.Context, uid uint, server string) error {
	var expired int64 = 60 * 60 * 24
	_, err := g.Redis().Set(ctx, m.userServerKey(uid), server, gredis.SetOption{
		TTLOption: gredis.TTLOption{
			EX: &expired,
		},
	})
	return err
}

func (m *manager) getUserServer(ctx context.Context, uid uint) (string, error) {
	val, err := g.Redis().Get(ctx, m.userServerKey(uid))
	if err != nil {
		return "", err
	}
	return val.String(), nil
}
func (m *manager) removeUserServer(ctx context.Context, uid uint) error {
	_, err := g.Redis().Del(ctx, m.userServerKey(uid))
	return err
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

func (m *manager) isCallLocal(forceLocal ...bool) bool {
	local := false
	if len(forceLocal) > 0 {
		local = forceLocal[0]
	}
	return local || !m.cluster
}

func (m *manager) getOnlineUserIds(ctx context.Context, customerId uint, forceLocal ...bool) ([]uint, error) {
	if m.isCallLocal(forceLocal...) {
		s := m.getSpread(customerId)
		allConn := s.getAll()
		ids := make([]uint, 0)
		for _, conn := range allConn {
			if conn.getCustomerId() == customerId {
				ids = append(ids, conn.getUserId())
			}
		}
		return ids, nil
	}
	idArr := garray.NewIntArray(true)
	err := service.Grpc().CallAll(ctx, func(client v1.ChatClient) {
		r, err := client.GetOnlineUserIds(ctx, &v1.GetOnlineUserIdsRequest{
			CustomerId: uint32(customerId),
			Type:       m.types,
		})
		if err != nil {
			g.Log().Errorf(ctx, "%+v", err)
		} else {
			idArr.Append(gconv.Ints(r.Uid)...)
		}
	})
	if err != nil {
		return nil, err
	}
	return gconv.Uints(idArr.Slice()), nil
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

func (m *manager) isLocalOnline(customerId uint, uid uint) bool {
	return m.connExist(customerId, uid)
}

// SendAction 给客户端发送消息
func (m *manager) SendAction(a *api.ChatAction, clients ...iWsConn) {
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

func (m *manager) getConnInfo(ctx context.Context, customerId, uid uint, forceLocal ...bool) (bool, string) {
	if m.isCallLocal(forceLocal...) {
		conn, exist := m.getConn(customerId, uid)
		if exist {
			return true, conn.getPlatform()
		} else {
			return false, ""
		}
	}
	server, _ := m.getUserServer(ctx, uid)
	if server != "" {
		r, err := service.Grpc().Client(server).GetConnInfo(ctx, &v1.GetConnInfoRequest{
			UserId:     uint32(uid),
			CustomerId: uint32(customerId),
			Type:       m.types,
		})
		if err == nil {
			return r.Exist, r.Platform
		}
		g.Log().Errorf(ctx, "%+v", err)
	}
	return false, ""
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
			if m.cluster {
				err := m.removeUserServer(ctx, conn.getUserId())
				if err != nil {
					g.Log().Errorf(ctx, "%+v", err)
				}
			}
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
	m.addConn(client)
	m.noticeRepeatConnect(client.getUser(), client.getUuid())
	if m.cluster {
		err := m.setUserServer(ctx, user.getPrimaryKey(), service.Grpc().GetServerName())
		if err != nil {
			return err
		}
	}
	client.run()
	<-timer
	err := m.trigger(ctx, eventRegister, eventArg{
		conn: client,
	})
	return err
}
func (m *manager) noticeRead(ctx context.Context, customerId, uid uint, msgIds []uint, forceLocal ...bool) (err error) {
	if m.isCallLocal(forceLocal...) {
		conn, exist := m.getConn(customerId, uid)
		if exist {
			act := action.newReadAction(msgIds)
			m.SendAction(act, conn)
		}
		return nil
	}
	server, err := m.getUserServer(ctx, uid)
	if err != nil {
		return err
	}
	if server != "" {
		_, err = service.Grpc().Client(server).NoticeRead(ctx, &v1.NoticeReadRequest{
			CustomerId: uint32(customerId),
			UserId:     uint32(uid),
			MsgId:      gconv.Uint32s(msgIds),
			Type:       m.types,
		})
		if err != nil {
			return
		}
	}
	return nil
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
