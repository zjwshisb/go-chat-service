package chat

import (
	"gf-chat/internal/consts"
	"gf-chat/internal/contract"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/official"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"sort"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
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
func (m *adminManager) deliveryMessage(msg *relation.CustomerChatMessages) {
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

func (m *adminManager) sendWaiting(admin *relation.CustomerAdmins, user contract.IChatUser) {
	offline := official.Offline{
		First:    "有新的用户咨询待接入",
		Remark:   "点击卡片查看详情",
		Keyword1: user.GetUsername(),
		Keyword3: "网页",
		Keyword4: time.Now().Format("2006-01-02 15:04:05"),
	}
	service.OfficialMsg().Chat(admin, offline)
}

// 发送公众号提醒
func (m *adminManager) sendOffline(admin *relation.CustomerAdmins, msg *relation.CustomerChatMessages) {
	tags := ""
	offline := official.Offline{
		First:    "发来了一条信息",
		Remark:   "点击卡片查看详情",
		Keyword1: msg.User.Username,
		Keyword2: tags,
		Keyword3: "网页",
		Keyword4: time.Now().Format("2006-01-02 15:04:05"),
	}
	service.OfficialMsg().Chat(admin, offline)
}

// 处理离线消息
func (m *adminManager) handleOffline(msg *relation.CustomerChatMessages) {
	userM.triggerMessageEvent(consts.AutoRuleSceneAdminOffline, msg, &user{Entity: msg.User})
	admin := service.Admin().FirstRelation(msg.AdminId)
	if admin != nil {
		// 离线消息
		message := service.ChatMessage().NewOffline(admin)
		if message != nil {
			message.UserId = msg.UserId
			message.SessionId = msg.SessionId
			service.ChatMessage().SaveRelationOne(message)
			userM.DeliveryMessage(message)
		}
		m.sendOffline(admin, msg)
	}
}

// 处理消息
func (m *adminManager) handleMessage(payload *chatConnMessage) {
	msg := payload.Msg
	conn := payload.Conn
	if msg.UserId > 0 {
		if !service.ChatRelation().IsUserValid(conn.GetUserId(), msg.UserId) {
			conn.Deliver(service.Action().NewErrorMessage("该用户已失效，无法发送消息"))
			return
		}
		session := service.ChatSession().ActiveOne(msg.UserId, conn.GetUserId(), nil)
		if session == nil {
			conn.Deliver(service.Action().NewErrorMessage("无效的用户"))
			return
		}
		msg.AdminId = conn.GetUserId()
		msg.Source = consts.MessageSourceAdmin
		msg.ReceivedAt = gtime.New()
		msg.SessionId = session.Id
		service.ChatMessage().SaveRelationOne(msg)
		_ = service.ChatRelation().UpdateUser(msg.AdminId, msg.UserId)
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
	sessions := service.ChatSession().GetUnAcceptModel(customerId)
	sessionIds := slice.Map(sessions, func(index int, item *relation.CustomerChatSessions) uint {
		return item.Id
	})
	userMap := make(map[uint]*chat.WaitingUser)
	messages := make([]entity.CustomerChatMessages, 0)
	_ = dao.CustomerChatMessages.Ctx(gctx.New()).Where("session_id in (?)", sessionIds).
		Where("source", consts.MessageSourceUser).
		Order("id").
		Scan(&messages)
	for _, session := range sessions {
		userMap[session.UserId] = &chat.WaitingUser{
			Username:     session.User.Username,
			Avatar:       "",
			UserId:       session.User.Id,
			MessageCount: 0,
			Description:  "",
			Messages:     make([]chat.SimpleMessage, 0),
			LastTime:     session.QueriedAt,
			SessionId:    session.Id,
		}
	}
	for _, m := range messages {
		userMap[m.UserId].Messages = append(userMap[m.UserId].Messages, chat.SimpleMessage{
			Type:    m.Type,
			Time:    m.ReceivedAt,
			Content: m.Content,
		})
		userMap[m.UserId].MessageCount += 1
	}

	waitingUser := maputil.Values(userMap)
	sort.Slice(waitingUser, func(i, j int) bool {
		return waitingUser[i].LastTime.Unix > waitingUser[j].LastTime.Unix
	})
	adminConns := m.GetAllConn(customerId)
	action := service.Action().NewWaitingUsers(waitingUser)
	m.SendAction(action, adminConns...)
}

func (m *adminManager) broadcastOnlineAdmins(gid uint) {
	m.broadcastLocalOnlineAdmins(gid)
}

func (m *adminManager) broadcastLocalOnlineAdmins(customerId uint) {
	admins := service.Admin().GetChatAll(customerId)
	data := make([]chat.CustomerAdmin, 0, len(admins))
	for _, c := range admins {
		conn, online := m.GetConn(customerId, c.Id)
		platform := ""
		if online {
			platform = conn.GetPlatform()
		}
		data = append(data, chat.CustomerAdmin{
			Avatar:        service.Admin().GetAvatar(c),
			Username:      c.Username,
			Online:        online,
			Id:            c.Id,
			AcceptedCount: service.ChatRelation().GetActiveCount(c.Id),
			Platform:      platform,
		})
	}
	conns := m.GetAllConn(customerId)
	m.SendAction(service.Action().NewAdminsAction(data), conns...)
}

func (m *adminManager) noticeRate(message *entity.CustomerChatMessages) {
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
	adminId := service.ChatRelation().GetUserValidAdmin(uid)
	admin := service.Admin().First(adminId)
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
	adminId := service.ChatRelation().GetUserValidAdmin(uid)
	admin := service.Admin().First(adminId)
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
		transfers := service.ChatTransfer().GetRelations(do.CustomerChatTransfers{
			ToAdminId:  adminId,
			AcceptedAt: nil,
			CanceledAt: nil,
		})
		data := slice.Map(transfers, func(index int, item *relation.CustomerChatTransfer) chat.Transfer {
			return service.ChatTransfer().RelationToChat(item)
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
