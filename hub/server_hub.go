package hub

import (
	"errors"
	"fmt"
	"sync"
	"ws/db"
	"ws/models"
	"ws/util"
)

type serverHub struct {
	Clients map[int64]*Client
	Lock sync.RWMutex
	util.Hook
}

func (hub *serverHub) setup() {
	hub.RegisterHook(serverLogin, func(i ...interface{}) {
		hub.broadcastOnlineList()
	})
	hub.RegisterHook(serverLogout, func(i ...interface{}) {
		hub.broadcastOnlineList()
	})
	hub.RegisterHook(userLogin, func(i ...interface{}) {
		uclient, ok := i[0].(*UClient)
		if ok  {
			if uclient.ServerId > 0{
				serverClient, err := hub.GetClient(uclient.ServerId)
				if err == nil {
					serverClient.Send<- models.NewUserOnlineAction(uclient.UserId)
				}
			} else {
				hub.broadcastUserWaitingCount()
			}
		}
	})
	hub.RegisterHook(userLogout, func(i ...interface{}) {
		uclient, ok := i[0].(*UClient)
		if ok {
			if uclient.ServerId > 0 {
				serverClient, err := hub.GetClient(uclient.ServerId)
				if err == nil {
					serverClient.Send<- models.NewUserOfflineAction(uclient.UserId)
				}
			} else {
				hub.broadcastUserWaitingCount()
			}
		}
	})
}
func (hub *serverHub) Logout(c *Client)  {
	hub.Lock.Lock()
	defer func() {
		hub.Lock.Unlock()
		hub.TriggerHook(serverLogin, c)
	}()
	delete(hub.Clients, c.UserId)
}

func (hub *serverHub) Login(c *Client)  {
	hub.Lock.Lock()
	defer func() {
		hub.Lock.Unlock()
		hub.TriggerHook(serverLogout)
	}()
	hub.Clients[c.UserId] = c
	c.Run()
}

func (hub *serverHub) GetClient(id int64) (client *Client,err error) {
	hub.Lock.RLock()
	defer hub.Lock.RUnlock()
	client, ok := hub.Clients[id]
	if !ok {
		err = errors.New("client not exists")
	}
	return
}
// 广播待接入的客户数量
func (hub *serverHub) broadcastUserWaitingCount()  {
	hub.Lock.RLock()
	defer hub.Lock.RUnlock()
	count := len(Hub.User.Waiting)
	fmt.Println(count)
	act := models.NewUserWaitingCountAction(count)
	for _, client := range hub.Clients {
		fmt.Println(client)
		client.Send<- act
	}
}
// 广播在线客服列表
func (hub *serverHub) broadcastOnlineList() {
	defer hub.Lock.RUnlock()
	hub.Lock.RLock()
	if len(hub.Clients) >  0 {
		var ids []int64
		var broadcastData []interface{}
		for _, c := range hub.Clients {
			ids = append(ids, c.UserId)
		}
		var users = make([]models.ServerUser, 100)
		db.Db.Find(&users, ids)
		for _, v := range users {
			broadcastData = append(broadcastData, map[string]interface{}{
				"user_id": v.ID,
				"username": v.Username,
			})
		}
		for _, c := range hub.Clients {
			c.Send<- models.NewServiceOnlineListAction(map[string]interface{}{
				"list": broadcastData,
			})
		}
	}
}


