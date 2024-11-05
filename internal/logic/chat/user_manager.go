package chat

import (
	"fmt"
	"gf-chat/internal/consts"
	"gf-chat/internal/contract"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"time"
)

type userManager struct {
	*manager
}

func (s *userManager) run() {
	s.Run()
	go s.handleReceiveMessage()
}

// DeliveryMessage 投递消息
// 查询user是否在本机上，是则直接投递
// 查询user当前server，如果存在则投递到该channel上
// 最后则说明user不在线，处理相关逻辑
func (s *userManager) DeliveryMessage(message *relation.CustomerChatMessages) {
	userConn, exist := s.GetConn(message.CustomerId, message.UserId)
	switch message.Type {
	case consts.MessageTypeRate:
		session := service.ChatSession().First(gctx.New(), do.CustomerChatSessions{Id: message.SessionId})
		if session != nil {
			service.ChatSession().Close(session, false, true)
		}
	}
	if exist {
		userConn.Deliver(service.Action().NewReceiveAction(message))
		return
	}
	s.handleOffline(message)
}

// NoticeQueueLocation 等待人数
func (s *userManager) NoticeQueueLocation(conn iWsConn) {
	uid := conn.GetUserId()
	uTime := service.ChatManual().GetTime(uid, conn.GetCustomerId())
	count := service.ChatManual().GetCountByTime(conn.GetCustomerId(), "-inf",
		strconv.FormatFloat(uTime, 'f', 0, 64))
	conn.Deliver(service.Action().NewWaitingUserCount(count - 1))
}

func (s *userManager) BroadcastQueueLocation(customerId int) {
	s.BroadcastLocalQueueLocation(customerId)
}

// BroadcastLocalQueueLocation 广播前面等待人数
func (s *userManager) BroadcastLocalQueueLocation(customerId int) {
	conns := s.GetAllConn(customerId)
	for _, conn := range conns {
		s.NoticeQueueLocation(conn)
	}
}

// 从conn接受消息并处理
func (s *userManager) handleReceiveMessage() {
	for {
		payload := <-s.ConnMessages
		go s.handleMessage(payload)
	}
}

// 处理离线逻辑
func (s *userManager) handleOffline(msg *relation.CustomerChatMessages) {
	_ = service.SubscribeMsg().Send(gctx.New(), msg.CustomerId, msg.UserId)
}

// 处理消息
func (s *userManager) handleMessage(payload *chatConnMessage) {
	msg := payload.Msg
	conn := payload.Conn
	msg.Source = consts.MessageSourceUser
	msg.UserId = conn.GetUserId()
	msg.AdminId = service.ChatRelation().GetUserValidAdmin(msg.UserId)
	service.ChatMessage().SaveRelationOne(msg)
	// 发送回执
	conn.Deliver(service.Action().NewReceiptAction(msg))
	// 有对应的客服对象
	if msg.AdminId > 0 {
		// 更新会话有效期
		session := service.ChatSession().ActiveOne(msg.UserId, msg.AdminId, nil)
		if session == nil {
			return
		}
		_ = service.ChatRelation().UpdateUser(msg.AdminId, msg.UserId)
		msg.SessionId = session.Id
		service.ChatMessage().SaveRelationOne(msg)
		s.triggerMessageEvent(consts.AutoRuleSceneAdminOnline, msg, conn.GetUser())
		adminM.deliveryMessage(msg)
	} else {
		s.triggerMessageEvent(consts.AutoRuleSceneNotAccepted, msg, conn.GetUser())
		// 转接adminId
		transferAdminId := service.ChatTransfer().GetUserTransferId(conn.GetCustomerId(), conn.GetUserId())
		if transferAdminId == 0 {
			// 在代人工接入列表中
			inManual := service.ChatManual().IsIn(conn.GetUserId(), conn.GetCustomerId())
			if inManual {
				session := service.ChatSession().ActiveNormalOne(msg.UserId, 0)
				if session != nil {
					msg.SessionId = session.Id
				}
				service.ChatMessage().SaveRelationOne(msg)
				adminM.broadcastWaitingUser(conn.GetCustomerId())
			} else {
				// 不在代人工接入列表中
				if service.ChatSetting().GetIsAutoTransferManual(conn.GetCustomerId()) { // 如果自动转人工
					session := s.addToManual(conn.GetUser())
					if session != nil {
						msg.SessionId = session.Id
					}
					service.ChatMessage().SaveRelationOne(msg)
					adminM.broadcastWaitingUser(conn.GetCustomerId())
					s.BroadcastQueueLocation(conn.GetCustomerId())
				}
			}
		} else {
			session := service.ChatSession().ActiveTransferOne(conn.GetUserId(), transferAdminId)
			if session != nil {
				msg.SessionId = session.Id
				service.ChatMessage().SaveRelationOne(msg)
			}
		}
	}
}

