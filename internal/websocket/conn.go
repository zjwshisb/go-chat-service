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
	onnSendSuccess   = iota
	onReceiveMessage
)

type AnonymousConn interface {
	ping()
	readMsg()
	sendMsg()
	close()
	run()
	Deliver(action *action.Action)
}
type Authentic interface {
	GetUserId() int64
}

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
	go c.ping()
}

func (c *BaseConn) ping() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			c.send <- action.NewPing()
		case <-c.closeSignal:
			ticker.Stop()
			goto END
		}
	}
END:
}

func (c *BaseConn) close() {
	c.Once.Do(func() {
		_ = c.conn.Close()
		close(c.closeSignal)
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
					if act.Action == action.ReceiveMessageAction {
						msg, ok := act.Data.(*models.Message)
						if ok {
							msg.SendAt = time.Now().Unix()
							databases.Db.Save(msg)
						}
					}
					go c.Call(onnSendSuccess, act)
				} else {
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
