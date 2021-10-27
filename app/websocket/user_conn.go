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
				msg.AdminId = chat.UserService.GetValidAdmin(c.GetUserId())
				// 发送回执
				c.Deliver(NewReceiptAction(msg))
				// 有对应的客服对象
				if msg.AdminId > 0 {
					// 更新会话有效期
					session := chat.SessionService.Get(c.GetUserId(), msg.AdminId)
					if session == nil {
						return
					}
					addTime := chat.SettingService.GetServiceSessionSecond()
					_ = chat.AdminService.UpdateUser(msg.AdminId,msg.UserId, addTime)
					msg.SessionId = session.Id
					messageRepo.Save(msg)
					adminConn, exist := AdminHub.GetConn(msg.AdminId)
					// 客服在线
					if exist {
						UserHub.triggerMessageEvent(models.SceneAdminOnline, msg, session)
						adminConn.Deliver(NewReceiveAction(msg))
					} else { // 客服不在线
						admin := adminRepo.First([]Where{
							{
								Filed: "id = ?",
								Value: msg.AdminId,
							},
						})
						UserHub.triggerMessageEvent(models.SceneAdminOffline, msg, session)
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
							duration := chat.SettingService.GetOfflineDuration()
							if (lastOnline.Unix() + duration) < time.Now().Unix() {
								chat.SessionService.Close(session, false, true)
								noticeMessage := admin.GetBreakMessage(c.GetUserId(), session.Id)
								c.Deliver(NewReceiveAction(noticeMessage))
							}
						}
					}
				} else { // 没有客服对象
					if chat.TransferService.GetUserTransferId(c.GetUserId()) == 0 {
						if chat.ManualService.IsIn(c.GetUserId()) {
							session := chat.SessionService.Get(c.GetUserId(), 0)
							if session != nil {
								msg.SessionId = session.Id
							}
							messageRepo.Save(msg)
							AdminHub.BroadcastWaitingUser()
						} else {
							if chat.SettingService.GetIsAutoTransferManual() { // 自动转人工
								session := UserHub.addToManual(c.GetUserId())
								if session != nil {
									msg.SessionId = session.Id
								}
								messageRepo.Save(msg)
								AdminHub.BroadcastWaitingUser()
							} else {
								messageRepo.Save(msg)
								UserHub.triggerMessageEvent(models.SceneNotAccepted, msg, nil)
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
		if chat.UserService.GetValidAdmin(c.GetUserId()) == 0 && !chat.ManualService.IsIn(c.GetUserId()) {
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
