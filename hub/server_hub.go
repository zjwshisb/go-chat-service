package hub

import (
	"errors"
	"sync"
	"ws/db"
	"ws/models"
)
type serverHub struct {
	Clients map[int64]*Client
	Lock sync.RWMutex
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
func (hub *serverHub) Logout(c *Client)  {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()
	delete(hub.Clients, c.UserId)
	go hub.broadcastOnlineList()
}
func (hub *serverHub) Login(c *Client)  {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()
	hub.Clients[c.UserId] = c
	c.Start()
	go hub.broadcastOnlineList()
}

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


