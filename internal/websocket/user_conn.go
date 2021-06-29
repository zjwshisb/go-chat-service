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
	"ws/internal/util"
)

type UserConn struct {
	BaseConn
	User auth.User
	CreatedAt int64
}
func (c *UserConn) GetUserId() int64 {
	return c.User.GetPrimaryKey()
}
func (c *UserConn) autoReply(content string) {
	rules := make([]*models.AutoRule, 0)
	databases.Db.Where("is_system", 0).
		Where("is_open", 1).
		Order("sort").
		Preload("Message").
		Find(&rules)
	LOOP: for _, rule := range rules {
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
		case models.ReplyTypeTransfer:
			_ = chat.AddToManual(c.GetUserId())
			ServiceHub.BroadcastWaitingUser()
			break LOOP
		case models.ReplyTypeMessage:
			msg := models.Message{
				UserId:      c.User.GetPrimaryKey(),
				ServiceId:   0,
				Type:        rule.Message.Type,
				Content:     rule.Message.Content,
				ReceivedAT:  time.Now().Unix(),
				SendAt:      0,
				Source:      models.SourceSystem,
				ReqId:       util.CreateReqId(),
				IsRead:      true,
				Avatar: chat.SystemAvatar(),
			}
			databases.Db.Save(&msg)
			c.Deliver(action.NewReceiveAction(&msg))
			break LOOP
		}
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
				databases.Db.Save(msg)
				c.Deliver(action.NewReceiptAction(msg))
				msg.ServiceId = chat.GetUserLastServerId(c.User.GetPrimaryKey())
				if msg.ServiceId > 0 {
					serviceClient, exist := ServiceHub.GetConn(msg.ServiceId)
					if exist {
						serviceClient.Deliver(action.NewReceiveAction(msg))
					}
				} else {
					c.autoReply(msg.Content)
				}
			}
		}
		break
	}
}
func (c *UserConn) Setup() {
	c.Register(onEnter, func(i ...interface{}) {
		rule := models.AutoRule{}
		query := databases.Db.Where("is_system", 1).
			Where("match", "enter").Preload("Message").First(&rule)
		if query.RowsAffected > 0 {
			if rule.Message != nil {
				msg := models.Message{
					UserId:      c.User.GetPrimaryKey(),
					ServiceId:   0,
					Type:        rule.Message.Type,
					Content:     rule.Message.Content,
					ReceivedAT:  time.Now().Unix(),
					SendAt:      0,
					Source:      models.SourceSystem,
					ReqId:       util.CreateReqId(),
					IsRead:      true,
					Avatar: chat.SystemAvatar(),
				}
				databases.Db.Save(&msg)
				c.Deliver(action.NewReceiveAction(&msg))
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
			conn: conn,
			closeSignal: make(chan interface{}),
			send: make(chan *action.Action, 100),
		},
	}
}