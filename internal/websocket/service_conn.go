package websocket

import (
	"github.com/gorilla/websocket"
	"time"
	"ws/internal/action"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/models"
)

type ServiceConn struct {
	User *models.BackendUser
	BaseConn
}
func (c *ServiceConn) onReceiveMessage(act *action.Action)  {
	switch act.Action {
	// 客服发送消息给用户
	case action.SendMessageAction:
		msg, err := act.GetMessage()
		if err == nil {
			if msg.UserId > 0 && len(msg.Content) != 0 && chat.CheckUserIdLegal(msg.UserId, c.User.GetPrimaryKey()) {
				msg.ServiceId = c.User.ID
				msg.Source = models.SourceBackendUser
				msg.ReceivedAT = time.Now().Unix()
				msg.Avatar = c.User.GetAvatarUrl()
				databases.Db.Save(msg)

				addTime := chat.GetUserSessionSecond()
				record := &models.QueryRecord{}
				databases.Db.Where("user_id = ?" , msg.UserId).
					Where("service_id = ?", msg.ServiceId).
					Order("id desc").First(record)
				if record.Id > 0 {
					record.BrokeAt = time.Now().Unix() + addTime
					databases.Db.Save(record)
				}
				_ = chat.UpdateUserServerId(msg.UserId, c.User.GetPrimaryKey(), addTime)
				c.Deliver(action.NewReceiptAction(msg))
				userConn, ok := UserHub.GetConn(msg.UserId)
				if ok { // 在线
					userConn.Deliver(action.NewReceiveAction(msg))
				}
			}
		}
		break
	}
}
func (c *ServiceConn) onSendSuccess(act *action.Action) {
	switch act.Action {
	case action.MoreThanOne:
		c.close()
		break
	case action.OtherLogin:
		c.close()
	}
}
func (c *ServiceConn) Setup() {
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
	c.Register(onSendSuccess, func(i ...interface{}) {
		length := len(i)
		if length >= 1 {
			ai := i[0]
			act, ok := ai.(*action.Action)
			if ok {
				c.onReceiveMessage(act)
			}
		}
	})
}
func (c *ServiceConn) GetUserId() int64 {
	return c.User.ID
}

func NewServiceConn(user *models.BackendUser, conn *websocket.Conn) *ServiceConn {
	return &ServiceConn{
		User: user,
		BaseConn: BaseConn{
			conn:        conn,
			closeSignal: make(chan interface{}),
			send:        make(chan *action.Action, 100),
		},
	}
}
