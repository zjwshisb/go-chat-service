package hub

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/action"
	"ws/db"
	"ws/models"
)

type UClient struct {
	Conn *websocket.Conn
	Send chan *action.Action
	User *models.User
	once sync.Once
	ServerId int64
	CloseSignal chan struct{}
	lock sync.RWMutex
	CreatedAt int64
}
func (c *UClient) Setup() {
	sid := c.User.GetLastServerId()
	if sid > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.ServerId = sid
	}
}
func (c *UClient) Run() {
	go c.sendMsg()
	go c.readMsg()
	go c.Ping()
}

func (c *UClient) Ping() {
	timer := time.NewTicker(time.Second)
	for {
		select {
		case <- timer.C:
			c.Send<- action.NewPing()
		case <- c.CloseSignal:
			timer.Stop()
			goto END
		}
	}
END:
}

func (c *UClient) SetServerId(sid int64) (err error) {
	c.lock.Lock()
	err = c.User.SetServerId(sid)
	if err == nil {
		c.ServerId = sid
	}
	c.lock.Unlock()
	return
}

func (c *UClient) close() {
	c.once.Do(func() {
		_ = c.Conn.Close()
		close(c.CloseSignal)
		Hub.User.Logout(c)
	})
}


func (c *UClient) readMsg() {
	for {
		_, msgStr, err := c.Conn.ReadMessage()
		if err != nil {
			c.close()
			break
		}
		var act = &action.Action{}
		err = act.UnMarshal(msgStr)
		if err == nil {
			switch act.Action {
			case "message":
				msg, err := act.GetMessage()
				act.Data["user_id"] = c.User.ID
				if err == nil {
					msg.IsServer = false
					msg.ReceivedAT = time.Now().Unix()
					msg.UserId = c.User.ID
					if c.ServerId == 0 { // 用户没有被客服接入时
						msg.ServiceId = 0
						db.Db.Save(msg)
					} else { // 用户被客服接入
						msg.ServiceId = c.ServerId
						db.Db.Save(msg)
						act.Message = msg
						sClient, ok := Hub.Server.GetClient(c.ServerId)
						if ok {
							sClient.Send<- act
						}
					}
					receipt := action.NewReceipt(act)
					c.Send<- receipt
				}
			}
		}

	}
}

func (c *UClient) sendMsg() {
	for {
		select {
		case act := <-c.Send:
			msg, err := act.Marshal()
			if err == nil {
				err := c.Conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					c.close()
					goto END
				}
				if act.Message != nil {
					act.Message.SendAt = time.Now().Unix()
					db.Db.Save(act.Message)
				}
			}
		case <-c.CloseSignal:
			goto END
		}
	}
END:
}