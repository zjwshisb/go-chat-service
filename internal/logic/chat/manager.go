package chat

import (
	"context"
	"fmt"
	"gf-chat/api"
	v1 "gf-chat/api/chat/v1"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

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
	noticeRepeatConnect(ctx context.Context, uid, customerId uint, newUuid string, forceLocal ...bool) error
	getAllConn(customerId uint) []iWsConn
	register(ctx context.Context, conn *websocket.Conn, user iChatUser, platform string) error
	unregister(connect iWsConn)
	removeConn(user iChatUser)
	getOnlineUserIds(ctx context.Context, gid uint, forceLocal ...bool) ([]uint, error)
	setUserServer(ctx context.Context, uid uint, server string) error
	getUserServer(ctx context.Context, uid uint) (string, error)
	removeUserServer(ctx context.Context, uid uint) error
	getConnInfo(ctx context.Context, customerId, uid uint, forceLocal ...bool) (bool, string, error)
}

type connManager interface {
	connContainer
	run()
	ping()
	SendAction(act *api.ChatAction, conn ...iWsConn)
	handleMessage(ctx context.Context, conn iWsConn, msg *model.CustomerChatMessage)
	noticeRead(ctx context.Context, customerId uint, uid uint, msgIds []uint, forceLocal ...bool) error
}

type eventArg struct {
	conn iWsConn
	msg  *model.CustomerChatMessage
}

func newManager(shareCount uint, pingDuration time.Duration, cluster bool, types string) *manager {
	return &manager{
		shardCount:   shareCount,
		events:       nil,
		pingDuration: pingDuration,
		cluster:      cluster,
		types:        types,
	}
}

