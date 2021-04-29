package hub

import (
	"github.com/gorilla/websocket"
	"time"
	"ws/action"
	"ws/db"
	"ws/models"
)

type UserConn struct {
	BaseConn
	User *models.User
	ServiceId int64
	CreatedAt int64
}
func (c *UserConn) GetUserId() int64 {
	return c.User.ID
}
func (c *UserConn) Setup() {
	sid := c.User.GetLastServerId()
	if sid > 0 {
		c.ServiceId = sid
	}
	c.Register(onReceiveMessage, func(i ...interface{}) {
		length := len(i)
		if length >= 1 {
			ai := i[0]
			act, ok := ai.(*action.Action)
			if ok {
				switch act.Action {
				case action.SendMessageAction:
					msg, err := act.GetMessage()
					if err == nil {
						if len(msg.Content) != 0 {
							msg.IsServer = false
							msg.UserId = c.User.ID
							msg.ReceivedAT = time.Now().Unix()
							c.Deliver(action.NewReceiptAction(msg))
							msg.ServiceId = c.ServiceId
							db.Db.Save(msg)
							if c.ServiceId > 0 {
								serviceClient, exist := ServiceHub.GetConn(c.ServiceId)
								if exist {
									serviceClient.Deliver(action.NewReceiveAction(msg))
								}
							}
						}
					}
					break
				}
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