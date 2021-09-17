package websocket

import (
	"ws/app/chat"
	"ws/app/models"
)

type userHub struct {
	BaseHub
}


func (userHub *userHub) Run() {
	userHub.Register(UserLogin, func(i ...interface{}) {
		length := len(i)
		if length >= 1 {
			ai := i[0]
			conn, ok := ai.(Conn)
			if ok {
				uid := conn.GetUserId()
				sid := chat.GetUserLastAdminId(uid)
				if sid > 0 {
					adminConn, exist := AdminHub.GetConn(sid)
					if exist {
						adminConn.Deliver(NewUserOnline(uid))
					}
				}
			}
		}
	})
	userHub.Register(UserLogout, func(i ...interface{}) {
		length := len(i)
		if length >= 1 {
			ai := i[0]
			conn, ok := ai.(Conn)
			if ok {
				uid := conn.GetUserId()
				sid := chat.GetUserLastAdminId(uid)
				if sid > 0 {
					adminConn, exist := AdminHub.GetConn(sid)
					if exist {
						adminConn.Deliver(NewUserOffline(uid))
					}
				}
			}
		}
	})
	userHub.BaseHub.Run()
}
// 加入人工列表
func (userHub *userHub) addToManual(uid int64) *models.ChatSession {
	if !chat.IsInManual(uid) {
		onlineServerCount := len(AdminHub.Clients)
		if onlineServerCount == 0 { // 如果没有在线客服
			rule := autoRuleRepo.GetAdminAllOffLine()
			if rule != nil {
				switch rule.ReplyType {
				case models.ReplyTypeMessage:
					msg := rule.GetReplyMessage(uid)
					if msg != nil {
						messageRepo.Save(msg)
						rule.AddCount()
						conn, exist := userHub.GetConn(uid)
						if exist {
							conn.Deliver(NewReceiveAction(msg))
						}
						return nil
					}
				default:
				}
			}
		}
		_ = chat.AddToManual(uid)
		AdminHub.BroadcastWaitingUser()
		session := chat.GetSession(uid, 0)
		if session == nil {
			session = chat.CreateSession(uid, models.ChatSessionTypeNormal)
		}
		return session
	}
	return nil

}

// 触发事件
func (userHub *userHub) triggerMessageEvent(scene string, message *models.Message, session *models.ChatSession) {
	rules := autoRuleRepo.GetAllActiveNormal()
	if session == nil {
		session = &models.ChatSession{}
	}
	for _, rule := range rules {
		if rule.IsMatch(message.Content) && rule.SceneInclude(scene) {
			switch rule.ReplyType {
			// 转接人工客服
			case models.ReplyTypeTransfer:
				session := userHub.addToManual(message.UserId)
				if session != nil {
					message.SessionId = session.Id
					messageRepo.Save(message)
				}
				AdminHub.BroadcastWaitingUser()
			// 回复消息
			case models.ReplyTypeMessage:
				msg := rule.GetReplyMessage(message.UserId)
				if msg != nil {
					msg.SessionId = session.Id
					messageRepo.Save(msg)
					conn, exist := userHub.GetConn(msg.UserId)
					if exist {
						conn.Deliver(NewReceiveAction(msg))
					}
				}
			//触发事件
			case models.ReplyTypeEvent:
				switch rule.Key {
				case "break":
					adminId := chat.GetUserLastAdminId(message.UserId)
					if adminId > 0 {
						_ = chat.RemoveUserAdminId(message.UserId, adminId)
					}
					msg := rule.GetReplyMessage(message.UserId)
					if msg != nil {
						msg.SessionId = session.Id
						messageRepo.Save(msg)
						conn, exist := userHub.GetConn(msg.UserId)
						if exist {
							conn.Deliver(NewReceiveAction(msg))
						}
					}
				}
			}
			rule.AddCount()
			return
		}
	}
}
