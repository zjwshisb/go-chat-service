package hub

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/db"
	"ws/models"
)
type UClient struct {
	Conn *websocket.Conn
	Send chan *models.Action
	UserId int64
	once sync.Once
	ServerId int64
	CloseSignal chan struct{}
	lock sync.RWMutex
}
func (c *UClient) Setup() {
	ctx := context.Background()
	cmd := db.Redis.Get(ctx, c.CacheServerIdKey())
	if ServerId, err := cmd.Int64(); err == nil {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.ServerId = ServerId
		db.Redis.SetEX(ctx, c.CacheServerIdKey(), ServerId, time.Hour * 24 * 2)
	}
}
func (c *UClient) Run() {
	go c.sendMsg()
	go c.readMsg()
	go c.Ping()
}

func (c *UClient) waitingSendMessageKey() string {
	return fmt.Sprintf("user:%d:waiting-send-messages", c.UserId)
}

func (c *UClient) CacheServerIdKey() string {
	return fmt.Sprintf("user:%d:server", c.UserId)
}

func (c *UClient) Ping() {
	timer := time.NewTicker(time.Second)
	for {
		select {
		case <- timer.C:
			c.Send<- models.NewPingAction()
		case <- c.CloseSignal:
			timer.Stop()
			goto END
		}
	}
END:
}

func (c *UClient) setServed(sid int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.ServerId > 0 {
		return errors.New("user is accept")
	}
	ctx := context.Background()
	cmd := db.Redis.SetEX(ctx, c.CacheServerIdKey(), sid, time.Hour * 24 * 2)
	if cmd.Err() != nil {
		return errors.New("persist error")
	}
	c.ServerId = sid
	return nil
}

func (c *UClient) close() {
	c.once.Do(func() {
		_ = c.Conn.Close()
		close(c.CloseSignal)
		Hub.User.Logout(c)
	})
}
func (c *UClient) getWaitingMsg() (messages []models.Message) {
	db.Db.Where("user_id = ?" , c.UserId).Where("service_id", 0).Find(&messages)
	return
}

func (c *UClient) readMsg() {
	for {
		_, msgStr, err := c.Conn.ReadMessage()
		if err != nil {
			c.close()
			break
		}
		var act = &models.Action{}
		err = act.UnMarshal(msgStr)
		if err == nil {
			switch act.Action {
			case "message":
				msg, err := models.NewFromAction(act)
				if err == nil {
					act.Data["id"] = msg.Id
					msg.ServiceId = c.UserId
					msg.IsServer = false
					msg.ReceivedAT = time.Now().Unix()
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
					receipt := models.NewReceiptAction(act)
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