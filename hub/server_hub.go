package hub

import (
	"sort"
	"sync"
	"ws/action"
	"ws/util"
)

type serverClientMap struct {
	Clients map[int64] *Client
	Lock    sync.RWMutex
}
// 获取客户端
func (m *serverClientMap) GetClient(id int64) (client *Client, ok bool) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	client, ok = m.Clients[id]
	return
}
// 移除客户端
func (m *serverClientMap) RemoveClient(id int64) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	delete(m.Clients, id)
}
// 注册客户端
func (m *serverClientMap) AddClient(client *Client) ()  {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Clients[client.User.ID] = client
}
// 所有客户端
func (m *serverClientMap) getAllClient() []*Client {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	r := make([]*Client, 0)
	for _, c := range m.Clients {
		r = append(r, c)
	}
	return r
}
type serverHub struct {
	util.Hook
	serverClientMap
}
// 初始化
func (hub *serverHub) setup() {
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
// 发送消息
func (hub *serverHub) SendAction(a  *action.Action, clients ...*Client) {
	for _,c := range clients {
		c.Send<- a
	}
}
// 登出客户端
func (hub *serverHub) Logout(c *Client) {
	hub.RemoveClient(c.User.ID)
	hub.TriggerHook(serverLogout, c)
}
// 登入客户端
func (hub *serverHub) Login(c *Client) {
	if old, ok := hub.GetClient(c.User.ID); ok{
		old.close()
	}
	hub.AddClient(c)
	c.Send<- hub.getWaitingUsersAction()
	go c.Run()
	hub.TriggerHook(serverLogin, c)
}
// 通知用户上线
func (hub *serverHub) noticeUserOnline(uClient *UClient) {
	serverClient, ok := hub.GetClient(uClient.ServerId)
	if ok {
		act := action.NewUserOnline(uClient.User.ID)
		hub.SendAction(act, serverClient)
	}
}
// 通知用户下线
func (hub *serverHub) noticeUserOffOnline(uClient *UClient) {
	serverClient, ok := hub.GetClient(uClient.ServerId)
	if ok {
		act := action.NewUserOffline(uClient.User.ID)
		hub.SendAction(act, serverClient)
	}
}
// 待接入用户
func (hub *serverHub) getWaitingUsersAction() *action.Action {
	d := make([]map[string]interface{}, 0)
	s := make([]*UClient, 0)
	for _, c := range Hub.User.WaitingClient.Clients {
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
	act := action.NewWaitingUsers(d)
	return act
}
// 广播待接入的客户数量
func (hub *serverHub) broadcastWaitingUsers() {
	act := hub.getWaitingUsersAction()
	client := hub.getAllClient()
	hub.SendAction(act, client...)
}
