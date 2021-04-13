package hub

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/action"
	"ws/db"
	"ws/models"
)

type Client struct {
	Conn        *websocket.Conn
	User      *models.ServerUser
	Once        sync.Once
	Send        chan *action.Action
	CloseSignal chan struct{}
}


func (c *Client) Run() {
	c.SendUserListAction()
	go c.ReadMsg()
	go c.SendMsg()
	go c.ping()
}

func (c *Client) SendUserListAction()  {
	users := c.User.GetChatUsers()
	for _, user := range users {
		if user.LastChatTime < time.Now().Unix() - 24 * 60 * 60 * 2 {
			user.Disabled = true
		}
		if _, ok := Hub.User.AcceptedClient.GetClient(user.ID); ok {
			user.Online = true
		}
		if _, ok := Hub.User.WaitingClient.GetClient(user.ID); ok {
			user.Online = true
		}
	}
	userAction := action.NewServerUserList(users)
	c.Send<- userAction
}


func (c *Client) close() {
	c.Once.Do(func() {
		_ = c.Conn.Close()
		c.CloseSignal <- struct{}{}
		Hub.Server.Logout(c)
	})
}

func (c *Client) ping(){
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			c.Send<- action.NewPing()
		case <-c.CloseSignal:
			ticker.Stop()
			goto END
		}
	}
	END:
}

// 接入用户
func (c *Client) Accept(uid int64) (user *models.User, err error){
	uClient, exist := Hub.User.WaitingClient.GetClient(uid)
	if !exist {
		err = errors.New("用户端已离线")
		return
	}
	if err := uClient.SetServerId(c.User.ID); err == nil {
		Hub.User.Change2accept(uClient)
		c.User.UpdateChatUser(uid)
		Hub.Server.broadcastWaitingUsers()
		user = uClient.User
	}
	return
}

func (c *Client) handleReadAction(a *action.Action) (err error) {
	switch a.Action {
	case "message":
		msg, err := a.GetMessage()
		if err == nil {
			if msg.UserId > 0 {
				msg.ServiceId = c.User.ID
				msg.IsServer = true
				msg.ReceivedAT = time.Now().Unix()
				db.Db.Save(msg)
				a.Message = msg
				receipt := action.NewReceipt(a)
				c.Send<- receipt
				UClient, ok := Hub.User.AcceptedClient.GetClient(msg.UserId)
				if ok { // 在线
					UClient.Send<- a
				}
			}

		}
	}
	return
}
func (c *Client) onSendSuccess(act action.Action) {
	if act.Message != nil {
		act.Message.SendAt = time.Now().Unix()
		db.Db.Save(act.Message)
		c.User.UpdateChatUser(act.Message.UserId)
	}
}
func (c *Client) ReadMsg() {
	var msg = make(chan []byte, 50)
	for {
		go func() {
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				c.close()
			} else {
				msg <- message
			}
		}()
		select {
		case <-c.CloseSignal:
			goto END
		case msgStr := <-msg:
			var act = &action.Action{}
			err := act.UnMarshal(msgStr)
			if err == nil {
				err = c.handleReadAction(act)
			}
		}
	}
END:
}
func (c *Client) SendMsg() {
	for {
		select {
		case act := <-c.Send:
			msgStr, err := act.Marshal()
			if err == nil {
				err := c.Conn.WriteMessage(websocket.TextMessage, msgStr)
				if err == nil { // 发送成功
					c.onSendSuccess(*act)
				} else {
					c.close()
					goto END
				}

			}
		case <-c.CloseSignal:
			goto END
		}
	}
END:
}
