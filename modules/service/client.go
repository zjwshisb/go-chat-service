package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/db"
	"ws/modules"
)

type Client struct {
	Conn        *websocket.Conn
	UserId      int64
	isClose     bool
	once        sync.Once
	Send        chan *modules.Action
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

func (c *Client) start() {
	H.Login <- c
	go c.ReadMsg()
	go c.SendMsg()
	go c.getMsg()
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
			var ctx = context.Background()
			var action *modules.Action
			err := action.UnMarshal(message)
			if err == nil {
				// todo
				db.Redis.RPush(ctx, fmt.Sprintf("user:%d:message", 1), message)
			}
		}
	}
END:
}
func (c *Client) SendMsg() {
	for {
		select {
		case action := <-c.Send:
			err := c.Conn.WriteMessage(websocket.TextMessage, action.Marshal())
			if err != nil {
				c.close()
				goto END
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
			val, err := db.Redis.BLPop(ctx, time.Second*1, fmt.Sprintf("server:%d:message", c.UserId)).Result()
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
			action := &modules.Action{}
			if err := action.UnMarshal([]byte(m)); err != nil {
				c.Send <- action
			}
		}
	}
END:
}
