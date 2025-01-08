package chat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	grpc "gf-chat/api/chat/v1"
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
)

func newUserManager(cluster bool) *userManager {
	userM = &userManager{
		newManager(10, 100, time.Minute, cluster, consts.WsTypeUser),
	}
	userM.on(eventRegister, userM.onRegister)
	userM.on(eventUnRegister, userM.unRegisterHook)
	userM.on(eventMessage, userM.handleMessage)
	return userM
}

type userManager struct {
	*manager
}

func (s *userManager) deliveryMessage(ctx context.Context, message *model.CustomerChatMessage) error {
	if !s.cluster {
		return s.deliveryLocalMessage(ctx, message)
	}
	server, err := s.getUserServer(ctx, message.UserId)
	if err != nil {
		return err
	}
	if server != "" {
		_, err = service.Grpc().Client(server).SendUserMessage(ctx, &grpc.SendUserMessageRequest{MsgId: uint32(message.Id)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *userManager) deliveryLocalMessage(ctx context.Context, message *model.CustomerChatMessage) error {
	userConn, exist := s.getConn(message.CustomerId, message.UserId)
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
		userConn.deliver(action.newReceive(message))
	}
	return nil
}

// noticeQueueLocation 等待人数
func (s *userManager) noticeQueueLocation(ctx context.Context, conn iWsConn) (err error) {
	uTime, err := manual.getAddTime(ctx, conn.getUserId(), conn.getCustomerId())
	if err != nil {
		return
	}
	count, err := manual.getCountByTime(ctx, conn.getCustomerId(), "-inf",
		strconv.FormatFloat(uTime, 'f', 0, 64))
	if err != nil {
		return
	}
	conn.deliver(action.newWaitingUserCount(count - 1))
	return
}

func (s *userManager) broadcastQueueLocation(ctx context.Context, customerId uint) error {
	isSHowQueue, err := service.ChatSetting().GetIsUserShowQueue(ctx, customerId)
	if err != nil {
		return err
	}
	if isSHowQueue {
		return s.broadcastLocalQueueLocation(ctx, customerId)
	}
	return nil
}

// broadcastLocalQueueLocation 广播前面等待人数
func (s *userManager) broadcastLocalQueueLocation(ctx context.Context, customerId uint) error {
	conns := s.getAllConn(customerId)
	for _, conn := range conns {
		if manual.isInSet(ctx, conn.getUserId(), conn.getCustomerId()) {
			err := s.noticeQueueLocation(ctx, conn)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// 处理消息
func (s *userManager) handleMessage(ctx context.Context, arg eventArg) (err error) {
	msg := arg.msg
	conn := arg.conn
	msg.Source = consts.MessageSourceUser
	msg.UserId = conn.getUserId()
	msg.AdminId = service.ChatRelation().GetUserValidAdmin(ctx, msg.UserId)
	msg, err = service.ChatMessage().Insert(ctx, msg)
	if err != nil {
		return
	}
	// 发送回执
	conn.deliver(action.newReceipt(msg))
	if msg.AdminId > 0 {
		// 获取有效会话
		session, err := service.ChatSession().FirstActive(ctx, msg.UserId, msg.AdminId, nil)
		if err != nil {
			return err
		}
		// 更新有效时间
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
		adminOnline, _ := adminM.getConnInfo(ctx, msg.CustomerId, msg.AdminId)
		if adminOnline {
			g.Dump(adminOnline)
			err = userM.triggerMessageEvent(ctx, consts.AutoRuleSceneAdminOnline, msg, conn)
			if err != nil {
				return err
			}
			return adminM.deliveryMessage(ctx, msg)
		} else {
			err = adminM.handleOffline(ctx, msg, conn)
			if err != nil {
				return err
			}
		}

	} else {
		// 触发自动回复事件
		err = s.triggerMessageEvent(ctx, consts.AutoRuleSceneNotAccepted, msg, conn)
		if err != nil {
			g.Log().Errorf(ctx, "%+v", err)
		}
		var transferAdminId uint
		// 转接adminId
		transferAdminId, err = service.ChatTransfer().GetUserTransferId(ctx, conn.getCustomerId(), conn.getUserId())
		if err != nil {
			_ = service.ChatTransfer().RemoveUser(ctx, conn.getCustomerId(), conn.getUserId())
			return
		}
		var session *model.CustomerChatSession
		if transferAdminId == 0 {
			// 在代人工接入列表中
			inManual := manual.isInSet(ctx, conn.getUserId(), conn.getCustomerId())
			if inManual {
				session, err = service.ChatSession().FirstNormal(ctx, msg.UserId, 0)
				if err != nil {
					_ = manual.removeFromSet(ctx, conn.getUserId(), conn.getCustomerId())
					return
				}
				msg.SessionId = session.Id
				_, err = service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
					SessionId: msg.SessionId,
				})
				if err != nil {
					return
				}
				err = adminM.broadcastWaitingUser(ctx, conn.getCustomerId())
				if err != nil {
					return
				}
			} else {
				// 不在代人工接入列表中
				var isAutoAdd bool
				isAutoAdd, err = service.ChatSetting().GetIsAutoTransferManual(ctx, conn.getCustomerId())
				if err != nil {
					return
				}
				if isAutoAdd { // 如果自动转人工
					session, err = s.addToManual(ctx, conn.getUser())
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
					err = adminM.broadcastWaitingUser(ctx, conn.getCustomerId())
					if err != nil {
						return
					}
					err = s.broadcastQueueLocation(ctx, conn.getCustomerId())
					if err != nil {
						return
					}
				}
			}
		} else {
			session, err = service.ChatSession().FirstTransfer(ctx, conn.getUserId(), transferAdminId)
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
	if service.ChatRelation().GetUserValidAdmin(ctx, conn.getUserId()) > 0 {
		return
	}
	cacheKey := fmt.Sprintf("welcome:%d", conn.getUserId())
	val, err := gcache.Get(ctx, cacheKey)
	if err != nil {
		return
	}
	if !val.IsNil() {
		return
	}
	var rule *model.CustomerChatAutoRule
	rule, err = service.AutoRule().GetEnterRule(ctx, conn.getCustomerId())
	if err != nil {
		return
	}
	var autoMsg *model.CustomerChatAutoMessage
	autoMsg, err = service.AutoRule().GetMessage(ctx, rule)
	if err != nil {
		return
	}
	var entityMessage *model.CustomerChatMessage
	entityMessage, err = service.AutoMessage().ToChatMessage(ctx, autoMsg)
	if err != nil {
		return
	}
	entityMessage.UserId = conn.getUserId()
	_, err = service.ChatMessage().Save(ctx, entityMessage)
	if err != nil {
		return err
	}
	rule.Count++
	err = service.AutoRule().IncrTriggerCount(ctx, rule)
	if err != nil {
		return
	}
	conn.deliver(action.newReceive(entityMessage))
	err = gcache.Set(ctx, cacheKey, gtime.Now().String(), time.Minute*10)
	if err != nil {
		return
	}
	return
}

func (s *userManager) unRegisterHook(ctx context.Context, arg eventArg) error {
	adminM.noticeUserOffline(arg.conn.getUser())
	return nil
}

// 链接建立后的额外操作
// 如果已经在待接入人工列表中，则推送当前队列位置
// 如果不在待接入人工列表中且没有设置客服，则触发进入事件
func (s *userManager) onRegister(ctx context.Context, arg eventArg) error {
	adminM.noticeUserOnline(ctx, arg.conn)
	var err error
	if manual.isInSet(ctx, arg.conn.getUserId(), arg.conn.getCustomerId()) {
		err = s.noticeQueueLocation(ctx, arg.conn)
	} else {
		err = s.triggerEnterEvent(ctx, arg.conn)
	}
	return err
}

// 加入人工列表
func (s *userManager) addToManual(ctx context.Context, user iChatUser) (session *model.CustomerChatSession, err error) {
	if manual.isInSet(ctx, user.getPrimaryKey(), user.getCustomerId()) {
		err = gerror.New("is in")
		return
	}
	onlineServerCount := adminM.getOnlineTotal(user.getCustomerId())
	if onlineServerCount == 0 { // 如果没有在线客服
		rule, _ := service.AutoRule().GetSystemOne(ctx, user.getCustomerId(), consts.AutoRuleMatchAdminAllOffLine)
		if rule != nil {
			switch rule.ReplyType {
			case consts.AutoRuleReplyTypeMessage:
				autoMessage, _ := service.AutoRule().GetMessage(ctx, rule)
				if autoMessage != nil {
					message, err := service.AutoMessage().ToChatMessage(ctx, autoMessage)
					if err != nil {
						return nil, err
					}
					err = service.AutoRule().IncrTriggerCount(ctx, rule)
					if err != nil {
						return nil, err
					}
					message.UserId = user.getPrimaryKey()
					message, err = service.ChatMessage().Insert(ctx, message)
					if err != nil {
						return nil, err
					}
					err = s.deliveryLocalMessage(ctx, message)
					if err != nil {
						return nil, err
					}
					return nil, nil
				}
			}
		}
	}
	err = manual.addToSet(ctx, user.getPrimaryKey(), user.getCustomerId())
	if err != nil {
		return
	}
	session, err = service.ChatSession().FirstNormal(ctx, user.getPrimaryKey(), 0)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		} else {
			session, err = service.ChatSession().Create(ctx, user.getPrimaryKey(), user.getCustomerId(), consts.ChatSessionTypeNormal)
		}
	}
	if session == nil {
		return
	}
	message := service.ChatMessage().NewNotice(session, "正在为你转接人工客服")
	_, err = service.ChatMessage().Save(ctx, message)
	if err != nil {
		return nil, err
	}
	err = s.deliveryLocalMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return
}

// 触发事件
func (s *userManager) triggerMessageEvent(ctx context.Context, scene string, message *model.CustomerChatMessage, userConn iWsConn) (err error) {
	user := userConn.getUser()
	rules, err := service.AutoRule().AllActive(ctx, user.getCustomerId())
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
			isTransfer, err = service.ChatTransfer().IsInTransfer(ctx, user.getCustomerId(), user.getPrimaryKey())
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
			err = s.broadcastQueueLocation(ctx, message.CustomerId)
			if err != nil {
				return
			}
			err = adminM.broadcastWaitingUser(ctx, user.getCustomerId())
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
			msg, err := service.AutoMessage().ToChatMessage(ctx, autoMessage)
			if err != nil {
				return err
			}
			msg.UserId = message.UserId
			msg.SessionId = message.SessionId
			_, err = service.ChatMessage().Save(ctx, msg)
			if err != nil {
				return err
			}
			s.SendAction(action.newReceive(msg), userConn)
			_ = service.AutoRule().IncrTriggerCount(ctx, rule)
			return nil
		}
	}
	return nil
}
