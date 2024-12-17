package chat

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"sort"
)

const TypeAdmin = "admin"

func newAdminManager() *adminManager {
	adminM = &adminManager{
		&manager{
			shardCount:   10,
			connMessages: make(chan *chatConnMessage, 100),
			types:        TypeAdmin,
		},
	}
	adminM.onRegister = adminM.registerHook
	adminM.onUnRegister = adminM.unregisterHook
	return adminM
}

type adminManager struct {
	*manager
}

func (m *adminManager) run() {
	m.Run()
	go m.handleReceiveMessage()
}

// DeliveryMessage
// 投递消息
// 查询admin是否在线，是则直接投递
// 最后则说明admin不在线，处理离线逻辑
func (m *adminManager) deliveryMessage(ctx context.Context, msg *model.CustomerChatMessage, userConn iWsConn) error {
	adminConn, exist := m.GetConn(msg.CustomerId, msg.AdminId)
	if exist { // admin在线
		err := userM.triggerMessageEvent(ctx, consts.AutoRuleSceneAdminOnline, msg, userConn)
		if err != nil {
			return nil
		}
		adminConn.Deliver(action.newReceive(msg))
		return nil
	}
	return m.handleOffline(ctx, msg, userConn)
}

// 从管道接受消息并处理
func (m *adminManager) handleReceiveMessage() {
	for {
		payload := <-m.connMessages
		go func() {
			ctx := gctx.New()
			err := m.handleMessage(ctx, payload)
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
			}
		}()
	}
}

func (m *adminManager) sendWaiting(admin *model.CustomerAdmin, user IChatUser) {

}

func (m *adminManager) sendOffline(admin *model.CustomerAdmin, msg *model.CustomerChatMessage) {

}

// 处理离线消息
func (m *adminManager) handleOffline(ctx context.Context, msg *model.CustomerChatMessage, userConn iWsConn) error {
	err := userM.triggerMessageEvent(ctx, consts.AutoRuleSceneAdminOffline, msg, userConn)
	if err != nil {
		return err
	}
	admin, err := service.Admin().First(ctx, do.CustomerAdmins{Id: msg.AdminId})
	if err != nil {
		return err
	}
	message, err := service.ChatMessage().NewOffline(ctx, admin)
	if err != nil {
		return err
	}
	if message != nil {
		message.UserId = msg.UserId
		message.SessionId = msg.SessionId
		message, err = service.ChatMessage().Insert(ctx, message)
		if err != nil {
			return err
		}
		err = userM.DeliveryMessage(ctx, message)
		if err != nil {
			return err
		}
	}
	m.sendOffline(admin, msg)
	return nil
}

// 处理消息
func (m *adminManager) handleMessage(ctx context.Context, payload *chatConnMessage) error {
	msg := payload.Msg
	conn := payload.Conn
	if msg.UserId > 0 {
		if !service.ChatRelation().IsUserValid(ctx, conn.GetUserId(), msg.UserId) {
			conn.Deliver(action.newErrorMessage("该用户已失效，无法发送消息"))
			return gerror.New("该用户已失效，无法发送消息")
		}
		session, _ := service.ChatSession().FirstActive(ctx, msg.UserId, conn.GetUserId(), nil)
		if session == nil {
			conn.Deliver(action.newErrorMessage("无效的用户"))
			return gerror.New("无效的用户")
		}
		msg.AdminId = conn.GetUserId()
		msg.Source = consts.MessageSourceAdmin
		msg.SessionId = session.Id
		msg, err := service.ChatMessage().Insert(ctx, msg)
		if err != nil {
			return err
		}
		_ = service.ChatRelation().UpdateUser(ctx, msg.AdminId, msg.UserId)
		// 服务器回执d
		conn.Deliver(action.newReceipt(msg))
		return userM.DeliveryMessage(ctx, msg)
	}
	return nil
}

func (m *adminManager) registerHook(conn iWsConn) {
	ctx := gctx.New()
	err := m.broadcastOnlineAdmins(ctx, conn.GetCustomerId())
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
	}
	err = m.broadcastWaitingUser(ctx, conn.GetCustomerId())
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
	}
	err = m.noticeUserTransfer(ctx, conn.GetCustomerId(), conn.GetUserId())
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
	}
}

// conn断开连接后，更新admin的最后在线时间
func (m *adminManager) unregisterHook(conn iWsConn) {
	ctx := gctx.New()
	err := service.Admin().UpdateLastOnline(ctx, conn.GetUserId())
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
	}
	err = m.broadcastOnlineAdmins(ctx, conn.GetCustomerId())
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
	}
}

func (m *adminManager) broadcastWaitingUser(ctx context.Context, customerId uint) error {
	return m.broadcastLocalWaitingUser(ctx, customerId)
}

func (m *adminManager) broadcastLocalWaitingUser(ctx context.Context, customerId uint) (err error) {
	sessions, err := service.ChatSession().GetUnAccepts(ctx, customerId)
	if err != nil {
		return
	}
	sessionIds := slice.Map(sessions, func(index int, item *model.CustomerChatSession) uint {
		return item.Id
	})
	userMap := make(map[uint]*api.ChatWaitingUser)
	messages, err := service.ChatMessage().All(ctx, do.CustomerChatMessages{
		Source:    consts.MessageSourceUser,
		SessionId: sessionIds,
	}, nil, "id")
	if err != nil {
		return
	}
	for _, session := range sessions {
		userMap[session.UserId] = &api.ChatWaitingUser{
			Username:     session.User.Username,
			Avatar:       "",
			UserId:       session.User.Id,
			MessageCount: 0,
			Description:  "",
			Messages:     make([]api.ChatSimpleMessage, 0),
			LastTime:     session.QueriedAt,
			SessionId:    session.Id,
		}
	}
	for _, m := range messages {
		userMap[m.UserId].Messages = append(userMap[m.UserId].Messages, api.ChatSimpleMessage{
			Type:    m.Type,
			Time:    m.ReceivedAt,
			Content: m.Content,
		})
		userMap[m.UserId].MessageCount += 1
	}

	waitingUser := maputil.Values(userMap)
	sort.Slice(waitingUser, func(i, j int) bool {
		return waitingUser[i].LastTime.Unix() > waitingUser[j].LastTime.Unix()
	})
	adminConns := m.GetAllConn(customerId)
	act := action.newWaitingUsers(waitingUser)
	m.SendAction(act, adminConns...)
	return
}

