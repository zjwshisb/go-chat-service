package hub

import (
	"sort"
	"sync"
	"ws/db"
	"ws/models"
	"ws/util"
)

type serverHub struct {
	Clients map[int64]*Client
	Lock    sync.RWMutex
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
		uClient, ok := i[0].(*UClient)
		if ok {
			if uClient.ServerId > 0 {
				hub.noticeUserOnline(uClient)
			} else {
				hub.broadcastWaitingUsers()
			}
		}
	})
	hub.RegisterHook(userLogout, func(i ...interface{}) {
		uClient, ok := i[0].(*UClient)
		if ok {
			if uClient.ServerId > 0 {
				hub.noticeUserOffOnline(uClient)
			} else {
				hub.broadcastWaitingUsers()
			}
		}
	})
}
func (hub *serverHub) Logout(c *Client) {
	hub.Lock.Lock()
	defer func() {
		hub.Lock.Unlock()
		hub.TriggerHook(serverLogin, c)
	}()
	delete(hub.Clients, c.User.ID)
}

func (hub *serverHub) Login(c *Client) {
	hub.Lock.Lock()
	defer func() {
		hub.Lock.Unlock()
		hub.TriggerHook(serverLogout)
	}()
	if old, ok := hub.Clients[c.User.ID]; ok{
		old.close()
	}
	hub.Clients[c.User.ID] = c
	c.Run()
}

func (hub *serverHub) GetClient(id int64) (client *Client, ok bool) {
	hub.Lock.RLock()
	defer hub.Lock.RUnlock()
	client, ok = hub.Clients[id]
	return
}
func (hub *serverHub) noticeUserOnline(uClient *UClient) {
	serverClient, ok := hub.GetClient(uClient.ServerId)
	if ok {
		serverClient.Send <- models.NewUserOnlineAction(uClient.User.ID)
	}
}
func (hub *serverHub) noticeUserOffOnline(uClient *UClient) {
	serverClient, ok := hub.GetClient(uClient.ServerId)
	if ok {
		serverClient.Send <- models.NewUserOfflineAction(uClient.User.ID)
	}
}

// 广播待接入的客户数量
func (hub *serverHub) broadcastWaitingUsers() {
	Hub.User.Lock.RLock()
	defer Hub.User.Lock.RUnlock()
	d := make([]map[string]interface{}, 0)
	s := make([]*UClient, 0)
	for _, c := range Hub.User.Waiting {
		s = append(s, c)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].CreatedAt < s[j].CreatedAt
	})
	for _, client := range s {
		i := make(map[string]interface{})
		i["id"] = client.User.ID
		i["username"] = client.User.Username
		d = append(d, i)
	}
	act := models.NewWaitingUsersAction(d)
	for _, client := range hub.Clients {
		client.Send <- act
	}
}

// 广播在线客服列表
func (hub *serverHub) broadcastOnlineList() {
	defer hub.Lock.RUnlock()
	hub.Lock.RLock()
	if len(hub.Clients) > 0 {
		var ids []int64
		var broadcastData []interface{}
		for _, c := range hub.Clients {
			ids = append(ids, c.User.ID)
		}
		var users = make([]models.ServerUser, 100)
		db.Db.Find(&users, ids)
		for _, v := range users {
			broadcastData = append(broadcastData, map[string]interface{}{
				"user_id":  v.ID,
				"username": v.Username,
			})
		}
		for _, c := range hub.Clients {
			c.Send <- models.NewServiceOnlineListAction(map[string]interface{}{
				"list": broadcastData,
			})
		}
	}
}
