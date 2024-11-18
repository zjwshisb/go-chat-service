package chat

import (
	"fmt"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"
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
			err := service.ChatSession().Close(ctx, session, false, true)
			if err != nil {
				g.Log().Error(ctx, err)
			}
		}
	}
	if exist {
		userConn.Deliver(newReceiveAction(message))
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
	conn.Deliver(newWaitingUserCountAction(count - 1))
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
		go func() {
			err := s.handleMessage(payload)
			if err != nil {
				g.Log().Error(gctx.New(), err)
			}
		}()
	}
}

// 处理离线逻辑
func (s *userManager) handleOffline(msg *model.CustomerChatMessage) {
	// todo
}

// 处理消息
func (s *userManager) handleMessage(payload *chatConnMessage) (err error) {
	ctx := gctx.New()
	msg := payload.Msg
	conn := payload.Conn
	msg.Source = consts.MessageSourceUser
	msg.UserId = conn.GetUserId()
	msg.AdminId = service.ChatRelation().GetUserValidAdmin(ctx, msg.UserId)
	_, err = service.ChatMessage().Save(ctx, msg)
	if err != nil {
		return
	}
	// 发送回执
	conn.Deliver(newReceiptAction(msg))
	// 有对应的客服对象
	if msg.AdminId > 0 {
		// 更新会话有效期
		session, err := service.ChatSession().FirstActive(ctx, msg.UserId, msg.AdminId, nil)
		if err != nil {
			return err
		}
		err = service.ChatRelation().UpdateUser(ctx, msg.AdminId, msg.UserId)
		if err != nil {
			return err
		}
		msg.SessionId = session.Id
		_, err = service.ChatMessage().Save(ctx, msg)
		if err != nil {
			return err
		}
		err = s.triggerMessageEvent(consts.AutoRuleSceneAdminOnline, msg, conn.GetUser())
		if err != nil {
			return nil
		}
		adminM.deliveryMessage(msg)
	} else {
		err = s.triggerMessageEvent(consts.AutoRuleSceneNotAccepted, msg, conn.GetUser())
		if err != nil {
			return nil
		}
		// 转接adminId
		transferAdminId, err := service.ChatTransfer().GetUserTransferId(ctx, conn.GetCustomerId(), conn.GetUserId())
		if err != nil {
			return err
		}
		if transferAdminId == 0 {
			// 在代人工接入列表中
			inManual := service.ChatManual().IsIn(conn.GetUserId(), conn.GetCustomerId())
			if inManual {
				session, err := service.ChatSession().FirstNormal(ctx, msg.UserId, 0)
				if err != nil {
					return err
				}
				msg.SessionId = session.Id
				_, err = service.ChatMessage().Save(ctx, msg)
				if err != nil {
					return err
				}
				adminM.broadcastWaitingUser(conn.GetCustomerId())
			} else {
				// 不在代人工接入列表中
				isAutoAdd, err := service.ChatSetting().GetIsAutoTransferManual(ctx, conn.GetCustomerId())
				if err != nil {
					return err
				}
				if isAutoAdd { // 如果自动转人工
					session, err := s.addToManual(conn.GetUser())
					if err != nil {
						return err
					}
					if session != nil {
						msg.SessionId = session.Id
					}
					_, err = service.ChatMessage().Save(ctx, msg)
					if err != nil {
						return err
					}
					adminM.broadcastWaitingUser(conn.GetCustomerId())
					s.BroadcastQueueLocation(conn.GetCustomerId())
				}
			}
		} else {
			session, err := service.ChatSession().FirstTransfer(ctx, conn.GetUserId(), transferAdminId)
			if err != nil {
				return err
			}
			msg.SessionId = session.Id
			_, err = service.ChatMessage().Save(ctx, msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 触发进入事件，只有没有对应客服的情况下触发，10分钟内多触发一次
func (s *userManager) triggerEnterEvent(conn iWsConn) (err error) {
	ctx := gctx.New()
	if service.ChatRelation().GetUserValidAdmin(ctx, conn.GetUserId()) > 0 {
		return
	}

	cacheKey := fmt.Sprintf("welcome:%d", conn.GetUserId())
	val, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		return
	}
	if val.String() == "" {
		rule, err := service.AutoRule().GetEnterRule(ctx, conn.GetCustomerId())
		if err != nil {
			return err
		}
		autoMsg, err := service.AutoRule().GetMessage(ctx, rule)
		if err != nil {
			return err
		}
		var entityMessage *model.CustomerChatMessage
		entityMessage, err = service.AutoMessage().ToChatMessage(autoMsg)
		if err != nil {
			return err
		}
		entityMessage.UserId = conn.GetUserId()
		_, err = service.ChatMessage().Save(ctx, entityMessage)
		if err != nil {
			return err
		}
		rule.Count++
		_ = service.AutoRule().Increment(ctx, rule)
		conn.Deliver(newReceiveAction(entityMessage))
		_ = gcache.Set(ctx, cacheKey, 1, time.Minute*10)
	}
	return
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
		err := s.triggerEnterEvent(conn)
		if err != nil {
			g.Log().Error(gctx.New(), err)
		}
	}
}

// 加入人工列表
func (s *userManager) addToManual(user IChatUser) (session *model.CustomerChatSession, err error) {
	ctx := gctx.New()
	if !service.ChatManual().IsIn(user.GetPrimaryKey(), user.GetCustomerId()) {
		onlineServerCount := adminM.GetOnlineTotal(user.GetCustomerId())
		if onlineServerCount == 0 { // 如果没有在线客服
			rule, _ := service.AutoRule().GetSystemOne(ctx, user.GetCustomerId(), consts.AutoRuleMatchAdminAllOffLine)
			if rule != nil {
				switch rule.ReplyType {
				case consts.AutoRuleReplyTypeMessage:
					autoMessage, err := service.AutoRule().GetMessage(ctx, rule)
					if err != nil {
						return nil, err
					}
					message, err := service.AutoMessage().ToChatMessage(autoMessage)
					if err != nil {
						return nil, err
					}
					err = service.AutoRule().Increment(ctx, rule)
					if err != nil {
						return nil, err
					}
					message.UserId = user.GetPrimaryKey()
					_, err = service.ChatMessage().Save(ctx, message)
					if err != nil {
						return nil, err
					}
					s.DeliveryMessage(message)
					return nil, nil
				default:
				}
			}
		}
		_ = service.ChatManual().Add(user.GetPrimaryKey(), user.GetCustomerId())
		session, _ = service.ChatSession().FirstNormal(ctx, user.GetPrimaryKey(), 0)
		if session == nil {
			session, err = service.ChatSession().Create(ctx, user.GetPrimaryKey(), user.GetCustomerId(), consts.ChatSessionTypeNormal)
			if err != nil {
				return nil, err
			}
		}
		message := service.ChatMessage().NewNotice(session, "正在为你转接人工客服")
		_, err := service.ChatMessage().Save(ctx, message)
		if err != nil {
			return nil, err
		}
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
func (s *userManager) triggerMessageEvent(scene string, message *model.CustomerChatMessage, user IChatUser) (err error) {
	ctx := gctx.New()
	rules, err := service.AutoRule().AllActive(ctx, user.GetCustomerId())
	if err != nil {
		return
	}
	for _, rule := range rules {
		isMatch := service.AutoRule().IsMatch(rule, scene, message.Content)
		if isMatch {
			switch rule.ReplyType {
			// 转接人工客服
			case consts.AutoRuleReplyTypeTransfer:
				isTransfer, err := service.ChatTransfer().IsInTransfer(ctx, user.GetCustomerId(), user.GetPrimaryKey())
				if err != nil {
					return err
				}
				if !isTransfer {
					session, err := s.addToManual(user)
					if err != nil {
						return err
					}
					message.SessionId = session.Id
					_, err = service.ChatMessage().Save(ctx, message)
					if err != nil {
						return err
					}
					adminM.broadcastWaitingUser(message.CustomerId)
					s.BroadcastQueueLocation(message.CustomerId)
					adminM.broadcastWaitingUser(user.GetCustomerId())
					err = service.AutoRule().Increment(ctx, rule)
					if err != nil {
						return err
					}
					return nil
				}

			// 回复消息
			case consts.AutoRuleReplyTypeMessage:
				autoMessage, err := service.AutoRule().GetMessage(ctx, rule)
				if err != nil {
					return err
				}
				msg, err := service.AutoMessage().ToChatMessage(autoMessage)
				if err != nil {
					return err
				}
				msg.UserId = message.UserId
				msg.SessionId = message.SessionId
				_, err = service.ChatMessage().Save(ctx, msg)
				if err != nil {
					return err
				}
				conn, exist := s.GetConn(user.GetCustomerId(), user.GetPrimaryKey())
				if exist {
					s.SendAction(newReceiveAction(msg), conn)
				}
				_ = service.AutoRule().Increment(ctx, rule)
				return nil
			}
		}
	}
	return nil
}
