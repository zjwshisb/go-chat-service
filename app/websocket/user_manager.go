package websocket

import (
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"time"
	"ws/app/chat"
	"ws/app/log"
	"ws/app/models"
	"ws/app/mq"
	"ws/app/wechat"
	"ws/configs"
)

type userManager struct {
	manager
}

var UserManager *userManager

func init() {
	UserManager = &userManager{
		manager{
			Clients:      make(map[int64]Conn),
			Channel:      configs.App.Name + "-user",
			ConnMessages: make(chan *ConnMessage, 100),
			userChannelCacheKey: "user:%d:channel",
			groupCacheKey: "user:channel:group",
		},
	}
	UserManager.onRegister = UserManager.registerHook
	UserManager.Run()
}

func (userManager *userManager) Run() {
	userManager.manager.Run()
	go userManager.handleReceiveMessage()
	if userManager.isCluster() {
		go userManager.handleRemoteMessage()
	}
}
// 投递消息
// 查询user是否在本机上，是则直接投递
// 查询user当前channel，如果存在则投递到该channel上
// 最后则说明user不在线，处理相关逻辑
func (userManager *userManager) DeliveryMessage(msg *models.Message) {
	userConn, exist := UserManager.GetConn(msg.UserId)
	if exist {
		userConn.Deliver(NewReceiveAction(msg))
	} else if userManager.isCluster() {
		userChannel := userManager.getUserChannel(msg.UserId)
		if userChannel != "" {
			_ = userManager.publish(userChannel, &mq.Payload{
				Data: msg.Id,
				Types: "message",
			})
		}
	} else {
		userManager.handleOffline(msg)
	}
}
// 从管道接受消息并处理
func (userManager *userManager) handleReceiveMessage() {
	for {
		payload := <- userManager.ConnMessages
		go userManager.handleMessage(payload)
	}
}
// 处理远程消息
func (userManager *userManager) handleRemoteMessage()  {
	sub := mq.Mq().Subscribe(userManager.GetSubscribeChannel())
	for {
		message, err := sub.ReceiveMessage()
		go func() {
			if err == nil {
				switch message.Types {
				case "message":
					mid := message.Data
					msg := messageRepo.First(mid)
					if msg != nil {
						client, exist := userManager.GetConn(msg.UserId)
						if exist {
							client.Deliver(NewReceiveAction(msg))
						} else {
							userManager.handleOffline(msg)
						}
					}
				}
			}
		}()
	}
}
// 处理离线逻辑
func (userManager *userManager) handleOffline(msg *models.Message) {
	hadSubscribe := chat.SubScribeService.IsSet(msg.UserId)
	user := userRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: msg.UserId,
		},
	})
	if hadSubscribe && user != nil && user.GetMpOpenId() != "" {
		err := wechat.GetMp().GetSubscribe().Send(&subscribe.Message{
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
		if err != nil {
			log.Log.Error(err.Error())
		} else {
			chat.SubScribeService.Remove(msg.UserId)
		}
	}
}
// 处理消息
func (userManager *userManager) handleMessage(payload *ConnMessage) {
	act := payload.Action
	conn := payload.Conn
	switch act.Action {
	case SendMessageAction:
		msg, err := act.GetMessage()
		if err == nil {
			if len(msg.Content) != 0 {
				msg.Source = models.SourceUser
				msg.UserId = conn.GetUserId()
				msg.ReceivedAT = time.Now().Unix()
				msg.User = conn.GetUser().(*models.User)
				msg.AdminId = chat.UserService.GetValidAdmin(conn.GetUserId())
				// 发送回执
				conn.Deliver(NewReceiptAction(msg))
				// 有对应的客服对象
				if msg.AdminId > 0 {
					// 更新会话有效期
					session := chat.SessionService.Get(conn.GetUserId(), msg.AdminId)
					if session == nil {
						return
					}
					addTime := chat.SettingService.GetServiceSessionSecond()
					_ = chat.AdminService.UpdateUser(msg.AdminId, msg.UserId, addTime)
					msg.SessionId = session.Id
					messageRepo.Save(msg)
					AdminManager.DeliveryMessage(msg)
				} else { // 没有客服对象
					if chat.TransferService.GetUserTransferId(conn.GetUserId()) == 0 {
						if chat.ManualService.IsIn(conn.GetUserId()) {
							session := chat.SessionService.Get(conn.GetUserId(), 0)
							if session != nil {
								msg.SessionId = session.Id
							}
							messageRepo.Save(msg)
							AdminManager.BroadcastWaitingUser()
						} else {
							if chat.SettingService.GetIsAutoTransferManual() { // 自动转人工
								session := UserManager.addToManual(conn.GetUserId())
								if session != nil {
									msg.SessionId = session.Id
								}
								messageRepo.Save(msg)
								AdminManager.BroadcastWaitingUser()
							} else {
								messageRepo.Save(msg)
								userManager.triggerMessageEvent(models.SceneNotAccepted, msg)
							}
						}
					}
				}
			}
		}
	}

}

func (userManager *userManager) publishWaitingCount()  {
	if userManager.isCluster() {

	} else {
		
	}
}
func (userManager *userManager) BroadcastWaitingCount()  {
	count := chat.ManualService.GetTotalCount()
	action := NewWaitingUserCount(count)
	conns := userManager.GetAllConn()
	userManager.SendAction(action, conns...)
}

// 链接建立后的额外操作
func (userManager *userManager) registerHook(conn Conn) {
	if chat.UserService.GetValidAdmin(conn.GetUserId()) == 0 && !chat.ManualService.IsIn(conn.GetUserId()) {
		rule := autoRuleRepo.GetEnter()
		if rule != nil {
			msg := rule.GetReplyMessage(conn.GetUserId())
			if msg != nil {
				messageRepo.Save(msg)
				rule.Count++
				autoRuleRepo.Save(rule)
				conn.Deliver(NewReceiveAction(msg))
			}
		}
	}
}



// 加入人工列表
func (userManager *userManager) addToManual(uid int64) *models.ChatSession {
	if !chat.ManualService.IsIn(uid) {
		onlineServerCount := len(AdminManager.Clients)
		if onlineServerCount == 0 { // 如果没有在线客服
			rule := autoRuleRepo.GetAdminAllOffLine()
			if rule != nil {
				switch rule.ReplyType {
				case models.ReplyTypeMessage:
					msg := rule.GetReplyMessage(uid)
					if msg != nil {
						messageRepo.Save(msg)
						rule.AddCount()
						conn, exist := userManager.GetConn(uid)
						if exist {
							conn.Deliver(NewReceiveAction(msg))
						}
						return nil
					}
				default:
				}
			}
		}
		_ = chat.ManualService.Add(uid)
		AdminManager.publishWaitingUser()
		session := chat.SessionService.Get(uid, 0)
		if session == nil {
			session = chat.SessionService.Create(uid, models.ChatSessionTypeNormal)
		}
		return session
	}
	return nil

}

// 触发事件
func (userManager *userManager) triggerMessageEvent(scene string, message *models.Message) {
	rules := autoRuleRepo.GetAllActiveNormal()
	for _, rule := range rules {
		if rule.IsMatch(message.Content) && rule.SceneInclude(scene) {
			switch rule.ReplyType {
			// 转接人工客服
			case models.ReplyTypeTransfer:
				session := userManager.addToManual(message.UserId)
				if session != nil {
					message.SessionId = session.Id
					messageRepo.Save(message)
				}
				AdminManager.BroadcastWaitingUser()
			// 回复消息
			case models.ReplyTypeMessage:
				msg := rule.GetReplyMessage(message.UserId)
				if msg != nil {
					msg.SessionId = message.SessionId
					messageRepo.Save(msg)
					conn, exist := userManager.GetConn(msg.UserId)
					if exist {
						conn.Deliver(NewReceiveAction(msg))
					}
				}
			//触发事件
			case models.ReplyTypeEvent:
				switch rule.Key {
				case "break":
					adminId := chat.UserService.GetValidAdmin(message.UserId)
					if adminId > 0 {
						_ = chat.AdminService.RemoveUser(adminId, message.UserId)
					}
					msg := rule.GetReplyMessage(message.UserId)
					if msg != nil {
						msg.SessionId = message.SessionId
						messageRepo.Save(msg)
						conn, exist := userManager.GetConn(msg.UserId)
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
