package websocket

import (
	"sort"
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
		hub.BroadcastServiceUser()
		hub.BroadcastWaitingUser()
	})
	hub.Register(UserLogout, func(i ...interface{}) {
		hub.BroadcastServiceUser()
	})
	hub.BaseHub.Run()
}

func (hub *adminHub) BroadcastWaitingUser() {
	manualUid := chat.GetManualUserIds()
	users := make([]models.User, 0)
	databases.Db.Where("id in ?", manualUid).
		Find(&users)
	messages := repositories.GetUnSendMessage(&repositories.Where{
		Filed: "user_id in ?",
		Value: manualUid,
	} )
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
	hub.SendAction(NewWaitingUsers(waitingUserSlice), conns...)
}

func (hub *adminHub) BroadcastServiceUser() {
	var serviceUsers []*models.Admin
	databases.Db.Find(&serviceUsers)
	conns := hub.GetAllConn()
	data := make([]models.AdminJson, 0, len(serviceUsers))
	for _, serviceUser := range serviceUsers {
		_, online := hub.GetConn(serviceUser.ID)
		data = append(data, models.AdminJson{
			Avatar:           serviceUser.GetAvatarUrl(),
			Username:         serviceUser.Username,
			Online:           online,
			Id:               serviceUser.GetPrimaryKey(),
		})
	}
	hub.SendAction(NewServiceUserAction(data), conns...)
}
