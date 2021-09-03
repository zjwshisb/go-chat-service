package websocket

import (
	"sort"
	"time"
	"ws/app/chat"
	"ws/app/models"
	"ws/app/repositories"
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
	manualUid := chat.GetManualUserIds()
	users := userRepo.Get([]Where{
		{
			Filed: "id in ?",
			Value: manualUid,
		},
	}, -1, []string{} )
	messages := messageRepo.GetUnSend([]*repositories.Where{
		{
			Filed: "user_id in ?",
			Value: manualUid,
		},
	})
	waitingUserMap := make(map[int64]*models.WaitingUserJson)
	for _, user := range users {
		waitingUserMap[user.GetPrimaryKey()] = &models.WaitingUserJson{
			Username:     user.GetUsername(),
			Avatar:       user.GetAvatarUrl(),
			Id:           user.GetPrimaryKey(),
			MessageCount: 0,
			Description:  "",
		}
	}
	for _, message := range messages {
		if wU, exist := waitingUserMap[message.UserId]; exist {
			if wU.LastTime == 0 {
				wU.LastTime = message.ReceivedAT
				wU.LastMessage = message.Content
				wU.MessageCount += 1
				wU.LastType = message.Type
			} else {
				wU.MessageCount += 1
			}
		}
	}
	waitingUserSlice := make([]*models.WaitingUserJson, 0, len(waitingUserMap))
	for _, user := range waitingUserMap {
		waitingUserSlice = append(waitingUserSlice, user)
	}
	sort.Slice(waitingUserSlice, func(i, j int) bool {
		return waitingUserSlice[i].LastTime > waitingUserSlice[j].LastTime
	})
	conns := hub.GetAllConn()
	for _, adminConnI := range conns {
		adminConn, ok := adminConnI.(*AdminConn)
		if ok {
			adminUserSlice := make([]*models.WaitingUserJson, 0)
			for _, userJson := range waitingUserSlice {
				if adminConn.User.AccessTo(userJson.Id) {
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
			})
		}
	}
	hub.SendAction(NewServiceUserAction(data), conns...)
}
