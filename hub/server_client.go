package hub

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"sync"
	"time"
	"ws/db"
	"ws/models"
)

type Client struct {
	Conn        *websocket.Conn
	UserId      int64
	Once        sync.Once
	Send        chan *models.Action
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
		c.CloseSignal <- struct{}{}
		Hub.Server.Logout(c)
	})
}

func (c *Client) ping(){
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			c.Send<- models.NewPingAction()
		case <-c.CloseSignal:
			ticker.Stop()
			goto END
		}
	}
	END:
}

func (c *Client) serverUserIdsKey() string {
	return fmt.Sprintf("server:%d:user-ids", c.UserId)
}
func (c *Client) getUserList() {
	//ctx := context.Background()
	//cmd := db.Redis.SMembers(ctx, c.serverUserIdsKey())
	//ids := cmd.Val()
}

func (c *Client) accept(uid int64) {
	uClient, ok := Hub.User.getClient(uid)
	if ok {
		if err := uClient.setServed(c.UserId); err == nil {
			uClient.ServerId = c.UserId
			messages := uClient.getWaitingMsg()
			ctx := context.Background()
			db.Redis.SAdd(ctx, c.serverUserIdsKey())
			for _, msg := range messages {
				msg.ServiceId = c.UserId
				db.Db.Save(msg)
				data := make(map[string]interface{})
				mapstructure.Decode(messages, data)
				c.Send<- &models.Action{
					Data: data,
					Action: "message",
				}
			}
		}
	}
}

func (c *Client) handleReadAction(a *models.Action) (err error) {
	switch a.Action {
	case "message":
		msg, err := models.NewFromAction(a)
		if err == nil {
			if msg.UserId > 0 {
				msg.ServiceId = c.UserId
				msg.IsServer = true
				msg.ReceivedAT = time.Now().Unix()
				db.Db.Save(msg)
				a.Message = msg
				receipt := models.NewReceiptAction(a)
				c.Send<- receipt
				UClient, ok := Hub.User.getClient(msg.UserId)
				if ok { // 在线
					UClient.Send<- a
				}
			}

		}
	}
	return
}
func (c *Client) handleSendAction(act models.Action) {
	if act.Message != nil {
		act.Message.SendAt = time.Now().Unix()
		db.Db.Save(act.Message)
	}
}
func (c *Client) ReadMsg() {
	var msg = make(chan []byte, 50)
	for {
		go func() {
			_, message, err := c.Conn.ReadMessage()
			fmt.Println(message)
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
			var act = &models.Action{}
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
					c.handleSendAction(*act)
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
