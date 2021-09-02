package websocket

import (
	"sort"
	"time"
	"ws/app/chat"
	"ws/app/databases"
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
				adminSetting := &models.AdminChatSetting{}
				databases.Db.Model(admin).Association("Setting").Find(adminSetting)
				if adminSetting.Id > 0 {
					adminSetting.LastOnline  = time.Now()
					databases.Db.Save(adminSetting)
				}
			}
		}
		hub.BroadcastAdmins()
	})
	hub.BaseHub.Run()
}
// 广播待接入用户
func (hub *adminHub) BroadcastWaitingUser() {
	manualUid := chat.GetManualUserIds()
	users := make([]models.User, 0)
	databases.Db.Where("id in ?", manualUid).
		Find(&users)
	repo := &repositories.MessageRepo{}
	messages := repo.GetUnSend([]*repositories.Where{
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
		transfers := make([]*models.ChatTransfer, 0)
		databases.Db.Where("to_admin_id = ?", adminId).
			Where("is_accepted = ?", 0).
			Where("is_canceled = ?", 0).
			Order("id desc").
			Preload("FromAdmin").
			Preload("User").
			Find(&transfers)
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
	databases.Db.Find(&serviceUsers)
	conns := hub.GetAllConn()
	data := make([]models.AdminJson, 0, len(serviceUsers))
	for _, serviceUser := range serviceUsers {
		_, online := hub.GetConn(serviceUser.ID)
		if online {
			data = append(data, models.AdminJson{
				Avatar:           serviceUser.GetAvatarUrl(),
				Username:         serviceUser.Username,
				Online:           online,
				Id:               serviceUser.GetPrimaryKey(),
			})
		}
	}
	hub.SendAction(NewServiceUserAction(data), conns...)
}
