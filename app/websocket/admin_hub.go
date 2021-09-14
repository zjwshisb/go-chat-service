package websocket

import (
	"sort"
	"time"
	"ws/app/chat"
	"ws/app/models"
)

type adminHub struct {
	BaseHub
}

func (hub *adminHub) Run() {
	hub.Register(UserLogin, func(i ...interface{}) {
		hub.BroadcastAdmins()
		hub.BroadcastWaitingUser()
		if len(i) > 0 {
			ii := i[0]
			if client, ok := ii.(Conn); ok {
				hub.BroadcastUserTransfer(client.GetUserId())
			}
		}
	})
	hub.Register(UserLogout, func(i ...interface{}) {
		if len(i) > 0 {
			ii := i[0]
			if client, ok := ii.(*AdminConn); ok {
				admin := client.User
				adminSetting := admin.GetSetting()
				adminSetting.LastOnline  = time.Now()
				adminRepo.SaveSetting(adminSetting)
			}
		}
		hub.BroadcastAdmins()
	})
	hub.BaseHub.Run()
}
// 广播待接入用户
func (hub *adminHub) BroadcastWaitingUser() {
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
	waitingUser :=  make([]*models.WaitingChatSessionJson, 0, len(sessions))
	for _, session := range sessions {
		userMap[session.UserId] = session.User
		lastMessage := &models.Message{}
		if len(session.Messages) > 0 {
			lastMessage = session.Messages[len(session.Messages) - 1]
		}
		waitingUser = append(waitingUser, &models.WaitingChatSessionJson{
			Username:     session.User.GetUsername(),
			Avatar:       session.User.GetAvatarUrl(),
			UserId:           session.User.GetPrimaryKey(),
			MessageCount: len(session.Messages),
			Description:  "",
			LastTime: lastMessage.ReceivedAT,
			LastMessage: lastMessage.Content,
			LastType: lastMessage.Type,
			SessionId: session.Id,
		})
	}
	sort.Slice(waitingUser, func(i, j int) bool {
		return waitingUser[i].LastTime > waitingUser[j].LastTime
	})
	conns := hub.GetAllConn()
	for _, adminConnI := range conns {
		adminConn, ok := adminConnI.(*AdminConn)
		if ok {
			adminUserSlice := make([]*models.WaitingChatSessionJson, 0)
			for _, userJson := range waitingUser {
				u := userMap[userJson.UserId]
				if adminConn.User.AccessTo(u) {
					adminUserSlice = append(adminUserSlice, userJson)
				}
			}
			adminConn.Deliver(NewWaitingUsers(adminUserSlice))
		}
	}
}
// 向admin推送待转接入的用户
func (hub *adminHub) BroadcastUserTransfer(adminId int64)   {
	client, exist := hub.GetConn(adminId)
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
		data := make([]*models.ChatTransferJson, 0, len(transfers))
		for _, transfer := range transfers {
			data = append(data, transfer.ToJson())
		}
		client.Deliver(NewUserTransfer(data))
	}
}
// 广播在线admin
func (hub *adminHub) BroadcastAdmins() {
	var serviceUsers []*models.Admin
	admins := adminRepo.Get([]Where{}, -1, []string{})
	conns := hub.GetAllConn()
	data := make([]models.AdminJson, 0, len(serviceUsers))
	for _, admin := range admins {
		_, online := hub.GetConn(admin.GetPrimaryKey())
		if online {
			data = append(data, models.AdminJson{
				Avatar:           admin.GetAvatarUrl(),
				Username:         admin.Username,
				Online:           online,
				Id:               admin.GetPrimaryKey(),
				AcceptedCount: chat.GetAdminUserActiveCount(admin.GetPrimaryKey()),
			})
		}
	}
	hub.SendAction(NewServiceUserAction(data), conns...)
}
