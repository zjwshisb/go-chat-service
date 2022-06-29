package websocket

import (
	"fmt"
	"github.com/spf13/viper"
	"sort"
	"time"
	"ws/app/chat"
	"ws/app/contract"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/resource"
	rpcClient "ws/app/rpc/client"
	"ws/app/util"
)

var AdminManager *adminManager

const TypeAdmin = "admin"

type adminManager struct {
	manager
}

func SetupAdmin() {
	AdminManager = &adminManager{
		manager: manager{
			shardCount:   10,
			ipAddr:       util.GetIPs()[0] + ":" + viper.GetString("Rpc.Port"),
			ConnMessages: make(chan *ConnMessage, 100),
			types:        TypeAdmin,
		},
	}
	AdminManager.onRegister = AdminManager.registerHook
	AdminManager.onUnRegister = AdminManager.unregisterHook
	AdminManager.Run()
}

func (m *adminManager) Run() {
	m.manager.Run()
	go m.handleReceiveMessage()
}

// DeliveryMessage
// 投递消息
// 查询admin是否在本机上，是则直接投递
// 查询admin当前channel，如果存在则投递到该channel上
// 最后则说明admin不在线，处理离线逻辑
func (m *adminManager) DeliveryMessage(msg *models.Message, isRemote bool) {
	adminConn, exist := m.GetConn(msg.GetAdmin())
	if exist { // admin在线且在当前服务上
		UserManager.triggerMessageEvent(models.SceneAdminOnline, msg)
		adminConn.Deliver(NewReceiveAction(msg))
		return
	} else if !isRemote && m.isCluster() {
		server := m.getUserServer(msg.AdminId) // 获取用户所在channel
		if server != "" {
			rpcClient.SendMessage(msg.Id, server)
			return
		}
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

// 处理离线消息
func (m *adminManager) handleOffline(msg *models.Message) {
	UserManager.triggerMessageEvent(models.SceneAdminOffline, msg)
	admin := repositories.AdminRepo.FirstById(msg.AdminId)
	setting := admin.GetSetting()
	if setting != nil {
		// 发送离线消息
		if setting.OfflineContent != "" {
			offlineMsg := setting.GetOfflineMsg(msg.UserId, msg.SessionId, msg.GroupId)
			offlineMsg.Admin = admin
			repositories.MessageRepo.Save(offlineMsg)
			UserManager.DeliveryMessage(offlineMsg, false)
		}
		// 判断是否自动断开
		lastOnline := setting.LastOnline
		duration := chat.SettingService.GetOfflineDuration(msg.GroupId)
		if (lastOnline.Unix() + duration) < time.Now().Unix() {
			chat.SessionService.Close(msg.SessionId, false, true)
			noticeMessage := admin.GetBreakMessage(msg.UserId, msg.SessionId) // 断开提醒
			noticeMessage.Save()
			UserManager.DeliveryMessage(noticeMessage, false)
		}
	}
}

// 处理消息
func (m *adminManager) handleMessage(payload *ConnMessage) {
	act := payload.Action
	conn := payload.Conn
	switch act.Action {
	// 客服发送消息给用户
	case SendMessageAction:
		msg, err := act.GetMessage()
		if err == nil {
			if msg.UserId > 0 && len(msg.Content) != 0 {
				if !chat.AdminService.IsUserValid(conn.GetUserId(), msg.UserId) {
					conn.Deliver(NewErrorMessage("该用户已失效，无法发送消息"))
					return
				}
				session := repositories.ChatSessionRepo.FirstActiveByUser(msg.UserId, conn.GetUserId())
				if session == nil {
					conn.Deliver(NewErrorMessage("无效的用户"))
					return
				}
				msg.GroupId = conn.GetGroupId()
				msg.AdminId = conn.GetUserId()
				msg.Source = models.SourceAdmin
				msg.ReceivedAT = time.Now().Unix()
				msg.Admin = conn.User.(*models.Admin)
				msg.SessionId = session.Id
				repositories.MessageRepo.Save(msg)
				_ = chat.AdminService.UpdateUser(msg.AdminId, msg.UserId)
				// 服务器回执d
				conn.Deliver(NewReceiptAction(msg))
				UserManager.DeliveryMessage(msg, false)
			}
		}
	}
}

func (m *adminManager) registerHook(conn Conn) {
	m.NoticeUserTransfer(conn.GetUser())
	m.BroadcastOnlineAdmins(conn.GetGroupId())
	m.BroadcastWaitingUser(conn.GetGroupId())
}

// conn断开连接后，更新admin的最后在线时间
func (m *adminManager) unregisterHook(conn Conn) {
	u := conn.GetUser()
	admin, ok := u.(*models.Admin)
	if ok {
		setting := admin.GetSetting()
		repositories.AdminRepo.UpdateSetting(setting, "last_online", time.Now())
	}
	m.BroadcastOnlineAdmins(conn.GetGroupId())
}

func (m *adminManager) BroadcastWaitingUser(groupId int64) {
	m.Do(func() {
		rpcClient.BroadcastWaitingUser(groupId)
	}, func() {
		m.BroadcastLocalWaitingUser(groupId)
	})
}

func (m *adminManager) BroadcastLocalWaitingUser(groupId int64) {
	sessions := repositories.ChatSessionRepo.GetWaitHandles()
	userMap := make(map[int64]*models.User)
	waitingUser := make([]*resource.WaitingChatSession, 0, len(sessions))
	for _, session := range sessions {
		userMap[session.UserId] = session.User
		msgs := make([]*resource.SimpleMessage, 0, len(session.Messages))
		for _, m := range session.Messages {
			msgs = append(msgs, &resource.SimpleMessage{
				Type:    m.Type,
				Time:    m.ReceivedAT,
				Content: m.Content,
			})
		}
		waitingUser = append(waitingUser, &resource.WaitingChatSession{
			Username:     session.User.GetUsername(),
			Avatar:       session.User.GetAvatarUrl(),
			UserId:       session.User.GetPrimaryKey(),
			MessageCount: len(session.Messages),
			Description:  "",
			Messages:     msgs,
			LastTime:     session.QueriedAt,
			SessionId:    session.Id,
		})
	}
	sort.Slice(waitingUser, func(i, j int) bool {
		return waitingUser[i].LastTime > waitingUser[j].LastTime
	})
	adminConns := m.GetAllConn(groupId)
	for _, conn := range adminConns {
		adminUserSlice := make([]*resource.WaitingChatSession, 0)
		for _, userJson := range waitingUser {
			u := userMap[userJson.UserId]
			admin := conn.GetUser().(*models.Admin)
			if admin.AccessTo(u) {
				adminUserSlice = append(adminUserSlice, userJson)
			}
		}
		conn.Deliver(NewWaitingUsers(adminUserSlice))
	}
}

func (m *adminManager) BroadcastOnlineAdmins(gid int64) {
	m.Do(func() {
		rpcClient.BroadcastOnlineAdmin(gid)
	}, func() {
		m.BroadcastLocalOnlineAdmins(gid)
	})
}

func (m *adminManager) BroadcastLocalOnlineAdmins(gid int64) {
	ids := m.GetOnlineUserIds(gid)
	admins := repositories.AdminRepo.Get([]*repositories.Where{{
		Filed: "id in ?",
		Value: ids,
	}}, -1, []string{}, []string{})
	data := make([]resource.Admin, 0, len(admins))
	for _, admin := range admins {
		data = append(data, resource.Admin{
			Avatar:        admin.GetAvatarUrl(),
			Username:      admin.Username,
			Online:        true,
			Id:            admin.GetPrimaryKey(),
			AcceptedCount: chat.AdminService.GetActiveCount(admin.GetPrimaryKey()),
		})
	}
	m.SendAction(NewAdminsAction(data), m.GetAllConn(gid)...)
}

func (m *adminManager) NoticeUserOffline(user contract.User) {
	m.Do(func() {
		adminId := chat.UserService.GetValidAdmin(user.GetPrimaryKey())
		server := m.getUserServer(adminId)
		fmt.Println(server)
		if server != "" {
			rpcClient.NoticeUserOffLine(user.GetPrimaryKey(), server)
		}
	}, func() {
		m.NoticeLocalUserOffline(user.GetPrimaryKey())
	})
}

func (m *adminManager) NoticeLocalUserOffline(uid int64) {
	adminId := chat.UserService.GetValidAdmin(uid)
	admin := repositories.AdminRepo.FirstById(adminId)
	if admin != nil {
		conn, exist := m.GetConn(admin)
		if exist {
			m.SendAction(NewUserOffline(uid), conn)
		}
	}
}

func (m *adminManager) NoticeUserOnline(user contract.User) {
	m.Do(func() {
		adminId := chat.UserService.GetValidAdmin(user.GetPrimaryKey())
		server := m.getUserServer(adminId)
		if server != "" {
			rpcClient.NoticeUserOnline(user.GetPrimaryKey(), server)
		}
	}, func() {
		m.NoticeLocalUserOnline(user.GetPrimaryKey())
	})
}

func (m *adminManager) NoticeLocalUserOnline(uid int64) {
	adminId := chat.UserService.GetValidAdmin(uid)
	admin := repositories.AdminRepo.FirstById(adminId)
	if admin != nil {
		conn, exist := m.GetConn(admin)
		if exist {
			m.SendAction(NewUserOnline(uid), conn)
		}
	}
}

//func (m *adminManager) NoticeRepeatConnect(admin contract.User) {
//	m.Do(func() {
//		//rpcClient.C
//	}, func() {
//		m.NoticeLocalRepeatConnect(admin)
//	})
//}
//
//func (m *adminManager) NoticeLocalRepeatConnect(admin contract.User) {
//	conn, exist := m.GetConn(admin)
//	if exist && conn.GetUuid() != m.GetUserUuid(admin) {
//		m.SendAction(NewOtherLogin(), conn)
//	}
//}

func (m *adminManager) NoticeUserTransfer(admin contract.User) {
	m.Do(func() {
		server := m.getUserServer(admin.GetPrimaryKey())
		if server != "" {
			rpcClient.NoticeUserTransfer(admin.GetPrimaryKey(), server)
		}
	}, func() {
		m.NoticeLocalUserTransfer(admin)
	})
}

func (m *adminManager) NoticeLocalUserTransfer(admin contract.User) {
	client, exist := m.GetConn(admin)
	if exist {
		transfers := repositories.TransferRepo.Get([]*repositories.Where{
			{
				Filed: "to_admin_id = ?",
				Value: admin.GetPrimaryKey(),
			},
			{
				Filed: "is_accepted = ?",
				Value: 0,
			},
			{
				Filed: "is_canceled",
				Value: 0,
			},
		}, -1, []string{"FromAdmin", "User"}, []string{"id desc"})
		data := make([]*resource.ChatTransfer, 0, len(transfers))
		for _, transfer := range transfers {
			data = append(data, transfer.ToJson())
		}
		client.Deliver(NewUserTransfer(data))
	}
}

// PublishUpdateSetting admin修改设置后通知conn 更新admin的设置信息
func (m *adminManager) PublishUpdateSetting(admin contract.User) {
	m.Do(func() {
		//channel := m.getUserChannel(admin.GetPrimaryKey())
		//if channel != "" {
		//	_ = m.publish(channel, &mq.Payload{
		//		Types: mq.TypeUpdateSetting,
		//		Data:  admin.GetPrimaryKey(),
		//	})
		//}
	}, func() {
		m.updateSetting(admin)
	})
}

// 更新设置
func (m *adminManager) updateSetting(admin contract.User) {
	conn, exist := m.GetConn(admin)
	if exist {
		u, ok := conn.GetUser().(*models.Admin)
		if ok {
			u.RefreshSetting()
		}
	}
}
