package websocket

import (
	"github.com/gorilla/websocket"
	"strings"
	"time"
	"ws/internal/action"
	"ws/internal/auth"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/models"
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
			onlineServerCount := len(ServiceHub.Clients)
			// 没有客服在线时
			if onlineServerCount == 0 {
				otherRule := models.AutoRule{}
				query := databases.Db.Where("is_system", 1).
					Where("match", models.MatchServiceAllOffLine).Preload("Message").First(&rule)
				if query.RowsAffected > 0 {
					msg := otherRule.GetMessages(c.User.GetPrimaryKey())
					switch otherRule.ReplyType {
					case models.ReplyTypeTransfer:
						UserHub.addToManual(c.GetUserId())
					case models.ReplyTypeMessage:
						if msg != nil {
							databases.Db.Save(msg)
							otherRule.Count++
							databases.Db.Save(&otherRule)
							c.Deliver(action.NewReceiveAction(msg))
						}
					}
				} else {
					UserHub.addToManual(c.GetUserId())
					break LOOP
				}
			} else {
				UserHub.addToManual(c.GetUserId())
				break LOOP
			}
		// 回复消息
		case models.ReplyTypeMessage:
			msg := rule.GetMessages(c.User.GetPrimaryKey())
			if msg != nil {
				databases.Db.Save(msg)
				c.Deliver(action.NewReceiveAction(msg))
			}
			break LOOP
		}
		rule.Count++
		databases.Db.Save(rule)
	}
}
func (c *UserConn) onReceiveMessage(act *action.Action) {
	switch act.Action {
	case action.SendMessageAction:
		msg, err := act.GetMessage()
		if err == nil {
			if len(msg.Content) != 0 {
				msg.Source = models.SourceUser
				msg.UserId = c.GetUserId()
				msg.ReceivedAT = time.Now().Unix()
				msg.Avatar = c.User.GetAvatarUrl()
				msg.ServiceId = chat.GetUserLastServerId(c.User.GetPrimaryKey())
				databases.Db.Save(msg)
				c.Deliver(action.NewReceiptAction(msg))
				// 有对应的客服对象
				if msg.ServiceId > 0 {
					serviceClient, exist := ServiceHub.GetConn(msg.ServiceId)
					if exist {
						serviceClient.Deliver(action.NewReceiveAction(msg))
					}
				} else {
					if chat.IsInManual(c.GetUserId()) {
						ServiceHub.BroadcastWaitingUser()
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
		rule := models.AutoRule{}
		query := databases.Db.
			Where("is_system", 1).
			Where("match", models.MatchEnter).
			Preload("Message").
			First(&rule)
		if query.RowsAffected > 0 {
			if rule.Message != nil {
				msg := rule.GetMessages(c.User.GetPrimaryKey())
				if msg != nil {
					databases.Db.Save(msg)
					rule.Count++
					databases.Db.Save(&rule)
					c.Deliver(action.NewReceiveAction(msg))
				}
			}
		}
	})
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
func NewUserConn(user *models.User, conn *websocket.Conn) *UserConn {
	return &UserConn{
		User: user,
		BaseConn: BaseConn{
			conn:        conn,
			closeSignal: make(chan interface{}),
			send:        make(chan *action.Action, 100),
		},
	}
}
