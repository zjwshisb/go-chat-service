package chat

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	grpc "gf-chat/api/chat/v1"
	"gf-chat/internal/cache"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
)

func newUserManager(cluster bool) *userManager {
	userM = &userManager{
		newManager(10, time.Minute, cluster, consts.WsTypeUser),
	}
	userM.on(eventRegister, userM.onRegister)
	userM.on(eventUnRegister, userM.onUnRegister)
	userM.on(eventMessage, userM.onMessage)
	return userM
}

type userManager struct {
	*manager
}

// deliveryMessage delivers a message to the appropriate user connection
// Returns error if delivery fails
// deliveryMessage delivers a chat message to a user's websocket connection.
// It handles both local and remote message delivery in a clustered environment.
//
// Parameters:
// - ctx: The context for the operation
// - message: The chat message to deliver containing user ID, customer ID etc.
// - forceLocal: Optional bool to force local delivery even in clustered mode
//
// The method:
// 1. Checks if the target user is connected locally or on a remote server
// 2. For local users, delivers directly to their websocket connection
// 3. For remote users, forwards the message via gRPC to the appropriate server
// 4. Handles special message types like ratings that require additional processing
//
// Returns error if message delivery fails
func (s *userManager) deliveryMessage(ctx context.Context, message *model.CustomerChatMessage, forceLocal ...bool) error {
	userLocal, server, err := s.isUserLocal(ctx, message.UserId)
	if err != nil {
		return err
	}
	if s.isCallLocal(forceLocal...) || userLocal {
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
	if server != "" {
		rpcClient := service.Grpc().Client(ctx, server)
		if rpcClient != nil {
			_, err = rpcClient.SendMessage(ctx, &grpc.SendMessageRequest{
				MsgId: uint32(message.Id),
				Type:  consts.WsTypeUser,
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// noticeQueueLocation 等待人数
func (s *userManager) noticeQueueLocation(ctx context.Context, conn iWsConn) (err error) {
	isSHowQueue, err := service.ChatSetting().GetIsUserShowQueue(ctx, conn.getCustomerId())
	if err != nil {
		return
	}
	if !isSHowQueue {
		return
	}
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

// broadcastQueueLocation broadcasts the current queue position to all connected users for a given customer ID.
// It checks if queue display is enabled for the customer, then either:
// 1. For local connections: Notifies all connected users of their position in queue
// 2. For cluster mode: Broadcasts the queue update to all servers via gRPC
//
// Parameters:
//   - ctx: The context for the request
//   - customerId: ID of the customer whose queue to broadcast
//   - forceLocal: Optional parameter to force local-only processing even in cluster mode
//
// Returns error if the broadcast fails
func (s *userManager) broadcastQueueLocation(ctx context.Context, customerId uint, forceLocal ...bool) (err error) {
	isSHowQueue, err := service.ChatSetting().GetIsUserShowQueue(ctx, customerId)
	if err != nil {
		return
	}
	if !isSHowQueue {
		return nil
	}
	if s.isCallLocal(forceLocal...) {
		conns := s.getAllConn(customerId)
		for _, conn := range conns {
			in, err := manual.isInSet(ctx, conn.getUserId(), conn.getCustomerId())
			if err != nil {
				return err
			}
			if in {
				err = s.noticeQueueLocation(ctx, conn)
				if err != nil {
					return err
				}
			}
		}
		return
	} else {
		err = service.Grpc().CallAll(ctx, func(client grpc.ChatClient) {
			_, err = client.BroadcastQueueLocation(ctx, &grpc.BroadcastQueueLocationRequest{
				CustomerId: uint32(customerId),
			})
			if err != nil {
				log.Errorf(ctx, "%+v", err)
			}
		})
		if err != nil {
			return
		}
	}
	return
}

// onMessage handles incoming messages from users
// It processes the message by:
// 1. Setting message source and user ID
// 2. Getting the valid admin for the user
// 3. Inserting the message into storage
// 4. Sending receipt to user
// 5. Handling message delivery to admin if online, or offline handling if not
//
// Parameters:
//   - ctx: The context for the request
//   - arg: Event arguments containing the message and connection
//
// Returns error if message processing fails
func (s *userManager) onMessage(ctx context.Context, arg eventArg) (err error) {
	msg := arg.msg
	conn := arg.conn
	msg.Source = consts.MessageSourceUser
	msg.UserId = conn.getUserId()
	msg.AdminId, err = relation.getUserValidAdmin(ctx, msg.UserId)
	if err != nil {
		return err
	}
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
		err = relation.updateUser(ctx, msg.AdminId, msg.UserId)
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
		adminOnline, _, err := adminM.getConnInfo(ctx, msg.CustomerId, msg.AdminId)
		if err != nil {
			return err
		}
		if adminOnline {
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
			log.Errorf(ctx, "%+v", err)
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
			inManual, err := manual.isInSet(ctx, conn.getUserId(), conn.getCustomerId())
			if err != nil {
				return err
			}
			if inManual {
				session, err = service.ChatSession().FirstNormal(ctx, msg.UserId, 0)
				if err != nil {
					_ = manual.removeFromSet(ctx, conn.getUserId(), conn.getCustomerId())
					return err
				}
				msg.SessionId = session.Id
				_, err = service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
					SessionId: msg.SessionId,
				})
				if err != nil {
					return err
				}
				err = adminM.broadcastWaitingUser(ctx, conn.getCustomerId())
				if err != nil {
					return err
				}
			} else {
				// 不在代人工接入列表中
				var isAutoAdd bool
				isAutoAdd, err = service.ChatSetting().GetIsAutoTransferManual(ctx, conn.getCustomerId())
				if err != nil {
					return err
				}
				if isAutoAdd { // 如果自动转人工
					session, err = s.addToManual(ctx, conn.getUser())
					if err != nil {
						return err
					}
					if session != nil {
						msg.SessionId = session.Id
						_, err = service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
							SessionId: msg.SessionId,
						})
					}
					if err != nil {
						return err
					}
					err = adminM.broadcastWaitingUser(ctx, conn.getCustomerId())
					if err != nil {
						return err
					}
					err = s.broadcastQueueLocation(ctx, conn.getCustomerId())
					if err != nil {
						return err
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
	adminId, err := relation.getUserValidAdmin(ctx, conn.getUserId())
	if adminId > 0 {
		return
	}
	cacheKey := fmt.Sprintf("welcome:%d", conn.getUserId())
	val, err := cache.Def.Get(ctx, cacheKey)
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

func (s *userManager) onUnRegister(ctx context.Context, arg eventArg) error {
	return adminM.noticeUserOffline(ctx, arg.conn.getUser().getPrimaryKey())
}

// 链接建立后的额外操作
// 如果已经在待接入人工列表中，则推送当前队列位置
// 如果不在待接入人工列表中且没有设置客服，则触发进入事件
func (s *userManager) onRegister(ctx context.Context, arg eventArg) (err error) {
	err = adminM.noticeUserOnline(ctx, arg.conn.getUserId(), arg.conn.getPlatform())
	if err != nil {
		return nil
	}
	in, err := manual.isInSet(ctx, arg.conn.getUserId(), arg.conn.getCustomerId())
	if err != nil {
		return err
	}
	if in {
		err = s.noticeQueueLocation(ctx, arg.conn)
	} else {
		err = s.triggerEnterEvent(ctx, arg.conn)
	}
	return err
}

// addToManual adds a user to the manual chat queue and handles the initial setup
// It checks if the user is already in queue, finds available admins, and handles offline admin scenarios
// Parameters:
//   - ctx: The context for the operation
//   - user: The chat user to add to manual queue
//
// Returns:
//   - session: The created chat session if successful
//   - error: Error if operation fails
func (s *userManager) addToManual(ctx context.Context, user iChatUser) (session *model.CustomerChatSession, err error) {
	in, err := manual.isInSet(ctx, user.getPrimaryKey(), user.getCustomerId())
	if err != nil {
		return
	}
	if in {
		err = gerror.NewCode(gcode.CodeBusinessValidationFailed, "is in")
		return
	}
	onlineAdminIds, err := adminM.getOnlineUserIds(ctx, user.getCustomerId())
	if err != nil {
		return
	}
	if len(onlineAdminIds) == 0 { // 如果没有在线客服
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
					err = s.deliveryMessage(ctx, message, true)
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
			if err != nil {
				return nil, err
			}
		}
	}
	if session == nil {
		return
	}
	message := service.ChatMessage().NewNotice(session, "正在为你转接人工客服")
	message, err = service.ChatMessage().Insert(ctx, message)
	if err != nil {
		return nil, err
	}
	err = s.deliveryMessage(ctx, message, true)
	if err != nil {
		return nil, err
	}

	return
}

// triggerMessageEvent handles triggering auto-reply rules based on message events.
// It checks if any auto-reply rules match the given scene and message content,
// and executes the appropriate action (transfer to manual service or send auto-reply).
//
// Parameters:
//   - ctx: The context for the operation
//   - scene: The scene identifier for matching rules (e.g. admin online/offline)
//   - message: The chat message that triggered the event
//   - userConn: The websocket connection of the user
//
// Returns error if any operation fails during rule processing
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