type manager struct {
	shardCount   uint                     // 分组数量, 默认 10
	shard        []*shard                 // 分组切片
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

// 保存用户websocket连接所在服务的redis key,
func (m *manager) userServerKey(id uint) string {
	return fmt.Sprintf(userServer, m.types, id)
}

// 设置用户websocket所在服务名称
func (m *manager) setUserServer(ctx context.Context, uid uint, server string) error {
	var expired int64 = 60 * 60 * 24
	_, err := g.Redis().Set(ctx, m.userServerKey(uid), server, gredis.SetOption{
		TTLOption: gredis.TTLOption{
			EX: &expired,
		},
	})
	return err
}

// 获取用户websocket所在服务名称
func (m *manager) getUserServer(ctx context.Context, uid uint) (string, error) {
	val, err := g.Redis().Get(ctx, m.userServerKey(uid))
	if err != nil {
		return "", err
	}
	if val.IsNil() {
		return "", err
	}
	return val.String(), nil
}

// 移除用户websocket所在服务名称
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

// handleMessage handles incoming messages
// Triggers the message event for the given connection and message
func (m *manager) handleMessage(ctx context.Context, conn iWsConn, msg *model.CustomerChatMessage) {
	err := m.trigger(ctx, eventMessage, eventArg{
		conn: conn,
		msg:  msg,
	})
	if err != nil {
		log.Errorf(ctx, "%+v", err)
	}
}

// noticeRepeatConnect checks if a user is already connected and sends a notification if they are
// connecting from multiple locations. It handles both local and remote (cluster) connections.
//
// Parameters:
//   - ctx: Context for the operation
//   - uid: User ID to check for repeat connections
//   - customerId: Customer ID the user belongs to
//   - newUuid: UUID of the new connection attempt
//   - forceLocal: Optional bool to force local-only check, ignoring cluster
//
// Returns error if the check fails
func (m *manager) noticeRepeatConnect(ctx context.Context, uid, customerId uint, newUuid string, forceLocal ...bool) error {
	userLocal, server, err := m.isUserLocal(ctx, uid)
	if err != nil {
		return err
	}
	if m.isCallLocal(forceLocal...) || userLocal {
		oldConn, ok := m.getConn(customerId, uid)
		if ok && oldConn.getUuid() != newUuid {
			m.SendAction(action.newMoreThanOne(), oldConn)
		}
		return nil
	} else if server != "" {
		rpcClient := service.Grpc().Client(ctx, server)
		if rpcClient != nil {
			_, err = rpcClient.NoticeRepeatConnect(ctx, &v1.NoticeRepeatConnectRequest{
				UserId:     uint32(uid),
				CustomerId: uint32(customerId),
				NewUid:     newUuid,
				Type:       m.types,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// isUserLocal checks if a user is connected locally or on a remote server
// Returns:
// - bool: true if user is local, false otherwise
// - string: server name if user is remote, empty string otherwise
// - error: error if retrieval fails
func (m *manager) isUserLocal(ctx context.Context, id uint) (bool, string, error) {
	if !m.cluster {
		return true, "", nil
	}
	server, err := m.getUserServer(ctx, id)
	if err != nil {
		return false, "", err
	}
	return server == service.Grpc().GetServerName(), server, nil
}

// isCallLocal checks if the local flag is set or if the cluster is disabled
// Returns:
// - bool: true if local flag is set or cluster is disabled, false otherwise
func (m *manager) isCallLocal(forceLocal ...bool) bool {
	local := false
	if len(forceLocal) > 0 {
		local = forceLocal[0]
	}
	return local || !m.cluster
}

// getOnlineUserIds retrieves user IDs of online users for a given customer ID
// Returns error if retrieval fails
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
	} else {
		idArr := garray.NewIntArray(true)
		err := service.Grpc().CallAll(ctx, func(client v1.ChatClient) {
			r, err := client.GetOnlineUserIds(ctx, &v1.GetOnlineUserIdsRequest{
				CustomerId: uint32(customerId),
				Type:       m.types,
			})
			if err != nil {
				log.Errorf(ctx, "%+v", err)
			} else {
				idArr.Append(gconv.Ints(r.Uid)...)
			}
		})
		if err != nil {
			return nil, err
		}
		return gconv.Uints(idArr.Slice()), nil
	}
}

// GetLocalOnlineTotal 获取本地groupId对应在线客户端数量
func (m *manager) GetLocalOnlineTotal(customerId uint) uint {
	s := m.getSpread(customerId)
	return s.getTotalCount()
}

// SendAction 给客户端发送消息
func (m *manager) SendAction(a *api.ChatAction, clients ...iWsConn) {
	for _, c := range clients {
		c.deliver(a)
	}
}

// GetConn 获取客户端
func (m *manager) getConn(customerId, uid uint) (iWsConn, bool) {
	s := m.getSpread(customerId)
	return s.get(uid)
}

// getConnInfo retrieves connection information for a given customer ID and user ID
// Returns:
// - bool: true if connection exists, false otherwise
// - string: platform of the connection
// - error: error if retrieval fails
func (m *manager) getConnInfo(ctx context.Context, customerId, uid uint, forceLocal ...bool) (bool, string, error) {
	userLocal, server, _ := m.isUserLocal(ctx, uid)
	if m.isCallLocal(forceLocal...) || userLocal {
		conn, exist := m.getConn(customerId, uid)
		if exist {
			return true, conn.getPlatform(), nil
		} else {
			return false, "", nil
		}
	}
	if server != "" {
		rpcClient := service.Grpc().Client(ctx, server)
		if rpcClient != nil {
			r, err := rpcClient.GetConnInfo(ctx, &v1.GetConnInfoRequest{
				UserId:     uint32(uid),
				CustomerId: uint32(customerId),
				Type:       m.types,
			})
			if err != nil {
				return false, "", err
			}
			return r.Exist, r.Platform, nil
		}
	}
	return false, "", nil
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

// unregister unregisters a client from the manager
// Removes the client from the connection map and triggers the unregister event
func (m *manager) unregister(conn iWsConn) {
	ctx := gctx.New()
	existConn, exist := m.getConn(conn.getCustomerId(), conn.getUserId())
	if exist {
		if existConn == conn {
			m.removeConn(conn.getUser())
			if m.cluster {
				err := m.removeUserServer(ctx, conn.getUserId())
				if err != nil {
					log.Errorf(ctx, "%+v", err)
				}
			}
			err := m.trigger(ctx, eventUnRegister, eventArg{
				conn: conn,
			})
			if err != nil {
				log.Errorf(ctx, "%+v", err)
			}
		}
	}
}

// register registers a client with the manager
// Adds the client to the connection map and triggers the register event
func (m *manager) register(ctx context.Context, conn *websocket.Conn, user iChatUser, platform string) error {
	client := newClient(conn, user, platform)
	client.manager = m
	timer := time.After(1 * time.Second)
	err := m.noticeRepeatConnect(ctx, user.getPrimaryKey(), user.getCustomerId(), client.getUuid())
	if err != nil {
		return err
	}
	m.addConn(client)
	if m.cluster {
		err := m.setUserServer(ctx, user.getPrimaryKey(), service.Grpc().GetServerName())
		if err != nil {
			return err
		}
	}
	client.run()
	<-timer
	err = m.trigger(ctx, eventRegister, eventArg{
		conn: client,
	})
	return err
}

// noticeRead notifies the read action for a given customer ID and user ID
// Returns error if notification fails
func (m *manager) noticeRead(ctx context.Context, customerId, uid uint, msgIds []uint, forceLocal ...bool) (err error) {
	userLocal, server, err := m.isUserLocal(ctx, uid)
	if err != nil {
		return err
	}
	if m.isCallLocal(forceLocal...) || userLocal {
		conn, exist := m.getConn(customerId, uid)
		if exist {
			act := action.newReadAction(msgIds)
			m.SendAction(act, conn)
		}
		return nil
	}
	if server != "" {
		rpcClient := service.Grpc().Client(ctx, server)
		if rpcClient != nil {
			_, err = rpcClient.NoticeRead(ctx, &v1.NoticeReadRequest{
				CustomerId: uint32(customerId),
				UserId:     uint32(uid),
				MsgId:      gconv.Uint32s(msgIds),
				Type:       m.types,
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// ping sends ping actions to all connections in the manager
// It runs in a separate goroutine and sends ping actions every 60 seconds
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
					// 超过60分钟没有活动，则关闭连接，节省资源
					duration := gtime.Now().Unix() - conn.getLastActive().Unix()
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

// run initializes the shard structure and starts the ping goroutine
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
	go m.ping()
}
