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
	User        *models.ServerUser
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

func (c *Client) SendUserListAction() {
	users := c.User.GetChatUsers()
	// 获取一周内的聊天记录
	last := time.Now().Unix() - 2 * 24 * 60 * 60 * 1000
	var messages []models.Message
	db.Db.Preload("ServerUser").
		Preload("User").
		Where("received_at > ?", last).
		Where("service_id = ?", c.User.ID).
		Find(&messages)
	for _, user := range users {
		for _, m := range messages {
			if m.UserId == user.ID {
				m.IsSuccess = true
				if m.IsServer {
					m.Avatar = m.ServerUser.GetAvatarUrl()
				}
				if !m.IsRead && !m.IsServer{
					user.Unread += 1
				}
				user.Messages = append(user.Messages, m)
			}
		}
		user.Disabled = !c.User.CheckChatUserLegal(user.ID)
		if _, ok := Hub.User.AcceptedClient.GetClient(user.ID); ok {
			user.Online = true
		}
		if _, ok := Hub.User.WaitingClient.GetClient(user.ID); ok {
			user.Online = true
		}
	}
	userAction := action.NewServerUserList(users)
	c.Send <- userAction
}

func (c *Client) close() {
	c.Once.Do(func() {
		_ = c.Conn.Close()
		close(c.CloseSignal)
		Hub.Server.Logout(c)
	})
}

func (c *Client) ping() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			c.Send <- action.NewPing()
		case <-c.CloseSignal:
			ticker.Stop()
			goto END
		}
	}
END:
}

// 接入用户
func (c *Client) Accept(uid int64) (user *models.User, err error) {
	uClient, exist := Hub.User.WaitingClient.GetClient(uid)
	if !exist {
		err = errors.New("用户端已离线")
		return
	}
	if err := uClient.SetServerId(c.User.ID); err == nil {
		Hub.User.Change2accept(uClient)
		_ = c.User.UpdateChatUser(uid)
		Hub.Server.broadcastWaitingUsers()
		user = uClient.User
	}
	return
}
func (c *Client) handleMessage(act *action.Action) {
	switch act.Action {
	case action.SendMessageAction:
		msg, err := act.GetMessage()
		if err == nil {
			if msg.UserId > 0 && len(msg.Content) != 0 && c.User.CheckChatUserLegal(msg.UserId) {
				msg.ServiceId = c.User.ID
				msg.IsServer = true
				msg.ReceivedAT = time.Now().Unix()
				db.Db.Save(msg)
				msg.Avatar = c.User.GetAvatarUrl()
				c.Send <- action.NewReceipt(msg)
				UClient, ok := Hub.User.AcceptedClient.GetClient(msg.UserId)
				if ok { // 在线
					UClient.Send <- action.NewReceiveAction(msg)
				}
			}
		}
		break
	}
	return
}
func (c *Client) onSendSuccess(act action.Action) {
	if act.Action == action.ReceiveMessageAction {
		msg, ok := act.Data.(*models.Message)
		if ok {
			msg.SendAt = time.Now().Unix()
			db.Db.Save(msg)
			_ = c.User.UpdateChatUser(msg.UserId)
		}
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
				c.handleMessage(act)
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
