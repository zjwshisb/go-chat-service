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
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"sort"
	"time"
)

func newAdminManager() *adminManager {
	adminM = &adminManager{
		&manager{
			shardCount:   10,
			connMessages: make(chan *chatConnMessage, 100),
			pingDuration: time.Second * 10,
		},
	}
	adminM.on(eventRegister, adminM.onRegister)
	adminM.on(eventUnRegister, adminM.onUnRegister)
	adminM.on(eventMessage, adminM.handleMessage)
	return adminM
}

type adminManager struct {
	*manager
}

// deliveryMessage
// 投递消息
// 查询admin是否在线，是则直接投递
// 最后则说明admin不在线，处理离线逻辑
func (m *adminManager) deliveryMessage(ctx context.Context, msg *model.CustomerChatMessage, userConn iWsConn) error {
	adminConn, exist := m.getConn(msg.CustomerId, msg.AdminId)
	if exist { // admin在线
		err := userM.triggerMessageEvent(ctx, consts.AutoRuleSceneAdminOnline, msg, userConn)
		if err != nil {
			return nil
		}
		adminConn.deliver(action.newReceive(msg))
		return nil
	} else {
		return m.handleOffline(ctx, msg, userConn)
	}
}

func (m *adminManager) sendWaiting(admin *model.CustomerAdmin, user iChatUser) {

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
		err = userM.deliveryMessage(ctx, message)
		if err != nil {
			return err
		}
	}
	m.sendOffline(admin, msg)
	return nil
}

// 处理消息
func (m *adminManager) handleMessage(ctx context.Context, arg eventArg) error {
	msg := arg.msg
	conn := arg.conn
	if msg.UserId > 0 {
		if !service.ChatRelation().IsUserValid(ctx, conn.getUserId(), msg.UserId) {
			conn.deliver(action.newErrorMessage("该用户已失效，无法发送消息"))
			return gerror.NewCode(gcode.CodeValidationFailed, "该用户已失效，无法发送消息")
		}
		session, _ := service.ChatSession().FirstActive(ctx, msg.UserId, conn.getUserId(), nil)
		if session == nil {
			conn.deliver(action.newErrorMessage("无效的用户"))
			return gerror.NewCode(gcode.CodeValidationFailed, "无效的用户")
		}
		msg.AdminId = conn.getUserId()
		msg.Source = consts.MessageSourceAdmin
		msg.SessionId = session.Id
		msg, err := service.ChatMessage().Insert(ctx, msg)
		if err != nil {
			return err
		}
		_ = service.ChatRelation().UpdateUser(ctx, msg.AdminId, msg.UserId)
		// 服务器回执d
		conn.deliver(action.newReceipt(msg))
		return userM.deliveryMessage(ctx, msg)
	}
	return nil
}

func (m *adminManager) onRegister(ctx context.Context, arg eventArg) error {
	err := m.broadcastOnlineAdmins(ctx, arg.conn.getCustomerId())
	if err != nil {
		return err
	}
	err = m.broadcastWaitingUser(ctx, arg.conn.getCustomerId())
	if err != nil {
		return err
	}
	err = m.noticeUserTransfer(ctx, arg.conn.getCustomerId(), arg.conn.getUserId())
	if err != nil {
		return err
	}
	return nil
}

// conn断开连接后，更新admin的最后在线时间
func (m *adminManager) onUnRegister(ctx context.Context, arg eventArg) error {
	err := service.Admin().UpdateLastOnline(ctx, arg.conn.getUserId())
	if err != nil {
		return err
	}
	err = m.broadcastOnlineAdmins(ctx, arg.conn.getCustomerId())
	if err != nil {
		return err
	}
	return nil
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
	adminConns := m.getAllConn(customerId)
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
	}, nil, nil)
	if err != nil {
		return err
	}
	data := make([]api.ChatCustomerAdmin, 0, len(admins))
	for _, c := range admins {
		conn, online := m.getConn(customerId, c.Id)
		platform := ""
		if online {
			platform = conn.getPlatform()
		}
		data = append(data, api.ChatCustomerAdmin{
			Username:      c.Username,
			Online:        online,
			Id:            c.Id,
			AcceptedCount: service.ChatRelation().GetActiveCount(ctx, c.Id),
			Platform:      platform,
		})
	}
	conns := m.getAllConn(customerId)
	m.SendAction(action.newAdmins(data), conns...)
	return nil
}

func (m *adminManager) noticeRate(message *model.CustomerChatMessage) {
	action := action.newRateAction(message)
	conn, exist := m.getConn(message.CustomerId, message.AdminId)
	if exist {
		conn.deliver(action)
	}
}

func (m *adminManager) noticeUserOffline(user iChatUser) {
	m.noticeLocalUserOffline(user.getPrimaryKey())
}

func (m *adminManager) noticeLocalUserOffline(uid uint) {
	adminId := service.ChatRelation().GetUserValidAdmin(gctx.New(), uid)
	admin, _ := service.Admin().First(gctx.New(), do.CustomerAdmins{Id: adminId})
	if admin != nil {
		conn, exist := m.getConn(admin.CustomerId, admin.Id)
		if exist {
			m.SendAction(action.newUserOffline(uid), conn)
		}
	}
}

func (m *adminManager) noticeUserOnline(ctx context.Context, conn iWsConn) {
	m.noticeLocalUserOnline(ctx, conn.getUserId(), conn.getPlatform())
}

func (m *adminManager) noticeLocalUserOnline(ctx context.Context, uid uint, platform string) {
	adminId := service.ChatRelation().GetUserValidAdmin(ctx, uid)
	admin, _ := service.Admin().First(ctx, do.CustomerAdmins{Id: adminId})
	if admin != nil {
		conn, exist := m.getConn(admin.CustomerId, admin.Id)
		if exist {
			m.SendAction(action.newUserOnline(uid, platform), conn)
		}
	}
}

//func (m *adminManager) noticeRepeatConnect(admin iChatUser) {
//	m.noticeLocalRepeatConnect(admin)
//}
//
//
//func (m *adminManager) noticeLocalRepeatConnect(admin iChatUser) {
//	conn, exist := m.getConn(admin.getCustomerId(), admin.getPrimaryKey())
//	if exist && conn.getUuid() != m.GetUserUuid(admin) {
//		m.SendAction(NewOtherLogin(), conn)
//	}
//}

func (m *adminManager) noticeUserTransfer(ctx context.Context, customerId, adminId uint) error {
	return m.noticeLocalUserTransfer(ctx, customerId, adminId)
}

func (m *adminManager) noticeLocalUserTransfer(ctx context.Context, customerId, adminId uint) error {
	client, exist := m.getConn(customerId, adminId)
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
		client.deliver(action.newUserTransfer(data))
	}
	return nil
}

// NoticeUpdateSetting admin修改设置后通知conn 更新admin的设置信息
func (m *adminManager) noticeUpdateSetting(customerId uint, setting *api.CurrentAdminSetting) {
	//m.updateSetting(customerId, setting)
}

// UpdateSetting 更新设置
func (m *adminManager) updateSetting(ctx context.Context, a *model.CustomerAdmin) {
	conn, exist := m.getConn(a.CustomerId, a.Id)
	if exist {
		u, ok := conn.getUser().(*admin)
		if ok {
			setting, err := service.Admin().FindSetting(ctx, a.Id, true)
			if err == nil {
				u.Entity.Setting = setting
			}
		}
	}
}
