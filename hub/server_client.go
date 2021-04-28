package hub

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/action"
	"ws/core/log"
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
	go c.ReadMsg()
	go c.SendMsg()
	go c.ping()
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
func (c *Client) Accept(user *models.User) {
	uClient, exist := Hub.User.WaitingClient.GetClient(user.ID)
	if exist { // 如果用户在线
		_ = uClient.SetServerId(c.User.ID)
		Hub.User.Change2accept(uClient)
	}
	_ = user.SetServerId(c.User.ID)
	_ = c.User.UpdateChatUser(user.ID)
}
// 消息处理
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
// 发送成功处理
func (c *Client) handleSendSuccess(act *action.Action) {
	if act.Action == action.ReceiveMessageAction {
		msg, ok := act.Data.(*models.Message)
		if ok {
			msg.SendAt = time.Now().Unix()
			db.Db.Save(msg)
			_ = c.User.UpdateChatUser(msg.UserId)
		}
	}
}
// 读消息
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
			} else {
				log.Log.Warning(err)
			}
		}
	}
END:
}
// 发消息
func (c *Client) SendMsg() {
	for {
		select {
		case act := <-c.Send:
			msgStr, err := act.Marshal()
			if err == nil {
				err := c.Conn.WriteMessage(websocket.TextMessage, msgStr)
				if err == nil { // 发送成功
					c.handleSendSuccess(act)
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
