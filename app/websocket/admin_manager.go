package websocket

import (
	"context"
	"encoding/json"
	"sort"
	"time"
	"ws/app/chat"
	"ws/app/databases"
	"ws/app/models"
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
			Clients: make(map[int64]Conn),
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

// 从管道接受消息并处理
func (m *adminManager) handleReceiveMessage() {
	for {
		payload := <- m.ConnMessages
		go m.handleMessage(payload)
	}
}
// 处理离线下消息
func (m *adminManager) handleOffline(msg *models.Message)  {
	admin := adminRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: msg.AdminId,
		},
	})
	UserManager.triggerMessageEvent(models.SceneAdminOffline, msg)
	setting := admin.GetSetting()
	if setting != nil {
		// 发送离线消息
		if setting.OfflineContent != "" {
			offlineMsg := setting.GetOfflineMsg(msg.UserId, msg.SessionId)
			offlineMsg.Admin = admin
			messageRepo.Save(offlineMsg)
			UserManager.DeliveryMessage(offlineMsg, nil)
		}
		// 判断是否自动断开
		lastOnline := setting.LastOnline
		duration := chat.SettingService.GetOfflineDuration()
		if (lastOnline.Unix() + duration) < time.Now().Unix() {
			chat.SessionService.Close(msg.SessionId, false, true)
			noticeMessage := admin.GetBreakMessage(msg.UserId, msg.SessionId) // 断开提醒
			noticeMessage.Save()
			UserManager.DeliveryMessage(noticeMessage, nil)
		}
	}
}
// 投递消息
// 查询admin是否在本机上，是则直接投递
// 查询admin当前channel，如果存在则投递到该channel上
// 最后则说明admin不在线，处理离线逻辑
func (m *adminManager) DeliveryMessage(msg *models.Message, from Conn)  {
	adminConn, exist := m.GetConn(msg.AdminId)
	if exist { // admin在线且在当前服务上
		UserManager.triggerMessageEvent(models.SceneAdminOnline, msg)
		adminConn.Deliver(NewReceiveAction(msg))
	} else {
		adminChannel := m.getUserChannel(msg.AdminId) // 获取用户所在channel
		if adminChannel != "" {
			_ = m.publish(adminChannel, newMessagePayload(msg.Id))
		} else {
			m.handleOffline(msg)
		}
	}
}

// 订阅本服务的channel， 处理消息
func (m *adminManager) handleRemoteMessage()  {
	ctx := context.Background()
	sub := databases.Redis.Subscribe(ctx, m.GetSubscribeChannel())
	for {
		msg, err := sub.ReceiveMessage(ctx)
		go func() {
			payload := &payload{}
			err = json.Unmarshal([]byte(msg.Payload), payload)
			switch payload.Types {
			case TypeWaitingUser:
				m.BroadcastWaitingUser()
			case TypeAdmin:
				m.BroadcastAdmins()
			case TypeMessage:
				mid := payload.Data
				msg := messageRepo.First(mid)
				if msg != nil {
					client, exist := m.GetConn(msg.AdminId)
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
				UserManager.DeliveryMessage(msg, conn)
			}
		}
	}
}

func (m *adminManager) Run()  {
	m.manager.Run()
	go m.handleRemoteMessage()
	go m.handleReceiveMessage()
}

func (m *adminManager) registerHook(conn Conn)  {
	m.BroadcastUserTransfer(conn.GetUserId())
	m.publishAdmins()
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

func (m *adminManager) publishWaitingUser() {
	channels := m.getAllChannel()
	for _, channel := range channels {
		_ = m.publish(channel, &payload{
			Types: TypeWaitingUser,
		})
	}
}

func (m *adminManager) publishAdmins()  {
	channels := m.getAllChannel()
	for _, channel := range channels {
		_ = m.publish(channel, &payload{
			Types: TypeAdmin,
		})
	}
}

// 广播待接入用户
func (m *adminManager) BroadcastWaitingUser() {
	sessions := sessionRepo.Get([]Where{
		{
			Filed: "admin_id = ?",
			Value: 0,
		},
		{
			Filed: "canceled_at = ?",
			Value: 0,
		},
	}, -1, []string{"User", "Messages"})
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
		lastMessage := &models.Message{}
		if len(session.Messages) > 0 {
			lastMessage = session.Messages[len(session.Messages) - 1]
		}
		waitingUser = append(waitingUser, &resource.WaitingChatSession{
			Username:     session.User.GetUsername(),
			Avatar:       session.User.GetAvatarUrl(),
			UserId:           session.User.GetPrimaryKey(),
			MessageCount: len(session.Messages),
			Description:  "",
			Messages: msgs,
			LastTime: lastMessage.ReceivedAT,
			SessionId: session.Id,
		})
	}
	sort.Slice(waitingUser, func(i, j int) bool {
		return waitingUser[i].LastTime > waitingUser[j].LastTime
	})
	adminConns := m.GetAllConn()
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
func (m *adminManager) BroadcastUserTransfer(adminId int64)   {
	client, exist := m.GetConn(adminId)
	if exist {
		transfers := transferRepo.Get([]Where{
			{
				Filed: "to_admin_id = ?",
				Value: adminId,
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
	var serviceUsers []*models.Admin
	admins := adminRepo.Get([]Where{}, -1, []string{})
	conns := m.GetAllConn()
	data := make([]resource.Admin, 0, len(serviceUsers))
	for _, admin := range admins {
		_, online := m.GetConn(admin.GetPrimaryKey())
		if online {
			data = append(data, resource.Admin{
				Avatar:           admin.GetAvatarUrl(),
				Username:         admin.Username,
				Online:           online,
				Id:               admin.GetPrimaryKey(),
				AcceptedCount: chat.AdminService.GetActiveCount(admin.GetPrimaryKey()),
			})
		}
	}
	m.SendAction(NewAdminsAction(data), conns...)
}
