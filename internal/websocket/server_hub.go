package websocket

import (
	"sort"
	"ws/internal/action"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/json"
	"ws/internal/models"
	"ws/internal/repositories"
)

type serviceHub struct {
	BaseHub
}

func (hub *serviceHub) Setup() {
	hub.Register(UserLogin, func(i ...interface{}) {
		hub.BroadcastServiceUser()
		hub.BroadcastWaitingUser()
	})
	hub.Register(UserLogout, func(i ...interface{}) {
		hub.BroadcastServiceUser()
	})
}
func (hub *serviceHub) BroadcastWaitingUser() {
	manualUid := chat.GetManualUserIds()
	users := make([]models.User, 0)
	databases.Db.Where("id in ?", manualUid).
		Find(&users)
	messages := repositories.GetUnSendMessage(&repositories.Where{
		Filed: "user_id in ?",
		Value: manualUid,
	}, &repositories.Where{
		Filed: "source = ?",
		Value: models.SourceUser,
	}, )
	waitingUserMap := make(map[int64]*json.WaitingUser)
	for _, user := range users {
		waitingUserMap[user.GetPrimaryKey()] = &json.WaitingUser{
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
	waitingUserSlice := make([]*json.WaitingUser, 0, len(waitingUserMap))
	for _, user := range waitingUserMap {
		waitingUserSlice = append(waitingUserSlice, user)
	}
	sort.Slice(waitingUserSlice, func(i, j int) bool {
		return waitingUserSlice[i].LastTime > waitingUserSlice[j].LastTime
	})
	conns := hub.GetAllConn()
	hub.SendAction(action.NewWaitingUsers(waitingUserSlice), conns...)
}

func (hub *serviceHub) BroadcastServiceUser() {
	var serviceUsers []*models.BackendUser
	databases.Db.Find(&serviceUsers)
	conns := hub.GetAllConn()
	data := make([]json.ChatServiceUser, 0, len(serviceUsers))
	for _, serviceUser := range serviceUsers {
		_, online := hub.GetConn(serviceUser.ID)
		data = append(data, json.ChatServiceUser{
			Avatar:           serviceUser.GetAvatarUrl(),
			Username:         serviceUser.Username,
			Online:           online,
			Id:               serviceUser.GetPrimaryKey(),
			TodayAcceptCount: serviceUser.GetTodayAcceptCount(),
		})
	}
	hub.SendAction(action.NewServiceUserAction(data), conns...)
}
