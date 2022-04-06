package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/time/rate"
	"sync"
	"time"
	"unicode/utf8"
	"ws/app/contract"
	"ws/app/exceptions"
	"ws/app/log"
	"ws/app/models"
	"ws/app/repositories"
)

func NewConn(user contract.User, conn *websocket.Conn, manager ConnManager) Conn {
	return &Client{
		conn:        conn,
		closeSignal: make(chan interface{}),
		send:        make(chan *Action, 100),
		manager:     manager,
		User:        user,
		uuid:        uuid.NewV4().String(),
		limiter:     rate.NewLimiter(5, 10),
	}
}

type Conn interface {
	readMsg()
	sendMsg()
	close()
	run()
	Deliver(action *Action)
	GetUserId() int64
	GetUser() contract.User
	GetUuid() string
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
	uuid    string
	Created int64
	limiter *rate.Limiter
}

func (c *Client) GetCreateTime() int64 {
	return c.Created
}

// GetGroupId 分组Id
func (c *Client) GetGroupId() int64 {
	return c.User.GetGroupId()
}

// GetUuid 每个连接的unique id
func (c *Client) GetUuid() string {
	return c.uuid
}

func (c *Client) GetUser() contract.User {
	return c.User
}

func (c *Client) GetUserId() int64 {
	return c.User.GetPrimaryKey()
}

func (c *Client) run() {
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

// 发送消息验证
func (c *Client) validate(data map[string]interface{}) error {
	if !c.limiter.Allow() {
		c.send <- NewErrorMessage("发送过于频繁，请慢一些")
	}
	content, exist := data["content"]
	if exist {
		s, ok := content.(string)
		if ok {
			length := utf8.RuneCountInString(s)
			if length == 0 {
				return errors.New("请勿发送空内容")
			}
			if length > 512 {
				return errors.New("内容长度必须小于512个字符")
			}
		}
	}
	return nil
}

// 从websocket读消息
func (c *Client) readMsg() {
	var msg = make(chan []byte, 50)
	for {
		go func() {
			_, message, err := c.conn.ReadMessage()
			// 读消息失败说明连接异常，调用close方法
			if err != nil {
				exceptions.Handler(err)
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
				data, ok := act.Data.(map[string]interface{})
				if ok {
					err = c.validate(data)
					if err != nil {
						c.send <- NewErrorMessage(err.Error())
					} else {
						log.Log.WithField("a-type", "websocket").
							WithField("b-type", c.manager.GetTypes()).
							WithField("c-type", "read-message").
							Infof("<user-id:%d><action:%s> %s",
								c.GetUserId(),
								act.Action,
								msgStr)
						c.manager.ReceiveMessage(&ConnMessage{
							Action: act,
							Conn:   c,
						})
					}
				}
			} else {
				exceptions.Handler(err)
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
			if err == nil {
				err := c.conn.WriteMessage(websocket.TextMessage, msgStr)
				log.Log.WithField("a-type", "websocket").
					WithField("b-type", c.manager.GetTypes()).
					WithField("b-type", "send-message").
					Infof("<user-id:%d><action:%s> %s",
						c.GetUserId(),
						act.Action,
						msgStr)
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
					exceptions.Handler(err)
					// 发送失败，close
					c.close()
					return
				}
			} else {
				exceptions.Handler(err)
			}
		case <-c.closeSignal:
			return
		}
	}
}
