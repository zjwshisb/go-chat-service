package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"sort"
	"time"
	"ws/app/chat"
	"ws/app/contract"
	"ws/app/log"
	"ws/app/models"
	"ws/app/mq"
	"ws/app/resource"
	"ws/app/util"
)

var AdminManager *adminManager

const (
	AdminConnChannelKey    = "admin:%d:channel"
	AdminManageGroupKey    = "admin:channel:group"
	AdminGroupKeepAliveKey = "admin:group:%d:alive"
)

type adminManager struct {
	manager
}

func NewAdminConn(user *models.Admin, conn *websocket.Conn) Conn {
	return &Client{
		conn:              conn,
		closeSignal:       make(chan interface{}),
		send:              make(chan *Action, 100),
		manager:           AdminManager,
		User:              user,
		uid:               uuid.NewV4().String(),
		groupKeepAliveKey: AdminGroupKeepAliveKey,
	}
}

func SetupAdmin() {
	AdminManager = &adminManager{
		manager: manager{
			groupCount:            10,
			Channel:               util.GetIPs()[0] + ":" + viper.GetString("Http.Port") + "-admin",
			ConnMessages:          make(chan *ConnMessage, 100),
			userChannelCacheKey:   AdminConnChannelKey,
			groupCacheKey:         AdminManageGroupKey,
			connGroupKeepAliveKey: AdminGroupKeepAliveKey,
		},
	}
	AdminManager.onRegister = AdminManager.registerHook
	AdminManager.onUnRegister = AdminManager.unregisterHook
	AdminManager.Run()
}

func (m *adminManager) Run() {
	m.manager.Run()
	go m.handleReceiveMessage()
	if m.isCluster() {
		go m.handleRemoteMessage()
	}
}

// DeliveryMessage
// 投递消息
// 查询admin是否在本机上，是则直接投递
// 查询admin当前channel，如果存在则投递到该channel上
// 最后则说明admin不在线，处理离线逻辑
func (m *adminManager) DeliveryMessage(msg *models.Message, remote bool) {
	adminConn, exist := m.GetConn(msg.GetAdmin())
	if exist { // admin在线且在当前服务上
		UserManager.triggerMessageEvent(models.SceneAdminOnline, msg)
		adminConn.Deliver(NewReceiveAction(msg))
		return
	} else if !remote && m.isCluster() {
		adminChannel := m.getUserChannel(msg.AdminId) // 获取用户所在channel
		if adminChannel != "" {
			_ = m.publish(adminChannel, mq.NewMessagePayload(msg.Id))
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
	admin := adminRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: msg.AdminId,
		},
	}, []string{})
	setting := admin.GetSetting()
	if setting != nil {
		// 发送离线消息
		if setting.OfflineContent != "" {
			offlineMsg := setting.GetOfflineMsg(msg.UserId, msg.SessionId)
			offlineMsg.Admin = admin
			messageRepo.Save(offlineMsg)
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

// 订阅本manger的channel， 处理消息
func (m *adminManager) handleRemoteMessage() {
	subscribe := mq.Mq().Subscribe(m.GetSubscribeChannel())
	defer subscribe.Close()
	for {
		message := subscribe.ReceiveMessage()
		go func() {
			switch message.Get("types").String() {
			case mq.TypeWaitingUser:
				fmt.Println(mq.TypeWaitingUser)
				gid := message.Get("data").Int()
				if gid > 0 {
					m.broadcastWaitingUser(gid)
				}
			case mq.TypeAdmin:
				gid := message.Get("data").Int()
				if gid > 0 {
					m.broadcastAdmins(gid)
				}
			case mq.TypeOtherLogin:
				uid := message.Get("data").Int()
				if uid > 0 {
					user := userRepo.First([]Where{{
						Filed: "id = ?",
						Value: uid,
					}}, []string{})
					if user != nil {
						m.handleRepeatLogin(user, true)
					}
				}

			case mq.TypeTransfer:
				adminId := message.Get("data").Int()
				if adminId > 0 {
					admin := adminRepo.First([]Where{
						{
							Filed: "id = ?",
							Value: adminId,
						},
					}, []string{})
					if admin != nil {
						m.broadcastUserTransfer(admin)
					}
				}
			case mq.TypeMessage:
				mid := message.Get("data").Int()
				msg := messageRepo.First(mid)
				if msg != nil {
					m.DeliveryMessage(msg , true)
				}
			}
		}()
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
				session := chat.SessionService.Get(msg.UserId, conn.GetUserId())
				if session == nil {
					conn.Deliver(NewErrorMessage("无效的用户"))
					return
				}
				msg.GroupId = conn.GetGroupId()
				msg.AdminId = conn.GetUserId()
				msg.Source = models.SourceAdmin
				msg.GroupId = conn.GetGroupId()
				msg.ReceivedAT = time.Now().Unix()
				msg.Admin = conn.User.(*models.Admin)
				msg.SessionId = session.Id
				messageRepo.Save(msg)
				_ = chat.AdminService.UpdateUser(msg.AdminId, msg.UserId)
				// 服务器回执d
				conn.Deliver(NewReceiptAction(msg))
				UserManager.DeliveryMessage(msg, false)
			}
		}
	}
}

func (m *adminManager) registerHook(conn Conn) {
	m.broadcastUserTransfer(conn.GetUser())
	m.PublishAdmins(conn.GetGroupId())
	m.broadcastWaitingUser(conn.GetGroupId())
}

func (m *adminManager) unregisterHook(conn Conn) {
	u := conn.GetUser()
	admin, ok := u.(*models.Admin)
	if ok {
		setting := admin.GetSetting()
		setting.LastOnline = time.Now()
		adminRepo.SaveSetting(setting)
	}
	m.PublishAdmins(conn.GetGroupId())
}

// PublishWaitingUser 推送待接入用户
func (m *adminManager) PublishWaitingUser(groupId int64) {
	if m.isCluster() {
		m.publishToAllChannel(&mq.Payload{
			Types: mq.TypeWaitingUser,
			Data:  groupId,
		})
	} else {
		m.broadcastWaitingUser(groupId)
	}
}
func (m *adminManager) PublishTransfer(admin contract.User) {
	if m.isCluster() {
		m.publishToAllChannel(&mq.Payload{
			Types: mq.TypeTransfer,
			Data:  admin.GetPrimaryKey(),
		})
	} else {
		m.broadcastUserTransfer(admin)
	}
}

// PublishAdmins 推送在线admin
func (m *adminManager) PublishAdmins(gid int64) {
	if m.isCluster() {
		m.publishToAllChannel(&mq.Payload{
			Types: mq.TypeAdmin,
			Data: gid,
		})
	} else {
		m.broadcastAdmins(gid)
	}
}


// 广播待接入用户
func (m *adminManager) broadcastWaitingUser(groupId int64) {
	log.Log.Info("广播待接入用户")
	sessions := sessionRepo.GetWaitHandles()
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

// 向admin推送待转接入的用户
func (m *adminManager) broadcastUserTransfer(admin contract.User) {
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
		}, -1, []string{"FromAdmin", "User"}, []string{"id desc"})
		data := make([]*resource.ChatTransfer, 0, len(transfers))
		for _, transfer := range transfers {
			data = append(data, transfer.ToJson())
		}
		client.Deliver(NewUserTransfer(data))
	}
}

// 广播在线admin
func (m *adminManager) broadcastAdmins(gid int64) {
	ids := m.GetOnlineUserIds(gid)
	admins := adminRepo.Get([]Where{{
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
