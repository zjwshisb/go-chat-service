package user

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/db"
	"ws/modules"
)

type UClient struct {
	Conn *websocket.Conn
	Send chan *modules.Action
	UserId int64
	isClose bool
	once sync.Once
}

func (c *UClient) close() {
	c.once.Do(func() {
		c.isClose = false
		_ = c.Conn.Close()
	})

}

func (c *UClient) ReadMsg() {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			c.close()
			break
		}
		var ctx = context.Background()
		db.Redis.RPush(ctx, fmt.Sprintf("server:%d:message", 1), message)
	}
}
func (c *UClient) GetMsg() {
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
func (c *UClient) SendMsg() {
	for {
		select {
		case action := <-c.Send:
			err := c.Conn.WriteMessage(websocket.TextMessage, action.Marshal())
			if err != nil {
				fmt.Println(err)
				c.close()
				goto END
			}
		}
	}
END:

}