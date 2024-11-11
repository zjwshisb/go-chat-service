package chat

import (
	"fmt"
	"gf-chat/internal/consts"
	"gf-chat/internal/contract"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
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
func (s *userManager) DeliveryMessage(message *model.CustomerChatMessage) {
	userConn, exist := s.GetConn(message.CustomerId, message.UserId)
	ctx := gctx.New()
	switch message.Type {
	case consts.MessageTypeRate:
		session, _ := service.ChatSession().First(ctx, do.CustomerChatSessions{Id: message.SessionId})
		if session != nil {
			service.ChatSession().Close(ctx, session, false, true)
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

func (s *userManager) BroadcastQueueLocation(customerId uint) {
	s.BroadcastLocalQueueLocation(customerId)
}

// BroadcastLocalQueueLocation 广播前面等待人数
func (s *userManager) BroadcastLocalQueueLocation(customerId uint) {
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
func (s *userManager) handleOffline(msg *model.CustomerChatMessage) {
	_ = service.SubscribeMsg().Send(gctx.New(), msg.CustomerId, msg.UserId)
}

// 处理消息
func (s *userManager) handleMessage(payload *chatConnMessage) {
	ctx := gctx.New()
	msg := payload.Msg
	conn := payload.Conn
	msg.Source = consts.MessageSourceUser
	msg.UserId = conn.GetUserId()
	msg.AdminId = service.ChatRelation().GetUserValidAdmin(ctx, msg.UserId)
	service.ChatMessage().Save(ctx, msg)
	// 发送回执
	conn.Deliver(service.Action().NewReceiptAction(msg))
	// 有对应的客服对象
	if msg.AdminId > 0 {
		// 更新会话有效期
		session, _ := service.ChatSession().ActiveOne(ctx, msg.UserId, msg.AdminId, nil)
		if session == nil {
			return
		}
		_ = service.ChatRelation().UpdateUser(ctx, msg.AdminId, msg.UserId)
		msg.SessionId = session.Id
		service.ChatMessage().Save(ctx, msg)
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
				session, _ := service.ChatSession().ActiveNormalOne(ctx, msg.UserId, 0)
				if session != nil {
					msg.SessionId = session.Id
				}
				service.ChatMessage().Save(ctx, msg)
				adminM.broadcastWaitingUser(conn.GetCustomerId())
			} else {
				// 不在代人工接入列表中
				if service.ChatSetting().GetIsAutoTransferManual(conn.GetCustomerId()) { // 如果自动转人工
					session, _ := s.addToManual(conn.GetUser())
					if session != nil {
						msg.SessionId = session.Id
					}
					service.ChatMessage().Save(ctx, msg)
					adminM.broadcastWaitingUser(conn.GetCustomerId())
					s.BroadcastQueueLocation(conn.GetCustomerId())
				}
			}
		} else {
			session, _ := service.ChatSession().ActiveTransferOne(ctx, conn.GetUserId(), transferAdminId)
			if session != nil {
				msg.SessionId = session.Id
				service.ChatMessage().Save(ctx, msg)
			}
		}
	}
}

// 触发进入事件，只有没有对应客服的情况下触发，10分钟内多触发一次
func (s *userManager) triggerEnterEvent(conn iWsConn) {
	ctx := gctx.New()
	if service.ChatRelation().GetUserValidAdmin(ctx, conn.GetUserId()) > 0 {
		return
	}

	cacheKey := fmt.Sprintf("welcome:%d", conn.GetUserId())
	val, _ := gcache.Get(ctx, cacheKey)
	if val.String() == "" {
		rule, _ := service.AutoRule().GetEnterRule(conn.GetCustomerId())
		if rule == nil {
			return
		}
		autoMsg := service.AutoRule().GetMessage(rule)
		if autoMsg == nil {
			return
		}
		entityMessage, err := service.AutoMessage().ToChatMessage(autoMsg)
		if err != nil {
			return
		}
		entityMessage.UserId = conn.GetUserId()
		service.ChatMessage().Save(ctx, entityMessage)
		rule.Count++
		_ = service.AutoRule().Increment(rule)
		conn.Deliver(service.Action().NewReceiveAction(entityMessage))
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
func (s *userManager) addToManual(user contract.IChatUser) (session *model.CustomerChatSession, err error) {
	ctx := gctx.New()
	if !service.ChatManual().IsIn(user.GetPrimaryKey(), user.GetCustomerId()) {
		onlineServerCount := adminM.GetOnlineTotal(user.GetCustomerId())
		if onlineServerCount == 0 { // 如果没有在线客服
			rule, _ := service.AutoRule().GetSystemOne(user.GetCustomerId(), consts.AutoRuleMatchAdminAllOffLine)
			if rule != nil {
				switch rule.ReplyType {
				case consts.AutoRuleReplyTypeMessage:
					autoMessage := service.AutoRule().GetMessage(rule)
					if autoMessage == nil {
						return nil, gerror.New("no message")
					}
					message, err := service.AutoMessage().ToChatMessage(autoMessage)
					if err != nil {
						return nil, err
					}
					service.ChatMessage().Save(ctx, message)
					service.AutoRule().Increment(rule)
					message.UserId = user.GetPrimaryKey()
					service.ChatMessage().Save(ctx, message)
					s.DeliveryMessage(message)
					return nil, nil
				default:
				}
			}
		}
		_ = service.ChatManual().Add(user.GetPrimaryKey(), user.GetCustomerId())
		session, _ := service.ChatSession().ActiveNormalOne(ctx, user.GetPrimaryKey(), 0)

		if session == nil {
			session, err = service.ChatSession().Create(ctx, user.GetPrimaryKey(), user.GetCustomerId(), consts.ChatSessionTypeNormal)
			if err != nil {
				return nil, err
			}
		}
		message := service.ChatMessage().NewNotice(session, "正在为你转接人工客服")
		service.ChatMessage().Save(ctx, message)
		s.DeliveryMessage(message)
		// 没有客服在线则发送公众号消息
		go func() {
			if onlineServerCount == 0 {
				admins, _ := service.Admin().All(ctx, do.CustomerAdmins{
					CustomerId: user.GetCustomerId(),
				}, nil, nil)
				for _, admin := range admins {
					adminM.sendWaiting(admin, user)
				}
			}
		}()
		return session, nil
	}
	return nil, gerror.New("is in")

}

// 触发事件
func (s *userManager) triggerMessageEvent(scene string, message *model.CustomerChatMessage, user contract.IChatUser) {
	rules := service.AutoRule().GetActiveByCustomer(message.CustomerId)
	ctx := gctx.New()
	for _, rule := range rules {
		isMatch := service.AutoRule().IsMatch(rule, scene, message.Content)
		if isMatch {
			switch rule.ReplyType {
			// 转接人工客服
			case consts.AutoRuleReplyTypeTransfer:
				transferId := service.ChatTransfer().GetUserTransferId(user.GetCustomerId(), user.GetPrimaryKey())
				if transferId == 0 {
					session, _ := s.addToManual(user)

					if session != nil {
						message.SessionId = session.Id
						service.ChatMessage().Save(ctx, message)
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
					msg, err := service.AutoMessage().ToChatMessage(autoMessage)
					if err == nil {
						msg.UserId = message.UserId
						msg.SessionId = message.SessionId
						service.ChatMessage().Save(ctx, msg)
						conn, exist := s.GetConn(user.GetCustomerId(), user.GetPrimaryKey())
						if exist {
							s.SendAction(service.Action().NewReceiveAction(msg), conn)
						}
						service.AutoRule().Increment(rule)
					}

					return
				}
			}
		}
	}
}
