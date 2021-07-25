package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"time"
	"ws/configs"
	"ws/internal/action"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/models"
	"ws/internal/repositories"
	"ws/internal/wechat"
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
				session := &models.ChatSession{}
				databases.Db.Where("user_id = ?" , msg.UserId).
					Where("service_id = ?", msg.ServiceId).
					Order("id desc").First(session)
				if session.Id <= 0 {
					return
				}
				sessionAddTime := chat.GetUserSessionSecond()
				session.BrokeAt = time.Now().Unix() + sessionAddTime
				databases.Db.Save(session)
				msg.ServiceId = c.User.ID
				msg.Source = models.SourceBackendUser
				msg.ReceivedAT = time.Now().Unix()
				msg.Avatar = c.User.GetAvatarUrl()
				msg.SessionId = session.Id
				databases.Db.Save(msg)
				_ = chat.UpdateUserServerId(msg.UserId, c.User.GetPrimaryKey(), sessionAddTime)
				c.Deliver(action.NewReceiptAction(msg))
				userConn, exist := UserHub.GetConn(msg.UserId)
				if exist { // 在线
					userConn.Deliver(action.NewReceiveAction(msg))
				} else {
					hadSubscribe := chat.IsSubScribe(msg.UserId)
					user, exist := repositories.GetUserById(msg.UserId)
					if hadSubscribe && exist && user.GetMpOpenId() != "" {
						_ = wechat.GetMp().GetSubscribe().Send(&subscribe.Message{
							ToUser:           user.GetMpOpenId(),
							TemplateID:       configs.Wechat.SubscribeTemplateIdOne,
							Page:             "/pages/chat/index",
							Data: map[string]*subscribe.DataItem{
								"thing1": {
									Value: msg.Content,
								},
								"thing2": {
									Value: "客服给你回复了一条消息",
								},
							},
						})
						chat.DelSubScribe(msg.UserId)
					}
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