func (m *adminManager) broadcastOnlineAdmins(ctx context.Context, gid uint) error {
	return m.broadcastLocalOnlineAdmins(ctx, gid)
}

func (m *adminManager) broadcastLocalOnlineAdmins(ctx context.Context, customerId uint) error {
	admins, err := service.Admin().All(ctx, do.CustomerAdmins{
		CustomerId: customerId,
	}, g.Slice{
		model.CustomerAdmin{}.Setting,
	}, nil)
	if err != nil {
		return err
	}
	data := make([]api.ChatCustomerAdmin, 0, len(admins))
	for _, c := range admins {
		conn, online := m.GetConn(customerId, c.Id)
		platform := ""
		if online {
			platform = conn.GetPlatform()
		}
		avatar, _ := service.Admin().GetAvatar(ctx, c)
		data = append(data, api.ChatCustomerAdmin{
			Avatar:        avatar,
			Username:      c.Username,
			Online:        online,
			Id:            c.Id,
			AcceptedCount: service.ChatRelation().GetActiveCount(ctx, c.Id),
			Platform:      platform,
		})
	}
	conns := m.GetAllConn(customerId)
	m.SendAction(action.newAdmins(data), conns...)
	return nil
}

func (m *adminManager) noticeRate(message *model.CustomerChatMessage) {
	action := action.newRateAction(message)
	conn, exist := m.GetConn(message.CustomerId, message.AdminId)
	if exist {
		conn.Deliver(action)
	}
}

func (m *adminManager) noticeUserOffline(user IChatUser) {
	m.noticeLocalUserOffline(user.GetPrimaryKey())
}

func (m *adminManager) noticeLocalUserOffline(uid uint) {
	adminId := service.ChatRelation().GetUserValidAdmin(gctx.New(), uid)
	admin, _ := service.Admin().First(gctx.New(), do.CustomerAdmins{Id: adminId})
	if admin != nil {
		conn, exist := m.GetConn(admin.CustomerId, admin.Id)
		if exist {
			m.SendAction(action.newUserOffline(uid), conn)
		}
	}
}

func (m *adminManager) noticeUserOnline(ctx context.Context, conn iWsConn) {
	m.noticeLocalUserOnline(ctx, conn.GetUserId(), conn.GetPlatform())
}

func (m *adminManager) noticeLocalUserOnline(ctx context.Context, uid uint, platform string) {
	adminId := service.ChatRelation().GetUserValidAdmin(ctx, uid)
	admin, _ := service.Admin().First(ctx, do.CustomerAdmins{Id: adminId})
	if admin != nil {
		conn, exist := m.GetConn(admin.CustomerId, admin.Id)
		if exist {
			m.SendAction(action.newUserOnline(uid, platform), conn)
		}
	}
}

//func (m *adminManager) noticeRepeatConnect(admin IChatUser) {
//	m.noticeLocalRepeatConnect(admin)
//}
//
//
//func (m *adminManager) noticeLocalRepeatConnect(admin IChatUser) {
//	conn, exist := m.GetConn(admin.GetCustomerId(), admin.GetPrimaryKey())
//	if exist && conn.GetUuid() != m.GetUserUuid(admin) {
//		m.SendAction(NewOtherLogin(), conn)
//	}
//}

func (m *adminManager) noticeUserTransfer(ctx context.Context, customerId, adminId uint) error {
	return m.noticeLocalUserTransfer(ctx, customerId, adminId)
}

func (m *adminManager) noticeLocalUserTransfer(ctx context.Context, customerId, adminId uint) error {
	client, exist := m.GetConn(customerId, adminId)
	if exist {
		transfers, err := service.ChatTransfer().All(ctx, g.Map{
			"to_admin_id":         adminId,
			"accepted_at is null": nil,
			"canceled_at is null": nil,
		}, g.Slice{
			model.CustomerChatTransfer{}.ToAdmin,
			model.CustomerChatTransfer{}.FormAdmin,
			model.CustomerChatTransfer{}.ToSession,
			model.CustomerChatTransfer{}.User,
		}, nil)
		if err != nil {
			return err
		}
		data := slice.Map(transfers, func(index int, item *model.CustomerChatTransfer) api.ChatTransfer {
			return service.ChatTransfer().ToApi(item)
		})
		client.Deliver(action.newUserTransfer(data))
	}
	return nil
}

// NoticeUpdateSetting admin修改设置后通知conn 更新admin的设置信息
func (m *adminManager) noticeUpdateSetting(customerId uint, setting *api.CurrentAdminSetting) {
	//m.updateSetting(customerId, setting)
}

// UpdateSetting 更新设置
func (m *adminManager) updateSetting(a *model.CustomerAdmin) {
	conn, exist := m.GetConn(a.CustomerId, a.Id)
	if exist {
		u, ok := conn.GetUser().(*admin)
		if ok {
			setting, err := service.Admin().FindSetting(gctx.New(), a.Id, true)
			if err == nil {
				u.Entity.Setting = setting
			}
		}
	}
}
