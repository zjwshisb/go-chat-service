package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/app/contract"
	"ws/app/log"
	"ws/app/models"
	"ws/app/repositories"
)

type Conn interface {
	readMsg()
	sendMsg()
	close()
	run()
	Deliver(action *Action)
	GetUserId() int64
	GetUser() contract.User
	GetUid() string
	GetGroupId() int64
	GetCreateTime() int64
}

type Client struct {
	conn        *websocket.Conn
	closeSignal chan interface{} // 连接断开后的广播通道，用于中断readMsg,sendMsg goroutine
	send        chan *Action     // 发送的消息chan
	sync.Once
	manager ConnManager
	User    contract.User
	uid     string
	Created int64
}

func (c *Client) GetCreateTime() int64 {
	return c.Created
}

// GetGroupId 分组Id
func (c *Client) GetGroupId() int64 {
	return c.User.GetGroupId()
}

// GetUid 每个连接的unique id
func (c *Client) GetUid() string {
	return c.uid
}

func (c *Client) GetUser() contract.User {
	return c.User
}

func (c *Client) GetUserId() int64 {
	return c.User.GetPrimaryKey()
}

func (c *Client) run() {
	c.Created = time.Now().Unix()
	go c.readMsg()
	go c.sendMsg()
}

//幂等的close方法 关闭连接，相关清理
func (c *Client) close() {
	c.Once.Do(func() {
		close(c.closeSignal)
		_ = c.conn.Close()
		c.manager.Unregister(c)
	})
}

// 从websocket读消息
func (c *Client) readMsg() {
	var msg = make(chan []byte, 50)
	for {
		go func() {
			_, message, err := c.conn.ReadMessage()
			// 读消息失败说明连接异常，调用close方法
			if err != nil {
				c.close()
			} else {
				msg <- message
			}
		}()
		select {
		case <-c.closeSignal:
			return
		case msgStr := <-msg:
			var act = &Action{}
			err := act.UnMarshal(msgStr)
			if err == nil {
				log.Log.Info(act)
				c.manager.ReceiveMessage(&ConnMessage{
					Action: act,
					Conn:   c,
				})
			} else {
				log.Log.Error(err)
			}
		}
	}
}

// Deliver 投递消息
func (c *Client) Deliver(act *Action) {
	c.send <- act
}

// 向websocket发消息
func (c *Client) sendMsg() {
	for {
		select {
		case act := <-c.send:
			msgStr, err := act.Marshal()
			log.Log.Info(act)
			if err == nil {
				err := c.conn.WriteMessage(websocket.TextMessage, msgStr)
				if err == nil {
					switch act.Action {
					case MoreThanOne:
						c.close()
					case OtherLogin:
						c.close()
					case ReceiveMessageAction:
						msg, ok := act.Data.(*models.Message)
						if ok {
							msg.SendAt = time.Now().Unix()
							repositories.MessageRepo.Save(msg)
						}
					default:
					}
				} else {
					// 发送失败，close
					c.close()
					return
				}
			} else {
				log.Log.Error(err)
			}
		case <-c.closeSignal:
			return
		}
	}
}
