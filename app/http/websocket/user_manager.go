package websocket

import (
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"github.com/spf13/viper"
	"strconv"
	"time"
	"ws/app/chat"
	"ws/app/contract"
	"ws/app/log"
	"ws/app/models"
	"ws/app/repositories"
	rpcClient "ws/app/rpc/client"
	"ws/app/util"
	"ws/app/wechat"
)

type userManager struct {
	manager
}

var UserManager *userManager

const TypeUser = "user"

func SetupUser() {
	UserManager = &userManager{
		manager{
			shardCount:   10,
			ipAddr:       util.GetIPs()[0] + ":" + viper.GetString("Rpc.Port"),
			ConnMessages: make(chan *ConnMessage, 100),
			types:        TypeUser,
		},
	}
	UserManager.onRegister = UserManager.registerHook
	UserManager.onUnRegister = UserManager.unRegisterHook
	UserManager.Run()
}

func (userManager *userManager) Run() {
	userManager.manager.Run()
	go userManager.handleReceiveMessage()
	userManager.Do(func() {
		//go userManager.handleRemoteMessage()
	}, nil)
}

// DeliveryMessage 投递消息
// 查询user是否在本机上，是则直接投递
// 查询user当前server，如果存在则投递到该channel上
// 最后则说明user不在线，处理相关逻辑
// remote 是否从消息队列读取的消息
func (userManager *userManager) DeliveryMessage(msg *models.Message, isRemote bool) {
	userConn, exist := UserManager.GetConn(msg.GetUser())
	if exist {
		userConn.Deliver(NewReceiveAction(msg))
		return
	} else if !isRemote && userManager.isCluster() {
		server := userManager.getUserServer(msg.UserId)
		if server != "" {
			rpcClient.SendMessage(msg.Id, server)
			return
		}
	}
	userManager.handleOffline(msg)
}

// 从conn接受消息并处理
func (userManager *userManager) handleReceiveMessage() {
	for {
		payload := <-userManager.ConnMessages
		go userManager.handleMessage(payload)
	}
}

// 订阅远程消息
//func (userManager *userManager) handleRemoteMessage() {
//	sub := mq.Subscribe(userManager.GetSubscribeChannel())
//	for {
//		message := sub.ReceiveMessage()
//		go func() {
//			log.Log.WithField("a-type", "publish/subscribe").
//				WithField("b-type", "subscribe").
//				WithField("c-type", "user").
//				Infof("<channel:%s><types:%s><data:%s>",
//					userManager.GetSubscribeChannel(),
//					message.Get("types"),
//					message.Get("data"))
//			switch message.Get("types").String() {
//			case mq.TypeMessage:
//				mid := message.Get("data").Int()
//				msg := repositories.MessageRepo.FirstById(mid)
//				if msg != nil {
//					userManager.DeliveryMessage(msg, true)
//				}
//			case mq.TypeWaitingUserCount:
//				gid := message.Get("data").Int()
//				if gid > 0 {
//					userManager.broadcastWaitingCount(gid)
//				}
//			}
//		}()
//	}
//}

// 处理离线逻辑
func (userManager *userManager) handleOffline(msg *models.Message) {
	hadSubscribe := chat.SubScribeService.IsSet(msg.UserId)
	user := repositories.UserRepo.FirstById(msg.UserId)
	if hadSubscribe && user != nil && user.GetMpOpenId() != "" {
		err := wechat.GetMp().GetSubscribe().Send(&subscribe.Message{
			ToUser:           user.GetMpOpenId(),
			TemplateID:       viper.GetString("Wechat.SubscribeTemplateIdOne"),
			Page:             viper.GetString("Wechat.ChatPath"),
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
				msg.GroupId = conn.GetGroupId()
				msg.ReceivedAT = time.Now().Unix()
				msg.User = conn.GetUser().(*models.User)
				msg.AdminId = chat.UserService.GetValidAdmin(conn.GetUserId())
				// 发送回执
				repositories.MessageRepo.Save(msg)
				conn.Deliver(NewReceiptAction(msg))
				// 有对应的客服对象
				if msg.AdminId > 0 {
					// 更新会话有效期
					session := repositories.ChatSessionRepo.FirstActiveByUser(conn.GetUserId(), msg.AdminId)
					if session == nil {
						return
					}
					_ = chat.AdminService.UpdateUser(msg.AdminId, msg.UserId)
					msg.SessionId = session.Id
					repositories.MessageRepo.Save(msg)
					AdminManager.DeliveryMessage(msg, false)
				} else { // 没有客服对象
					if chat.TransferService.GetUserTransferId(conn.GetUserId()) == 0 {
						if chat.ManualService.IsIn(conn.GetUserId(), conn.GetGroupId()) {
							session := repositories.ChatSessionRepo.FirstActiveByUser(conn.GetUserId(), 0)
							if session != nil {
								msg.SessionId = session.Id
							}
							repositories.MessageRepo.Save(msg)
							AdminManager.BroadcastWaitingUser(conn.GetGroupId())
						} else {
							if chat.SettingService.GetIsAutoTransferManual(conn.GetGroupId()) { // 自动转人工
								session := UserManager.addToManual(conn.GetUser())
								if session != nil {
									msg.SessionId = session.Id
								}
								repositories.MessageRepo.Save(msg)
								AdminManager.BroadcastWaitingUser(conn.GetGroupId())
								userManager.PublishWaitingCount(conn.GetGroupId())
							} else {
								repositories.MessageRepo.Save(msg)
								userManager.triggerMessageEvent(models.SceneNotAccepted, msg)
							}
						}
					}
				}
			}
		}
	}
}

