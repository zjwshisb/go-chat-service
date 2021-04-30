package hub

import (
	"fmt"
	"ws/action"
	"ws/db"
	"ws/models"
	"ws/resources"
)

type serviceHub struct {
	BaseHub
}

func (hub *serviceHub) Setup() {
	hub.Register(UserLogin, func(i ...interface{}) {
		if len(i) >= 1 {
			serviceClient := i[0].(*ServiceConn)
			fmt.Println(serviceClient.User.Username)
		}
		hub.BroadcastServiceUser()
	})
	hub.Register(UserLogout, func(i ...interface{}) {
		if len(i) >= 1 {
			serviceClient := i[0].(*ServiceConn)
			fmt.Println(serviceClient.User.Username)
		}
		hub.BroadcastServiceUser()
	})
}

func (hub *serviceHub) BroadcastServiceUser() {
	var serviceUsers []*models.ServiceUser
	db.Db.Find(&serviceUsers)
	conns := hub.GetAllConn()
	data := make([]resources.ChatServiceUser, 0)
	for _, serviceUser := range serviceUsers {
		_, online := hub.GetConn(serviceUser.ID)
		data = append(data, resources.ChatServiceUser{
			Avatar: serviceUser.GetAvatarUrl(),
			Username: serviceUser.Username,
			Online: online,
			Id: serviceUser.ID,
			TodayAcceptCount: serviceUser.GetTodayAcceptCount(),
		})
	}
	hub.SendAction(action.NewServiceUserAction(data), conns...)
}