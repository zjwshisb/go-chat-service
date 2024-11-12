package chat

import (
	"gf-chat/internal/consts"
	"gf-chat/internal/contract"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"sort"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

const TypeAdmin = "admin"

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
func (m *adminManager) deliveryMessage(msg *model.CustomerChatMessage) {
	adminConn, exist := m.GetConn(msg.CustomerId, msg.AdminId)
	if exist { // admin在线
		adminConn.Deliver(service.Action().NewReceiveAction(msg))
		return
	}
	m.handleOffline(msg)
}

// 从管道接受消息并处理
func (m *adminManager) handleReceiveMessage() {
	for {
		payload := <-m.ConnMessages
		go m.handleMessage(payload)
	}
}

func (m *adminManager) sendWaiting(admin *model.CustomerAdmin, user contract.IChatUser) {

}

func (m *adminManager) sendOffline(admin *model.CustomerAdmin, msg *model.CustomerChatMessage) {

}

// 处理离线消息
func (m *adminManager) handleOffline(msg *model.CustomerChatMessage) {
	userM.triggerMessageEvent(consts.AutoRuleSceneAdminOffline, msg, &user{Entity: msg.User})
	ctx := gctx.New()
	admin, err := service.Admin().First(ctx, do.CustomerAdmins{Id: msg.AdminId})
	if err != nil {
		return
	}
	message := service.ChatMessage().NewOffline(admin)
	if message != nil {
		message.UserId = msg.UserId
		message.SessionId = msg.SessionId
		_ = service.ChatMessage().SaveWithUpdate(ctx, message)
		userM.DeliveryMessage(message)
	}
	m.sendOffline(admin, msg)
}

// 处理消息
func (m *adminManager) handleMessage(payload *chatConnMessage) {
	msg := payload.Msg
	conn := payload.Conn
	ctx := gctx.New()
	if msg.UserId > 0 {
		if !service.ChatRelation().IsUserValid(gctx.New(), conn.GetUserId(), msg.UserId) {
			conn.Deliver(service.Action().NewErrorMessage("该用户已失效，无法发送消息"))
			return
		}
		session, _ := service.ChatSession().ActiveOne(ctx, msg.UserId, conn.GetUserId(), nil)
		if session == nil {
			conn.Deliver(service.Action().NewErrorMessage("无效的用户"))
			return
		}
		msg.AdminId = conn.GetUserId()
		msg.Source = consts.MessageSourceAdmin
		msg.ReceivedAt = gtime.New()
		msg.SessionId = session.Id
		_ = service.ChatMessage().SaveWithUpdate(ctx, msg)
		_ = service.ChatRelation().UpdateUser(gctx.New(), msg.AdminId, msg.UserId)
		// 服务器回执d
		conn.Deliver(service.Action().NewReceiptAction(msg))
		userM.DeliveryMessage(msg)
	}

}

func (m *adminManager) registerHook(conn iWsConn) {
	m.broadcastOnlineAdmins(conn.GetCustomerId())
	m.broadcastWaitingUser(conn.GetCustomerId())
	m.noticeUserTransfer(conn.GetCustomerId(), conn.GetUserId())
}

// conn断开连接后，更新admin的最后在线时间
func (m *adminManager) unregisterHook(conn iWsConn) {
	u := conn.GetUser()
	a, ok := u.(*admin)
	if ok {
		e := a.Entity
		e.Setting.LastOnline = gtime.New()
		dao.CustomerAdminChatSettings.Ctx(gctx.New()).Save(e.Setting)
	}
	m.broadcastOnlineAdmins(conn.GetCustomerId())
}

func (m *adminManager) broadcastWaitingUser(customerId uint) {
	m.broadcastLocalWaitingUser(customerId)
}

func (m *adminManager) broadcastLocalWaitingUser(customerId uint) {
	ctx := gctx.New()
	sessions, _ := service.ChatSession().GetUnAcceptModel(ctx, customerId)
	sessionIds := slice.Map(sessions, func(index int, item *model.CustomerChatSession) uint {
		return item.Id
	})
	userMap := make(map[uint]*model.ChatWaitingUser)
	messages := make([]entity.CustomerChatMessages, 0)
	_ = dao.CustomerChatMessages.Ctx(gctx.New()).Where("session_id in (?)", sessionIds).
		Where("source", consts.MessageSourceUser).
		Order("id").
		Scan(&messages)
	for _, session := range sessions {
		userMap[session.UserId] = &model.ChatWaitingUser{
			Username:     session.User.Username,
			Avatar:       "",
			UserId:       session.User.Id,
			MessageCount: 0,
			Description:  "",
			Messages:     make([]model.ChatSimpleMessage, 0),
			LastTime:     session.QueriedAt,
			SessionId:    session.Id,
		}
	}
	for _, m := range messages {
		userMap[m.UserId].Messages = append(userMap[m.UserId].Messages, model.ChatSimpleMessage{
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
	action := service.Action().NewWaitingUsers(waitingUser)
	m.SendAction(action, adminConns...)
}

func (m *adminManager) broadcastOnlineAdmins(gid uint) {
	m.broadcastLocalOnlineAdmins(gid)
}

func (m *adminManager) broadcastLocalOnlineAdmins(customerId uint) {
	admins, _ := service.Admin().All(gctx.New(), do.CustomerAdmins{
		CustomerId: customerId,
	}, g.Slice{
		model.CustomerAdmin{}.Setting,
	}, nil)
	data := make([]model.ChatCustomerAdmin, 0, len(admins))
	for _, c := range admins {
		conn, online := m.GetConn(customerId, c.Id)
		platform := ""
		if online {
			platform = conn.GetPlatform()
		}
		avatar, _ := service.Admin().GetAvatar(gctx.New(), c)
		data = append(data, model.ChatCustomerAdmin{
			Avatar:        avatar,
			Username:      c.Username,
			Online:        online,
			Id:            c.Id,
			AcceptedCount: service.ChatRelation().GetActiveCount(gctx.New(), c.Id),
			Platform:      platform,
		})
	}
	conns := m.GetAllConn(customerId)
	m.SendAction(service.Action().NewAdminsAction(data), conns...)
}

func (m *adminManager) noticeRate(message *model.CustomerChatMessage) {
	action := service.Action().NewRateAction(message)
	conn, exist := m.GetConn(message.CustomerId, message.AdminId)
	if exist {
		conn.Deliver(action)
	}
}

func (m *adminManager) noticeUserOffline(user contract.IChatUser) {
	m.noticeLocalUserOffline(user.GetPrimaryKey())
}

func (m *adminManager) noticeLocalUserOffline(uid uint) {
	adminId := service.ChatRelation().GetUserValidAdmin(gctx.New(), uid)
	admin, _ := service.Admin().First(gctx.New(), do.CustomerAdmins{Id: adminId})
	if admin != nil {
		conn, exist := m.GetConn(admin.CustomerId, admin.Id)
		if exist {
			m.SendAction(service.Action().NewUserOffline(uid), conn)
		}
	}
}

func (m *adminManager) noticeUserOnline(conn iWsConn) {
	m.noticeLocalUserOnline(conn.GetUserId(), conn.GetPlatform())
}

func (m *adminManager) noticeLocalUserOnline(uid uint, platform string) {
	adminId := service.ChatRelation().GetUserValidAdmin(gctx.New(), uid)
	admin, _ := service.Admin().First(gctx.New(), do.CustomerAdmins{Id: adminId})
	if admin != nil {
		conn, exist := m.GetConn(admin.CustomerId, admin.Id)
		if exist {
			m.SendAction(service.Action().NewUserOnline(uid, platform), conn)
		}
	}
}

//func (m *adminManager) noticeRepeatConnect(admin contract.IChatUser) {
//	m.noticeLocalRepeatConnect(admin)
//}
//
//
//func (m *adminManager) noticeLocalRepeatConnect(admin contract.IChatUser) {
//	conn, exist := m.GetConn(admin.GetCustomerId(), admin.GetPrimaryKey())
//	if exist && conn.GetUuid() != m.GetUserUuid(admin) {
//		m.SendAction(NewOtherLogin(), conn)
//	}
//}

func (m *adminManager) noticeUserTransfer(customerId, adminId uint) {
	m.noticeLocalUserTransfer(customerId, adminId)
}

func (m *adminManager) noticeLocalUserTransfer(customerId, adminId uint) {
	client, exist := m.GetConn(customerId, adminId)
	if exist {
		transfers, _ := service.ChatTransfer().All(gctx.New(), do.CustomerChatTransfers{
			ToAdminId:  adminId,
			AcceptedAt: nil,
			CanceledAt: nil,
		}, g.Slice{model.CustomerChatTransfer{}.FormAdmin,
			model.CustomerChatTransfer{}.ToAdmin,
			model.CustomerChatTransfer{}.FormAdmin,
			model.CustomerChatTransfer{}.ToSession,
			model.CustomerChatTransfer{}.User,
		}, nil)
		data := slice.Map(transfers, func(index int, item *model.CustomerChatTransfer) model.ChatTransfer {
			return service.ChatTransfer().ToChatTransfer(item)
		})
		client.Deliver(service.Action().NewUserTransfer(data))
	}
}

// NoticeUpdateSetting admin修改设置后通知conn 更新admin的设置信息
func (m *adminManager) noticeUpdateSetting(customerId uint, setting *entity.CustomerAdminChatSettings) {
	m.updateSetting(customerId, setting)
}

// UpdateSetting 更新设置
func (m *adminManager) updateSetting(customerId uint, setting *entity.CustomerAdminChatSettings) {
	conn, exist := m.GetConn(customerId, setting.AdminId)
	if exist {
		u, ok := conn.GetUser().(*admin)
		if ok {
			u.Entity.Setting = setting
		}
	}
}
