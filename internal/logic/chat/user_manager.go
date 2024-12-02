package chat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

func newUserManager() *userManager {
	userM = &userManager{
		&manager{
			shardCount:   10,
			connMessages: make(chan *chatConnMessage, 100),
			types:        "user",
		},
	}
	userM.onRegister = userM.registerHook
	userM.onUnRegister = userM.unRegisterHook
	return userM
}

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
func (s *userManager) DeliveryMessage(ctx context.Context, message *model.CustomerChatMessage) error {
	userConn, exist := s.GetConn(message.CustomerId, message.UserId)
	switch message.Type {
	case consts.MessageTypeRate:
		session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{Id: message.SessionId})
		if err != nil {
			return err
		}
		err = service.ChatSession().Close(ctx, session, false, true)
		if err != nil {
			return err
		}
	}
	if exist {
		userConn.Deliver(newReceiveAction(message))
		return nil
	}
	return s.handleOffline(ctx, message)
}

// NoticeQueueLocation 等待人数
func (s *userManager) NoticeQueueLocation(ctx context.Context, conn iWsConn) (err error) {
	uid := conn.GetUserId()
	uTime, err := manual.getAddTime(ctx, uid, conn.GetCustomerId())
	if err != nil {
		return
	}
	count, err := manual.getCountByTime(ctx, conn.GetCustomerId(), "-inf",
		strconv.FormatFloat(uTime, 'f', 0, 64))
	g.Dump(count)
	if err != nil {
		return
	}
	conn.Deliver(newWaitingUserCountAction(count - 1))
	return
}

func (s *userManager) BroadcastQueueLocation(ctx context.Context, customerId uint) error {
	return s.BroadcastLocalQueueLocation(ctx, customerId)
}

// BroadcastLocalQueueLocation 广播前面等待人数
func (s *userManager) BroadcastLocalQueueLocation(ctx context.Context, customerId uint) error {
	conns := s.GetAllConn(customerId)
	for _, conn := range conns {
		if manual.isInSet(ctx, conn.GetUserId(), conn.GetCustomerId()) {
			err := s.NoticeQueueLocation(ctx, conn)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// 从conn接受消息并处理
func (s *userManager) handleReceiveMessage() {
	for {
		payload := <-s.connMessages
		go func() {
			ctx := gctx.New()
			err := s.handleMessage(ctx, payload)
			if err != nil {
				g.Log().Error(ctx, err)
			}
		}()
	}
}

// 处理离线逻辑
func (s *userManager) handleOffline(ctx context.Context, msg *model.CustomerChatMessage) error {
	// todo
	return nil
}

// 处理消息
func (s *userManager) handleMessage(ctx context.Context, payload *chatConnMessage) (err error) {
	msg := payload.Msg
	conn := payload.Conn
	msg.Source = consts.MessageSourceUser
	msg.UserId = conn.GetUserId()
	msg.AdminId = service.ChatRelation().GetUserValidAdmin(ctx, msg.UserId)
	msg, err = service.ChatMessage().Insert(ctx, msg)
	if err != nil {
		return
	}
	// 发送回执
	conn.Deliver(newReceiptAction(msg))
	if msg.AdminId > 0 {
		session, err := service.ChatSession().FirstActive(ctx, msg.UserId, msg.AdminId, nil)
		if err != nil {
			return err
		}
		err = service.ChatRelation().UpdateUser(ctx, msg.AdminId, msg.UserId)
		if err != nil {
			return err
		}
		msg.SessionId = session.Id
		_, err = service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
			SessionId: msg.SessionId,
		})
		if err != nil {
			return err
		}
		err = s.triggerMessageEvent(ctx, consts.AutoRuleSceneAdminOnline, msg, conn.GetUser())
		if err != nil {
			return nil
		}
		return adminM.deliveryMessage(ctx, msg)
	} else {
		// 触发自动回复事件
		err = s.triggerMessageEvent(ctx, consts.AutoRuleSceneNotAccepted, msg, conn.GetUser())
		if err != nil {
			g.Log().Error(ctx, err)
		}
		var transferAdminId uint
		// 转接adminId
		transferAdminId, err = service.ChatTransfer().GetUserTransferId(ctx, conn.GetCustomerId(), conn.GetUserId())
		if err != nil {
			return
		}
		var session *model.CustomerChatSession
		if transferAdminId == 0 {
			// 在代人工接入列表中
			inManual := manual.isInSet(ctx, conn.GetUserId(), conn.GetCustomerId())

			if inManual {
				session, err = service.ChatSession().FirstNormal(ctx, msg.UserId, 0)
				if err != nil {
					return
				}
				msg.SessionId = session.Id
				_, err = service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
					SessionId: msg.SessionId,
				})
				if err != nil {
					return
				}
				err = adminM.broadcastWaitingUser(ctx, conn.GetCustomerId())
				if err != nil {
					return
				}
			} else {
				// 不在代人工接入列表中
				var isAutoAdd bool
				isAutoAdd, err = service.ChatSetting().GetIsAutoTransferManual(ctx, conn.GetCustomerId())
				if err != nil {
					return err
				}
				if isAutoAdd { // 如果自动转人工
					session, err = s.addToManual(ctx, conn.GetUser())
					if err != nil {
						return
					}
					if session != nil {
						msg.SessionId = session.Id
						_, err = service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
							SessionId: msg.SessionId,
						})
					}
					if err != nil {
						return
					}
					err = adminM.broadcastWaitingUser(ctx, conn.GetCustomerId())
					if err != nil {
						return
					}
					err = s.BroadcastQueueLocation(ctx, conn.GetCustomerId())
					if err != nil {
						return
					}
				}
			}
		} else {
			session, err = service.ChatSession().FirstTransfer(ctx, conn.GetUserId(), transferAdminId)
			if err != nil {
				return
			}
			msg.SessionId = session.Id
			_, err = service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
				SessionId: msg.SessionId,
			})
			if err != nil {
				return
			}
		}
	}
	return nil
}