// PublishWaitingCount 广播等待人数
func (userManager *userManager) PublishWaitingCount(groupId int64) {
	log.Log.WithField("type", "WEBSOCKET").
		Infof("<user><publish><waiting-count><group-id:%d>", groupId)
	userManager.Do(func() {
		//userManager.publishToAllChannel(&mq.Payload{
		//	Types: mq.TypeWaitingUserCount,
		//	Data:  groupId,
		//})
	}, func() {
		userManager.broadcastWaitingCount(groupId)
	})
}

// NoticeWaitingCount 推送前面等待人数
func (userManager *userManager) NoticeWaitingCount(conn Conn) {
	log.Log.WithField("type", "WEBSOCKET").
		Infof("<user><notice><waiting-count><user-id:%d>", conn.GetGroupId())
	uid := conn.GetUserId()
	uTime := chat.ManualService.GetTime(uid, conn.GetGroupId())
	count := chat.ManualService.GetCountByTime(conn.GetGroupId(), "-inf",
		strconv.FormatFloat(uTime, 'f', 0, 64))
	conn.Deliver(NewWaitingUserCount(count - 1))
}

// 广播前面等待人数
func (userManager *userManager) broadcastWaitingCount(gid int64) {
	log.Log.WithField("type", "WEBSOCKET").
		Infof("<user><broadcast><waiting-count><group-id:%d>", gid)
	conns := userManager.GetAllConn(gid)
	for _, conn := range conns {
		userManager.NoticeWaitingCount(conn)
	}
}
func (userManager *userManager) unRegisterHook(conn Conn) {
	AdminManager.NoticeUserOffline(conn.GetUser())
}

// 链接建立后的额外操作
// 如果已经在待接入人工列表中，则推送当前队列位置
// 如果不在待接入人工列表中且没有设置客服，则推送欢迎语
func (userManager *userManager) registerHook(conn Conn) {
	AdminManager.NoticeUserOnline(conn.GetUser())
	if chat.ManualService.IsIn(conn.GetUserId(), conn.GetGroupId()) {
		userManager.NoticeWaitingCount(conn)
	} else if chat.UserService.GetValidAdmin(conn.GetUserId()) == 0 {
		rule := repositories.AutoRuleRepo.GetEnterByGroup(conn.GetGroupId())
		if rule != nil {
			msg := rule.GetReplyMessage(conn.GetUserId())
			if msg != nil {
				repositories.MessageRepo.Save(msg)
				rule.Count++
				repositories.AutoRuleRepo.Save(rule)
				conn.Deliver(NewReceiveAction(msg))
			}
		}
	}
}

// 加入人工列表
func (userManager *userManager) addToManual(user contract.User) *models.ChatSession {
	if !chat.ManualService.IsIn(user.GetPrimaryKey(), user.GetGroupId()) {
		onlineServerCount := AdminManager.GetOnlineTotal(user.GetGroupId())
		if onlineServerCount == 0 { // 如果没有在线客服
			rule := repositories.AutoRuleRepo.GetAdminAllOffLine(user.GetGroupId())
			if rule != nil {
				switch rule.ReplyType {
				case models.ReplyTypeMessage:
					msg := rule.GetReplyMessage(user.GetPrimaryKey())
					if msg != nil {
						repositories.MessageRepo.Save(msg)
						rule.AddCount()
						userManager.DeliveryMessage(msg, false)
						return nil
					}
				default:
				}
			}
		}
		_ = chat.ManualService.Add(user.GetPrimaryKey(), user.GetGroupId())
		session := repositories.ChatSessionRepo.FirstActiveByUser(user.GetPrimaryKey(), 0)
		if session == nil {
			session = repositories.ChatSessionRepo.Create(user.GetPrimaryKey(),
				user.GetGroupId(),
				models.ChatSessionTypeNormal)
		}
		message := repositories.MessageRepo.NewNotice(session, "正在为你转接人工客服")
		repositories.MessageRepo.Save(message)
		userManager.DeliveryMessage(message, false)
		return session
	}
	return nil

}

// 触发事件
func (userManager *userManager) triggerMessageEvent(scene string, message *models.Message) {
	rules := repositories.AutoRuleRepo.GetAllActiveNormalByGroup(message.GroupId)
	for _, rule := range rules {
		if rule.IsMatch(message.Content) && rule.SceneInclude(scene) {
			switch rule.ReplyType {
			// 转接人工客服
			case models.ReplyTypeTransfer:
				session := userManager.addToManual(message.GetUser())
				if session != nil {
					message.SessionId = session.Id
					repositories.MessageRepo.Save(message)
				}
				AdminManager.BroadcastWaitingUser(message.GroupId)
				userManager.PublishWaitingCount(message.GroupId)
				AdminManager.BroadcastWaitingUser(message.GetUser().GetGroupId())
			// 回复消息
			case models.ReplyTypeMessage:
				msg := rule.GetReplyMessage(message.UserId)
				if msg != nil {
					msg.SessionId = message.SessionId
					repositories.MessageRepo.Save(msg)
					userManager.DeliveryMessage(msg, false)
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
						repositories.MessageRepo.Save(msg)
						userManager.DeliveryMessage(msg, false)
					}
				}
			}
			rule.AddCount()
			return
		}
	}
}
