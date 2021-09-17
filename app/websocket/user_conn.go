package websocket

import (
	"github.com/gorilla/websocket"
	"time"
	"ws/app/auth"
	"ws/app/chat"
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

// 转接到人工客服列表
func (c *UserConn) handleTransferToManual() *models.ChatSession {
	if !chat.IsInManual(c.GetUserId()) {
		onlineServerCount := len(AdminHub.Clients)
		if onlineServerCount == 0 { // 如果没有在线客服
			rule := autoRuleRepo.GetAdminAllOffLine()
			if rule != nil {
				switch rule.ReplyType {
				case models.ReplyTypeMessage:
					msg := rule.GetReplyMessage(c.User.GetPrimaryKey())
					if msg != nil {
						messageRepo.Save(msg)
						rule.AddCount()
						c.Deliver(NewReceiveAction(msg))
						return nil
					}
				default:
				}
			}
		}
		return UserHub.addToManual(c.GetUserId())
	}
	return nil
}
func (c *UserConn) triggerMessageEvent(scene string, message *models.Message, session *models.ChatSession) {
	rules := autoRuleRepo.GetAllActiveNormal()
	if session == nil {
		session = &models.ChatSession{}
	}
	for _, rule := range rules {
		if rule.IsMatch(message.Content) && rule.SceneInclude(scene) {
			switch rule.ReplyType {
			// 转接人工客服
			case models.ReplyTypeTransfer:
				session := c.handleTransferToManual()
				if session != nil {
					message.SessionId = session.Id
					messageRepo.Save(message)
				}
				AdminHub.BroadcastWaitingUser()
			// 回复消息
			case models.ReplyTypeMessage:
				msg := rule.GetReplyMessage(c.User.GetPrimaryKey())
				if msg != nil {
					msg.SessionId = session.Id
					messageRepo.Save(msg)
					c.Deliver(NewReceiveAction(msg))
				}
			//触发事件
			case models.ReplyTypeEvent:
				switch rule.Key {
				case "break":
					adminId := chat.GetUserLastAdminId(c.GetUserId())
					if adminId > 0 {
						_ = chat.RemoveUserAdminId(c.GetUserId(), adminId)
					}
					msg := rule.GetReplyMessage(c.User.GetPrimaryKey())
					if msg != nil {
						msg.SessionId = session.Id
						messageRepo.Save(msg)
						c.Deliver(NewReceiveAction(msg))
					}
				}
			}
			rule.AddCount()
			return
		}
	}
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
				msg.User = c.User.(*models.User)
				msg.AdminId = chat.GetUserLastAdminId(c.GetUserId())
				// 发送回执
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
					messageRepo.Save(msg)
					adminConn, exist := AdminHub.GetConn(msg.AdminId)
					// 客服在线
					if exist {
						c.triggerMessageEvent(models.SceneAdminOnline, msg, session)
						adminConn.Deliver(NewReceiveAction(msg))
					} else { // 客服不在线
						admin := adminRepo.First([]Where{
							{
								Filed: "id = ?",
								Value: msg.AdminId,
							},
						})
						c.triggerMessageEvent(models.SceneAdminOffline, msg, session)
						setting := admin.GetSetting()
						if setting != nil {
							// 发送离线消息
							if setting.OfflineContent != "" {
								offlineMsg := setting.GetOfflineMsg(c.GetUserId(), session.Id)
								offlineMsg.Admin = admin
								c.Deliver(NewReceiveAction(offlineMsg))
							}
							// 判断是否自动断开
							lastOnline := setting.LastOnline
							duration := chat.GetOfflineDuration()
							if (lastOnline.Unix() + duration) < time.Now().Unix() {
								chat.CloseSession(session, false, true)
								noticeMessage := admin.GetBreakMessage(c.GetUserId(), session.Id)
								c.Deliver(NewReceiveAction(noticeMessage))
							}
						}
					}
				} else { // 没有客服对象
					if chat.GetUserTransferId(c.GetUserId()) == 0 {
						if chat.IsInManual(c.GetUserId()) {
							session := chat.GetSession(c.GetUserId(), 0)
							if session != nil {
								msg.SessionId = session.Id
							}
							messageRepo.Save(msg)
							AdminHub.BroadcastWaitingUser()
						} else {
							isAutoTransfer, exist := chat.Settings[chat.IsAutoTransfer]
							if exist && isAutoTransfer.GetValue() == "1" { // 自动转人工
								session := c.handleTransferToManual()
								if session != nil {
									msg.SessionId = session.Id
								}
								messageRepo.Save(msg)
								AdminHub.BroadcastWaitingUser()
							} else {
								messageRepo.Save(msg)
								c.triggerMessageEvent(models.SceneNotAccepted, msg, nil)
							}
						}
					}

				}
			}
		}
	}
}

func (c *UserConn) Setup() {
	c.Register(onEnter, func(i ...interface{}) {
		if chat.GetUserLastAdminId(c.GetUserId()) == 0 {
			rule := autoRuleRepo.GetEnter()
			if rule != nil {
				msg := rule.GetReplyMessage(c.User.GetPrimaryKey())
				if msg != nil {
					messageRepo.Save(msg)
					rule.Count++
					autoRuleRepo.Save(rule)
					c.Deliver(NewReceiveAction(msg))
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
