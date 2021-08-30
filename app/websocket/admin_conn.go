package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"time"
	"ws/app/chat"
	"ws/app/databases"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/wechat"
	"ws/configs"
)
// admin websocket
type AdminConn struct {
	User *models.Admin
	BaseConn
}

func (c *AdminConn) UpdateSetting() {
	setting := &models.AdminChatSetting{}
	databases.Db.Model(c.User).Association("Setting").Find(setting)
	c.User.Setting = setting
}

func (c *AdminConn) GetUserId() int64 {
	return c.User.ID
}
// 从conn读取消息后的处理
func (c *AdminConn) onReceiveMessage(act *Action)  {
	switch act.Action {
	// 客服发送消息给用户
	case SendMessageAction:
		msg, err := act.GetMessage()
		if err == nil {
			if msg.UserId > 0 && len(msg.Content) != 0 {
				if !chat.CheckUserIdLegal(msg.UserId, c.User.GetPrimaryKey()) {
					c.Deliver(NewErrorMessage("该用户已失效，无法发送消息"))
					return
				}
				session := chat.GetSession(msg.UserId, c.GetUserId())
				if session == nil {
					c.Deliver(NewErrorMessage("无效的用户"))
					return
				}
				sessionAddTime := chat.GetUserSessionSecond()
				session.BrokeAt = time.Now().Unix() + sessionAddTime
				databases.Db.Save(session)
				msg.AdminId = c.GetUserId()
				msg.Source = models.SourceAdmin
				msg.ReceivedAT = time.Now().Unix()
				msg.Admin = c.User
				msg.SessionId = session.Id
				msg.Save()
				_ = chat.UpdateUserAdminId(msg.UserId, c.User.GetPrimaryKey(), sessionAddTime)
				// 服务器回执
				c.Deliver(NewReceiptAction(msg))
				userConn, exist := UserHub.GetConn(msg.UserId)
				if exist { // 用户在线
					userConn.Deliver(NewReceiveAction(msg))
				} else {  // 用户不在线
					hadSubscribe := chat.IsSubScribe(msg.UserId)
					user, exist := repositories.GetUserById(msg.UserId)
					if hadSubscribe && exist && user.GetMpOpenId() != "" {
						_ = wechat.GetMp().GetSubscribe().Send(&subscribe.Message{
							ToUser:           user.GetMpOpenId(),
							TemplateID:       configs.Wechat.SubscribeTemplateIdOne,
							Page:             configs.Wechat.ChatPath,
							MiniprogramState: "",
							Data: map[string]*subscribe.DataItem{
								"thing1": {
									Value: "请点击卡片查看",
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

func (c *AdminConn) Setup() {
	c.Register(onReceiveMessage, func(i ...interface{}) {
		length := len(i)
		if length >= 1 {
			ai := i[0]
			act, ok := ai.(*Action)
			if ok {
				c.onReceiveMessage(act)
			}
		}
	})
	// 意外断开
	c.Register(onClose, func(i ...interface{}) {
		AdminHub.Logout(c)
	})
}


func NewAdminConn(user *models.Admin, conn *websocket.Conn) *AdminConn {
	return &AdminConn{
		User: user,
		BaseConn: BaseConn{
			conn:        conn,
			closeSignal: make(chan interface{}),
			send:        make(chan *Action, 100),
		},
	}
}
