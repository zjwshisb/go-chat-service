package websocket

import (
	"github.com/gorilla/websocket"
	"time"
	"ws/internal/action"
	"ws/internal/auth"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/models"
)

type UserConn struct {
	BaseConn
	User auth.User
	CreatedAt int64
}
func (c *UserConn) GetUserId() int64 {
	return c.User.GetPrimaryKey()
}
func (c *UserConn) onReceiveMessage(act *action.Action) {
	switch act.Action {
	case action.SendMessageAction:
		msg, err := act.GetMessage()
		if err == nil {
			if len(msg.Content) != 0 {
				msg.IsServer = false
				msg.UserId = c.GetUserId()
				msg.ReceivedAT = time.Now().Unix()
				c.Deliver(action.NewReceiptAction(msg))
				msg.ServiceId = chat.GetUserLastServerId(c.User.GetPrimaryKey())
				databases.Db.Save(msg)
				if msg.ServiceId > 0 {
					serviceClient, exist := ServiceHub.GetConn(msg.ServiceId)
					if exist {
						serviceClient.Deliver(action.NewReceiveAction(msg))
					}
				} else {
					ServiceHub.BroadcastWaitingUser()
				}
			}
		}
		break
	}
}
func (c *UserConn) Setup() {
	c.Register(onReceiveMessage, func(i ...interface{}) {
		length := len(i)
		if length >= 1 {
			ai := i[0]
			act, ok := ai.(*action.Action)
			if ok {
				c.onReceiveMessage(act)
			}
		}
	})
	c.Register(onnSendSuccess, func(i ...interface{}) {
	})
}
func NewUserConn(user *models.User, conn *websocket.Conn) *UserConn {
	return &UserConn{
		User: user,
		BaseConn: BaseConn{
			conn: conn,
			closeSignal: make(chan interface{}),
			send: make(chan *action.Action, 100),
		},
	}
}