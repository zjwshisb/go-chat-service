package websocket

import (
	"github.com/gorilla/websocket"
	"strings"
	"time"
	"ws/app/auth"
	"ws/app/chat"
	"ws/app/databases"
	"ws/app/models"
)

type UserConn struct {
	BaseConn
	User      auth.User
	CreatedAt int64
}

func (c *UserConn) GetUserId() int64 {
	return c.User.GetPrimaryKey()
}

func (c *UserConn) autoReply(content string) {
	rules := make([]*models.AutoRule, 0)
	databases.Db.
		Where("is_system", 0).
		Where("is_open", 1).
		Order("sort").
		Preload("Message").
		Find(&rules)
LOOP:
	for _, rule := range rules {
		switch rule.MatchType {
		case models.MatchTypeAll:
			if rule.Match != content {
				continue
			}
		case models.MatchTypePart:
			if !strings.Contains(content, rule.Match) {
				continue
			}
		default:
			continue
		}
		switch rule.ReplyType {
		// 转接人工客服
		case models.ReplyTypeTransfer:
			onlineServerCount := len(AdminHub.Clients)
			// 没有客服在线时
			if onlineServerCount == 0 {
				otherRule := models.AutoRule{}
				query := databases.Db.Where("is_system", 1).
					Where("match", models.MatchServiceAllOffLine).Preload("Message").First(&rule)
				if query.RowsAffected > 0 {
					msg := otherRule.GetReplyMessage(c.User.GetPrimaryKey())
					switch otherRule.ReplyType {
					case models.ReplyTypeTransfer:
						UserHub.addToManual(c.GetUserId())
					case models.ReplyTypeMessage:
						if msg != nil {
							databases.Db.Save(msg)
							otherRule.Count++
							databases.Db.Save(&otherRule)
							c.Deliver(NewReceiveAction(msg))
						}
					}
				} else {
					UserHub.addToManual(c.GetUserId())
				}
			} else {
				UserHub.addToManual(c.GetUserId())
			}
		// 回复消息
		case models.ReplyTypeMessage:
			msg := rule.GetReplyMessage(c.User.GetPrimaryKey())
			if msg != nil {
				databases.Db.Save(msg)
				c.Deliver(NewReceiveAction(msg))
			}
		case models.ReplyTypeEvent:
			c.triggerEvent(rule.Key)
		}
		rule.Count++
		databases.Db.Save(rule)
		break LOOP
	}
}

func (c *UserConn) triggerEvent(key string)  {
	// todo
}

func (c *UserConn) onReceiveMessage(act *Action) {
	switch act.Action {
	case SendMessageAction:
		msg, err := act.GetMessage()
		if err == nil {
			if len(msg.Content) != 0 {
				msg.Source = models.SourceUser
				msg.UserId = c.GetUserId()
				msg.ReceivedAT = time.Now().Unix()
				msg.Avatar = c.User.GetAvatarUrl()
				msg.AdminId = chat.GetUserLastAdminId(c.GetUserId())
				c.Deliver(NewReceiptAction(msg))
				// 有对应的客服对象
				if msg.AdminId > 0 {
					// 更新会话有效期
					session := chat.GetSession(c.GetUserId(), msg.AdminId)
					if session == nil {
						return
					}
					addTime := chat.GetServiceSessionSecond()
					_ = chat.UpdateUserAdminId(msg.UserId, msg.AdminId, addTime)
					msg.SessionId = session.Id
					databases.Db.Save(msg)
					session.BrokeAt = time.Now().Unix() + addTime
					databases.Db.Save(session)
					adminConn, exist := AdminHub.GetConn(msg.AdminId)
					if exist {
						adminConn.Deliver(NewReceiveAction(msg))
					}
				} else {
					databases.Db.Save(msg)
					if chat.IsInManual(c.GetUserId()) {
						AdminHub.BroadcastWaitingUser()
					} else {
						isAutoTransfer, exist := chat.Settings[chat.IsAutoTransfer]
						if exist  && isAutoTransfer.GetValue() == "1"{
							if !chat.IsInManual(c.GetUserId()) {
								UserHub.addToManual(c.GetUserId())
							}
						} else {
							if !chat.IsInManual(c.GetUserId()) {
								c.autoReply(msg.Content)
							}
						}
					}

				}
			}
		}
		break
	}
}
func (c *UserConn) Setup() {
	c.Register(onEnter, func(i ...interface{}) {
		if chat.GetUserLastAdminId(c.GetUserId()) == 0 {
			rule := models.AutoRule{}
			query := databases.Db.
				Where("is_system", 1).
				Where("match", models.MatchEnter).
				Preload("Message").
				First(&rule)
			if query.RowsAffected > 0 {
				if rule.Message != nil {
					msg := rule.GetReplyMessage(c.User.GetPrimaryKey())
					if msg != nil {
						databases.Db.Save(msg)
						rule.Count++
						databases.Db.Save(&rule)
						c.Deliver(NewReceiveAction(msg))
					}
				}
			}
		}
	})
	c.Register(onClose, func(i ...interface{}) {
		UserHub.Logout(c)
	})
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
	c.Register(onSendSuccess, func(i ...interface{}) {
	})
}
func NewUserConn(user auth.User, conn *websocket.Conn) *UserConn {
	return &UserConn{
		User: user,
		BaseConn: BaseConn{
			conn:        conn,
			closeSignal: make(chan interface{}),
			send:        make(chan *Action, 100),
		},
	}
}
