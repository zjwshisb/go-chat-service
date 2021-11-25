package websocket

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
	"ws/app/auth"
	"ws/app/log"
	"ws/app/models"
)

const (
	// conn向客户端发送消息成功事件
	onSendSuccess  = iota
	// 客服端连接成功事件
	onEnter
	// conn读取客服端消息事件
	onReceiveMessage
	// 关闭
	onClose
)

type Conn interface {
	readMsg()
	sendMsg()
	close()
	run()
	Deliver(action *Action)
	GetUserId() int64
	GetUser() auth.User
	GetUid() string
	GetGroupId() int64
}

type Client struct {
	conn        *websocket.Conn
	closeSignal chan interface{} // 连接断开后的广播通道，用于中断readMsg,sendMsg goroutine
	send        chan *Action  // 发送的消息chan
	sync.Once
	manager ConnManager
	User auth.User
	uid string

}

func (c *Client) GetGroupId() int64 {
	return c.User.GetGroupId()
}
func (c *Client) GetUid() string  {
	return c.uid
}

func (c *Client) GetUser() auth.User  {
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
				log.Log.Info(string(message))
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
				c.manager.ReceiveMessage(&ConnMessage{
					Action: act,
					Conn: c,
				})
			} else {
				log.Log.Error(err)
			}
		}
	}
}
// 投递消息
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
				log.Log.Info(string(msgStr))
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
							messageRepo.Save(msg)
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


func NewUserConn(user auth.User, conn *websocket.Conn) Conn {
	return &Client{
		conn:        conn,
		closeSignal: make(chan interface{}),
		send:        make(chan *Action, 100),
		manager: UserManager,
		User: user,
		uid: uuid.NewV4().String(),
	}
}

func NewAdminConn(user *models.Admin, conn *websocket.Conn) Conn {
	return &Client{
		conn:        conn,
		closeSignal: make(chan interface{}),
		send:        make(chan *Action, 100),
		manager: AdminManager,
		User: user,
		uid: uuid.NewV4().String(),
	}
}