// 触发进入事件，只有没有对应客服的情况下触发，10分钟内多触发一次
func (s *userManager) triggerEnterEvent(ctx context.Context, conn iWsConn) (err error) {
	if service.ChatRelation().GetUserValidAdmin(ctx, conn.GetUserId()) > 0 {
		return
	}
	cacheKey := fmt.Sprintf("welcome:%d", conn.GetUserId())
	val, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		return
	}
	if !val.IsNil() {
		return
	}
	var rule *model.CustomerChatAutoRule
	rule, err = service.AutoRule().GetEnterRule(ctx, conn.GetCustomerId())
	if err != nil {
		return
	}
	var autoMsg *model.CustomerChatAutoMessage
	autoMsg, err = service.AutoRule().GetMessage(ctx, rule)
	if err != nil {
		return
	}
	var entityMessage *model.CustomerChatMessage
	entityMessage, err = service.AutoMessage().ToChatMessage(autoMsg)
	if err != nil {
		return
	}
	entityMessage.UserId = conn.GetUserId()
	_, err = service.ChatMessage().Save(ctx, entityMessage)
	if err != nil {
		return err
	}
	rule.Count++
	err = service.AutoRule().IncrTriggerCount(ctx, rule)
	if err != nil {
		return
	}
	conn.Deliver(newReceiveAction(entityMessage))
	err = gcache.Set(ctx, cacheKey, gtime.Now().String(), time.Minute*10)
	if err != nil {
		return
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
	ctx := gctx.New()
	adminM.noticeUserOnline(ctx, conn)
	var err error
	if manual.isInSet(ctx, conn.GetUserId(), conn.GetCustomerId()) {
		err = s.NoticeQueueLocation(ctx, conn)
	} else {
		err = s.triggerEnterEvent(ctx, conn)
	}
	if err != nil {
		g.Log().Error(ctx, err)
	}
}

// 加入人工列表
func (s *userManager) addToManual(ctx context.Context, user IChatUser) (session *model.CustomerChatSession, err error) {
	if manual.isInSet(ctx, user.GetPrimaryKey(), user.GetCustomerId()) {
		err = gerror.New("is in")
		return
	}
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
				err = service.AutoRule().IncrTriggerCount(ctx, rule)
				if err != nil {
					return nil, err
				}
				message.UserId = user.GetPrimaryKey()
				message, err = service.ChatMessage().Insert(ctx, message)
				if err != nil {
					return nil, err
				}
				err = s.DeliveryMessage(ctx, message)
				if err != nil {
					return nil, err
				}
				return nil, nil
			}
		}
	}
	err = manual.addToSet(ctx, user.GetPrimaryKey(), user.GetCustomerId())
	if err != nil {
		return
	}
	session, err = service.ChatSession().FirstNormal(ctx, user.GetPrimaryKey(), 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			session, err = service.ChatSession().Create(ctx, user.GetPrimaryKey(), user.GetCustomerId(), consts.ChatSessionTypeNormal)
			if err != nil {
				return nil, err
			}
		} else {
			return
		}
	}
	message := service.ChatMessage().NewNotice(session, "正在为你转接人工客服")
	_, err = service.ChatMessage().Save(ctx, message)
	if err != nil {
		return nil, err
	}
	err = s.DeliveryMessage(ctx, message)
	if err != nil {
		return nil, err
	}
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
	return
}

// 触发事件
func (s *userManager) triggerMessageEvent(ctx context.Context, scene string, message *model.CustomerChatMessage, user IChatUser) (err error) {
	rules, err := service.AutoRule().AllActive(ctx, user.GetCustomerId())
	if err != nil {
		return
	}
	for _, rule := range rules {
		isMatch := service.AutoRule().IsMatch(rule, scene, message.Content)
		if !isMatch {
			continue
		}
		switch rule.ReplyType {
		// 转接人工客服
		case consts.AutoRuleReplyTypeTransfer:
			var isTransfer bool
			isTransfer, err = service.ChatTransfer().IsInTransfer(ctx, user.GetCustomerId(), user.GetPrimaryKey())
			if err != nil {
				return
			}
			if isTransfer {
				return nil
			}
			var session *model.CustomerChatSession
			session, err = s.addToManual(ctx, user)
			if err != nil {
				return
			}
			message.SessionId = session.Id
			_, err = service.ChatMessage().Save(ctx, message)
			if err != nil {
				return
			}
			err = adminM.broadcastWaitingUser(ctx, message.CustomerId)
			if err != nil {
				return
			}
			err = s.BroadcastQueueLocation(ctx, message.CustomerId)
			if err != nil {
				return
			}
			err = adminM.broadcastWaitingUser(ctx, user.GetCustomerId())
			if err != nil {
				return
			}
			err = service.AutoRule().IncrTriggerCount(ctx, rule)
			if err != nil {
				return
			}
			return nil

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
			_ = service.AutoRule().IncrTriggerCount(ctx, rule)
			return nil
		}
	}
	return nil
}
