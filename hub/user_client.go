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
	isClose bool
	once sync.Once
	ServerId int64
	CloseSignal chan struct{}
	lock sync.RWMutex
}

func (c *UClient) Run() {
	ctx := context.Background()
	cmd := db.Redis.Get(ctx, c.CacheServerIdKey())
	if ServerId, err := cmd.Int64(); err == nil {
		c.ServerId = ServerId
		db.Redis.SetEX(ctx, c.CacheServerIdKey(), ServerId, time.Hour * 24 * 2)
	}
	go c.getMsg()
	go c.sendMsg()
	go c.readMsg()
}

func (c *UClient) waitingSendMessageKey() string {
	return fmt.Sprintf("user:%d:waiting-send-messages", c.UserId)
}

func (c *UClient) CacheServerIdKey() string {
	return fmt.Sprintf("user:%d:server", c.UserId)
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
		c.isClose = false
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
		var act models.Action
		err = act.UnMarshal(msgStr)
		if err == nil {
			switch act.Action {
			case "message":
				msg, err := models.NewFromAction(act)
				if err == nil {
					act.Message = msg
					msg.ServiceId = c.UserId
					msg.IsServer = false
					msg.ReceivedAT = time.Now().Unix()
					if c.ServerId == 0 { // 用户没有被客服接入时
						msg.ServiceId = 0
					} else { // 用户被客服接入
						msg.ServiceId = c.ServerId
						sClient, err := Hub.Server.GetClient(c.ServerId)
						if err == nil {
							sClient.Send<- &act
						}
					}
					db.Db.Save(msg)
					receipt := models.NewReceiptAction(act)
					c.Send<- receipt
				}

			}
		}

	}
}
func (c *UClient) getMsg() {
	var ctx = context.Background()
	for {
		val ,err := db.Redis.BLPop(ctx, time.Second * 10, fmt.Sprintf("user:%d:message", c.UserId)).Result()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}
		if c.isClose {
			break
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
			}
		case <-c.CloseSignal:
			goto END
		}
	}
END:
}