package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws/internal/action"
	"ws/internal/databases"
	"ws/internal/event"
	"ws/internal/log"
	"ws/internal/models"
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
	AnonymousConn
	Authentic
}

type BaseConn struct {
	conn        *websocket.Conn
	closeSignal chan interface{}
	send        chan *action.Action
	sync.Once
	event.BaseEvent
}
func (c *BaseConn) run() {
	go c.readMsg()
	go c.sendMsg()
	c.Call(onEnter)
}

func (c *BaseConn) close() {
	c.Once.Do(func() {
		close(c.closeSignal)
		_ = c.conn.Close()
		c.Call(onClose)
	})
}

// 读消息
func (c *BaseConn) readMsg() {
	var msg = make(chan []byte, 50)
	for {
		go func() {
			_, message, err := c.conn.ReadMessage()
			log.Log.Warning(string(message))
			if err != nil {
				log.Log.Error(err)
				c.close()
			} else {
				msg <- message
			}
		}()
		select {
		case <-c.closeSignal:
			goto END
		case msgStr := <-msg:
			var act = &action.Action{}
			err := act.UnMarshal(msgStr)
			if err == nil {
				c.Call(onReceiveMessage, act)
			} else {
				log.Log.Error(err)
			}
		}
	}
END:
}
func (c *BaseConn) Deliver(act *action.Action) {
	c.send <- act
}

// 发消息
func (c *BaseConn) sendMsg() {
	for {
		select {
		case act := <-c.send:
			msgStr, err := act.Marshal()
			if err == nil {
				log.Log.Warning(string(msgStr))
				err := c.conn.WriteMessage(websocket.TextMessage, msgStr)
				if err == nil {
					switch act.Action {
					case action.MoreThanOne:
						c.close()
					case action.OtherLogin:
						c.close()
					case action.ReceiveMessageAction:
						msg, ok := act.Data.(*models.Message)
						if ok {
							msg.SendAt = time.Now().Unix()
							databases.Db.Save(msg)
						}
					default:
					}
					go c.Call(onSendSuccess, act)
				} else {
					log.Log.Error(err)
					c.close()
					goto END
				}
			} else {
				log.Log.Error(err)
			}
		case <-c.closeSignal:
			goto END
		}
	}
END:
}
