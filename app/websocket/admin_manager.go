package websocket

import (
	"fmt"
	"sort"
	"time"
	"ws/app/auth"
	"ws/app/chat"
	"ws/app/models"
	"ws/app/mq"
	"ws/app/resource"
	"ws/configs"
)

var AdminManager *adminManager

type adminManager struct {
	manager
}

func init()  {
	AdminManager = &adminManager{
		manager: manager{
			groupCount: 10,
			Channel: configs.App.Name + "-admin",
			ConnMessages: make(chan *ConnMessage, 100),
			userChannelCacheKey: "admin:%d:channel",
			groupCacheKey: "admin:channel:group",
		},
	}
	AdminManager.onRegister = AdminManager.registerHook
	AdminManager.onUnRegister = AdminManager.unregisterHook
	AdminManager.Run()
}

func (m *adminManager) Run()  {
	m.manager.Run()
	go m.handleReceiveMessage()
	if m.isCluster() {
		go m.handleRemoteMessage()
	}
}
// 投递消息
// 查询admin是否在本机上，是则直接投递
// 查询admin当前channel，如果存在则投递到该channel上
// 最后则说明admin不在线，处理离线逻辑
func (m *adminManager) DeliveryMessage(msg *models.Message)  {
	adminConn, exist := m.GetConn(msg.GetAdmin())
	if exist { // admin在线且在当前服务上
		UserManager.triggerMessageEvent(models.SceneAdminOnline, msg)
		adminConn.Deliver(NewReceiveAction(msg))
	} else if m.isCluster() {
		adminChannel := m.getUserChannel(msg.AdminId) // 获取用户所在channel
		if adminChannel != "" {
			_ = m.publish(adminChannel, mq.NewMessagePayload(msg.Id))
		}
	} else {
		m.handleOffline(msg)
	}
}
// 从管道接受消息并处理
func (m *adminManager) handleReceiveMessage() {
	for {
		payload := <- m.ConnMessages
		go m.handleMessage(payload)
	}
}
// 处理离线消息
func (m *adminManager) handleOffline(msg *models.Message)  {
	UserManager.triggerMessageEvent(models.SceneAdminOffline, msg)
	admin := adminRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: msg.AdminId,
		},
	})
	setting := admin.GetSetting()
	if setting != nil {
		// 发送离线消息
		if setting.OfflineContent != "" {
			offlineMsg := setting.GetOfflineMsg(msg.UserId, msg.SessionId)
			offlineMsg.Admin = admin
			messageRepo.Save(offlineMsg)
			UserManager.DeliveryMessage(offlineMsg)
		}
		// 判断是否自动断开
		lastOnline := setting.LastOnline
		duration := chat.SettingService.GetOfflineDuration()
		if (lastOnline.Unix() + duration) < time.Now().Unix() {
			chat.SessionService.Close(msg.SessionId, false, true)
			noticeMessage := admin.GetBreakMessage(msg.UserId, msg.SessionId) // 断开提醒
			noticeMessage.Save()
			UserManager.DeliveryMessage(noticeMessage)
		}
	}
}


// 订阅本服务的channel， 处理消息
func (m *adminManager) handleRemoteMessage()  {
	subscribe := mq.Mq().Subscribe(m.GetSubscribeChannel())
	defer subscribe.Close()
	for {
		message := subscribe.ReceiveMessage()
		go func() {
			switch message.Get("Types").String() {
			case mq.TypeWaitingUser:
				gid := message.Get("data").Int()
				if gid> 0 {
					m.BroadcastWaitingUser(gid)
				}
			case mq.TypeAdmin:
				m.BroadcastAdmins()
			case mq.TypeAdminLogin:
			case mq.TypeTransfer:
				adminIdI := message.Get("data").Int()
				admin := adminRepo.First([]Where{
					{
						Filed: "id = ?",
						Value: adminIdI,
					},
				})
				if admin != nil {
					m.BroadcastUserTransfer(admin)
				}
			case mq.TypeMessage:
				mid := message.Get("data").Int()
				msg := messageRepo.First(mid)
				if msg != nil {
					client, exist := m.GetConn(msg.GetAdmin())
					if exist {
						client.Deliver(NewReceiveAction(msg))
					} else {
						m.handleOffline(msg)
					}
				}
			}
		}()
	}
}


