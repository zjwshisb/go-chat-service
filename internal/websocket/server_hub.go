package websocket

import (
	"sort"
	"ws/internal/action"
	"ws/internal/databases"
	"ws/internal/models"
	"ws/internal/repositories"
	"ws/internal/json"
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
	messages := repositories.GetUnSendMessage([]repositories.Where{}, []string{"User"})
	waitingUserMap := make(map[int64]*json.WaitingUser)
	for _, message := range messages {
		if message.User.ID != 0{
			if wU, exist := waitingUserMap[message.User.ID]; !exist {
				waitingUserMap[message.User.GetPrimaryKey()] =  &json.WaitingUser{
					Username:     message.User.GetUsername(),
					Avatar:       message.User.GetAvatarUrl(),
					Id:           message.User.GetPrimaryKey(),
					LastMessage:  message.Content,
					LastTime:     message.ReceivedAT,
					MessageCount: 1,
					Description:  "",
				}
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
	hub.SendAction(action.NewWaitingUsers(waitingUserSlice),  conns...)
}

func (hub *serviceHub) BroadcastServiceUser() {
	var serviceUsers []*models.BackendUser
	databases.Db.Find(&serviceUsers)
	conns := hub.GetAllConn()
	data := make([]json.ChatServiceUser, 0, len(serviceUsers))
	for _, serviceUser := range serviceUsers {
		_, online := hub.GetConn(serviceUser.ID)
		data = append(data, json.ChatServiceUser{
			Avatar: serviceUser.GetAvatarUrl(),
			Username: serviceUser.Username,
			Online: online,
			Id: serviceUser.GetPrimaryKey(),
			TodayAcceptCount: serviceUser.GetTodayAcceptCount(),
		})
	}
	hub.SendAction(action.NewServiceUserAction(data), conns...)
}