package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/db"
	"ws/models"
	"ws/util"
)
const MessageKey = "server:%d:message"
type Client struct {
	Conn        *websocket.Conn
	UserId      int64
	isClose     bool
	once        sync.Once
	Send        chan *util.Action
	closeSignal chan struct{}
}

func (c *Client) close() {
	c.once.Do(func() {
		_ = c.Conn.Close()
		c.closeSignal <- struct{}{}
		c.isClose = true
		H.Logout <- c
	})
}
func (c *Client) sendKey() string  {
	return fmt.Sprintf("server:%d:message", c.UserId)
}
func (c *Client) start() {
	H.Login <- c
	go c.ReadMsg()
	go c.SendMsg()
	go c.getMsg()
}

func (c *Client) handleReadAction(action util.Action) (err error) {
	switch action.Action {
	case "message":
		msg, err := models.NewFromAction(action)
		if err == nil {
			msg.ServiceId = c.UserId
			msg.IsServer = true
			db.Db.Save(msg)
		}
		receipt := util.NewReceiptAction(action)
		c.Send<- receipt
	}
	return
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
		case <-c.closeSignal:
			goto END
		case message := <-msg:
			var action util.Action
			err := action.UnMarshal(message)
			ctx := context.Background()
			if err == nil {
				err = c.handleReadAction(action)
				if err == nil {
					db.Redis.RPush(ctx, fmt.Sprintf("user:%d:message", 1), message)
				}
			}
		}
	}
END:
}
func (c *Client) SendMsg() {
	for {
		select {
		case action := <-c.Send:
			msg, err := action.Marshal()
			if err == nil {
				err := c.Conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					c.close()
					goto END
				}
			}
		case <-c.closeSignal:
			goto END
		}
	}
END:
}
func (c *Client) getMsg() {
	var msg = make(chan string, 50)
	var errChan = make(chan error)
	for {
		ctx := context.Background()
		go func() {
			val, err := db.Redis.BLPop(ctx, time.Second * 1, fmt.Sprintf(MessageKey, c.UserId)).Result()
			if err != redis.Nil {
				msg <- val[1]
			} else {
				errChan <- err
			}
		}()
		select {
		case <-c.closeSignal:
			goto END
		case <-errChan:
		case m := <-msg:
			action := &util.Action{}
			if err := action.UnMarshal([]byte(m)); err != nil {
				c.Send <- action
			}
		}
	}
END:
}