// 触发进入事件，只有没有对应客服的情况下触发，10分钟内多触发一次
func (s *userManager) triggerEnterEvent(conn iWsConn) {
	if service.ChatRelation().GetUserValidAdmin(conn.GetUserId()) > 0 {
		return
	}
	ctx := gctx.New()
	cacheKey := fmt.Sprintf("welcome:%d", conn.GetUserId())
	val, _ := gcache.Get(ctx, cacheKey)
	if val.String() == "" {
		rule := service.AutoRule().GetEnterRule(conn.GetCustomerId())
		if rule == nil {
			return
		}
		autoMsg := service.AutoRule().GetMessage(rule)
		if autoMsg == nil {
			return
		}
		entityMessage := service.AutoMessage().ToChatMessage(autoMsg)
		entityMessage.UserId = conn.GetUserId()
		relationMessage := service.ChatMessage().EntityToRelation(entityMessage)
		service.ChatMessage().SaveRelationOne(relationMessage)
		rule.Count++
		_ = service.AutoRule().Increment(rule)
		conn.Deliver(service.Action().NewReceiveAction(relationMessage))
		_ = gcache.Set(ctx, cacheKey, 1, time.Minute*10)
	}
}

func (s *userManager) unRegisterHook(conn iWsConn) {
	adminM.noticeUserOffline(conn.GetUser())
}

// 链接建立后的额外操作
// 如果已经在待接入人工列表中，则推送当前队列位置
// 如果不在待接入人工列表中且没有设置客服，则触发进入事件
func (s *userManager) registerHook(conn iWsConn) {
	adminM.noticeUserOnline(conn)
	if service.ChatManual().IsIn(conn.GetUserId(), conn.GetCustomerId()) {
		s.NoticeQueueLocation(conn)
	} else {
		s.triggerEnterEvent(conn)
	}
}

// 加入人工列表
func (s *userManager) addToManual(user contract.IChatUser) *entity.CustomerChatSessions {
	if !service.ChatManual().IsIn(user.GetPrimaryKey(), user.GetCustomerId()) {
		onlineServerCount := adminM.GetOnlineTotal(user.GetCustomerId())
		if onlineServerCount == 0 { // 如果没有在线客服
			rule := service.AutoRule().GetSystemOne(user.GetCustomerId(), consts.AutoRuleMatchAdminAllOffLine)
			if rule != nil {
				switch rule.ReplyType {
				case consts.AutoRuleReplyTypeMessage:
					autoMessage := service.AutoRule().GetMessage(rule)
					if autoMessage != nil {
						message := service.AutoMessage().ToChatMessage(autoMessage)
						if message != nil {
							service.ChatMessage().SaveOne(message)
							service.AutoRule().Increment(rule)
							message.UserId = user.GetPrimaryKey()
							relationMessage := service.ChatMessage().EntityToRelation(message)
							service.ChatMessage().SaveRelationOne(relationMessage)
							s.DeliveryMessage(relationMessage)
							return nil
						}
					}
				default:
				}
			}
		}
		_ = service.ChatManual().Add(user.GetPrimaryKey(), user.GetCustomerId())
		session := service.ChatSession().ActiveNormalOne(user.GetPrimaryKey(), 0)
		if session == nil {
			session = service.ChatSession().
				Create(user.GetPrimaryKey(), user.GetCustomerId(), consts.ChatSessionTypeNormal)
		}
		message := service.ChatMessage().NewNotice(session, "正在为你转接人工客服")
		service.ChatMessage().SaveOne(message)
		relationMessage := service.ChatMessage().EntityToRelation(message)
		s.DeliveryMessage(relationMessage)
		// 没有客服在线则发送公众号消息
		go func() {
			if onlineServerCount == 0 {
				admins := service.Admin().GetChatAll(user.GetCustomerId())
				for _, admin := range admins {
					adminM.sendWaiting(admin, user)
				}
			}
		}()
		return session
	}
	return nil

}

// 触发事件
func (s *userManager) triggerMessageEvent(scene string, message *relation.CustomerChatMessages, user contract.IChatUser) {
	rules := service.AutoRule().GetActiveByCustomer(message.CustomerId)
	for _, rule := range rules {
		isMatch := service.AutoRule().IsMatch(rule, scene, message.Content)
		if isMatch {
			switch rule.ReplyType {
			// 转接人工客服
			case consts.AutoRuleReplyTypeTransfer:
				transferId := service.ChatTransfer().GetUserTransferId(user.GetCustomerId(), user.GetPrimaryKey())
				if transferId == 0 {
					session := s.addToManual(user)
					if session != nil {
						message.SessionId = session.Id
						service.ChatMessage().SaveRelationOne(message)
					}
					adminM.broadcastWaitingUser(message.CustomerId)
					s.BroadcastQueueLocation(message.CustomerId)
					adminM.broadcastWaitingUser(user.GetCustomerId())
					service.AutoRule().Increment(rule)
					return
				}

			// 回复消息
			case consts.AutoRuleReplyTypeMessage:
				autoMessage := service.AutoRule().GetMessage(rule)
				if autoMessage != nil {
					msg := service.AutoMessage().ToChatMessage(autoMessage)
					msg.UserId = message.UserId
					msg.SessionId = message.SessionId
					relationMsg := service.ChatMessage().EntityToRelation(msg)
					service.ChatMessage().SaveRelationOne(relationMsg)
					conn, exist := s.GetConn(user.GetCustomerId(), user.GetPrimaryKey())
					if exist {
						s.SendAction(service.Action().NewReceiveAction(relationMsg), conn)
					}
					service.AutoRule().Increment(rule)
					return
				}
			}
		}
	}
}