// 处理消息
func (m *adminManager) handleMessage(payload *ConnMessage)  {
	act := payload.Action
	conn := payload.Conn
	switch act.Action {
	// 客服发送消息给用户
	case SendMessageAction:
		msg, err := act.GetMessage()
		fmt.Println(msg)
		if err == nil {
			if msg.UserId > 0 && len(msg.Content) != 0 {
				if !chat.AdminService.IsUserValid(conn.GetUserId(), msg.UserId) {
					conn.Deliver(NewErrorMessage("该用户已失效，无法发送消息"))
					return
				}
				session := chat.SessionService.Get(msg.UserId, conn.GetUserId())
				if session == nil {
					conn.Deliver(NewErrorMessage("无效的用户"))
					return
				}
				sessionAddTime := chat.SettingService.GetUserSessionSecond()
				msg.AdminId = conn.GetUserId()
				msg.Source = models.SourceAdmin
				msg.ReceivedAT = time.Now().Unix()
				msg.Admin = conn.User.(*models.Admin)
				msg.SessionId = session.Id
				messageRepo.Save(msg)
				_ = chat.AdminService.UpdateUser(msg.AdminId, msg.UserId, sessionAddTime)
				// 服务器回执d
				conn.Deliver(NewReceiptAction(msg))
				UserManager.DeliveryMessage(msg)
			}
		}
	}
}



func (m *adminManager) registerHook(conn Conn)  {
	m.BroadcastUserTransfer(conn.GetUser())
	m.publishAdmins()
	m.BroadcastWaitingUser(conn.GetGroupId())
}

func (m *adminManager) unregisterHook(conn Conn)  {
	u := conn.GetUser()
	admin , ok := u.(*models.Admin)
	if ok {
		setting := admin.GetSetting()
		setting.LastOnline = time.Now()
		adminRepo.SaveSetting(setting)
	}
	m.BroadcastAdmins()
}
// 推送待接入用户
func (m *adminManager) PublishWaitingUser(groupId int64) {
	if m.isCluster() {
		channels := m.getAllChannel()
		for _, channel := range channels {
			_ = m.publish(channel, &mq.Payload{
				Types: mq.TypeWaitingUser,
				Data: groupId,
			})
		}
	} else {
		m.BroadcastWaitingUser(groupId)
	}
}
func (m *adminManager) PublishTransfer(admin auth.User)  {
	if m.isCluster() {
		channel := m.getUserChannel(admin.GetPrimaryKey())
		if channel != "" {
			_ = m.publish(channel, &mq.Payload{
				Types: mq.TypeTransfer,
				Data: admin.GetPrimaryKey(),
			})
		}
	} else {
		m.BroadcastUserTransfer(admin)
	}
}
// 推送在线admin
// todo
func (m *adminManager) publishAdmins()  {
	if m.isCluster() {
		channels := m.getAllChannel()
		for _, channel := range channels {
			_ = m.publish(channel, &mq.Payload{
				Types: mq.TypeAdmin,
			})
		}
	} else {
		m.BroadcastAdmins()
	}
}
func (m *adminManager) publishLogin() {

}
// 广播待接入用户
func (m *adminManager) BroadcastWaitingUser(groupId int64) {
	sessions :=  sessionRepo.GetWaitHandles()
	userMap := make(map[int64]*models.User)
	waitingUser :=  make([]*resource.WaitingChatSession, 0, len(sessions))
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
			UserId:           session.User.GetPrimaryKey(),
			MessageCount: len(session.Messages),
			Description:  "",
			Messages: msgs,
			LastTime: session.QueriedAt,
			SessionId: session.Id,
		})
	}
	sort.Slice(waitingUser, func(i, j int) bool {
		return waitingUser[i].LastTime < waitingUser[j].LastTime
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

// 向admin推送待转接入的用户
func (m *adminManager) BroadcastUserTransfer(admin auth.User)   {
	client, exist := m.GetConn(admin)
	if exist {
		transfers := transferRepo.Get([]Where{
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
		}, -1, []string{"FromAdmin", "User"}, "id desc")
		data := make([]*resource.ChatTransfer, 0, len(transfers))
		for _, transfer := range transfers {
			data = append(data, transfer.ToJson())
		}
		client.Deliver(NewUserTransfer(data))
	}
}
// 广播在线admin
func (m *adminManager) BroadcastAdmins() {
	//var serviceUsers []*models.Admin
	//admins := adminRepo.Get([]Where{}, -1, []string{})
	//conns := m.GetAllConn()
	//data := make([]resource.Admin, 0, len(serviceUsers))
	//for _, admin := range admins {
	//	_, online := m.GetConn(admin.GetPrimaryKey())
	//	if online {
	//		data = append(data, resource.Admin{
	//			Avatar:           admin.GetAvatarUrl(),
	//			Username:         admin.Username,
	//			Online:           online,
	//			Id:               admin.GetPrimaryKey(),
	//			AcceptedCount: chat.AdminService.GetActiveCount(admin.GetPrimaryKey()),
	//		})
	//	}
	//}
	//m.SendAction(NewAdminsAction(data), conns...)
}
