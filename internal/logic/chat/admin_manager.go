package chat

import (
	"context"
	"gf-chat/api/backend/v1"
	grpc "gf-chat/api/chat/v1"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"sort"
	"time"
)

func newAdminManager(cluster bool) *adminManager {
	adminM = &adminManager{
		newManager(10, 100, time.Minute, cluster, consts.WsTypeAdmin),
	}
	adminM.on(eventRegister, adminM.onRegister)
	adminM.on(eventUnRegister, adminM.onUnRegister)
	adminM.on(eventMessage, adminM.onMessage)
	return adminM
}

type adminManager struct {
	*manager
}

// 投递消息
func (m *adminManager) deliveryMessage(ctx context.Context, msg *model.CustomerChatMessage, forceLocal ...bool) error {
	userLocal, server, err := m.isUserLocal(ctx, msg.AdminId)
	if err != nil {
		return err
	}
	if m.isCallLocal(forceLocal...) || userLocal {
		adminConn, exist := m.getConn(msg.CustomerId, msg.AdminId)
		if exist { // admin在线
			adminConn.deliver(action.newReceive(msg))
		}
		return nil
	}
	if server != "" {
		rpcClient := service.Grpc().Client(ctx, server)
		if rpcClient != nil {
			_, err = rpcClient.SendMessage(ctx, &grpc.SendMessageRequest{
				MsgId: uint32(msg.Id),
				Type:  consts.WsTypeAdmin,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 处理离线消息
func (m *adminManager) handleOffline(ctx context.Context, msg *model.CustomerChatMessage, userConn iWsConn) error {
	err := userM.triggerMessageEvent(ctx, consts.AutoRuleSceneAdminOffline, msg, userConn)
	if err != nil {
		return err
	}
	admin, err := service.Admin().First(ctx, do.CustomerAdmins{Id: msg.AdminId})
	if err != nil {
		return err
	}
	message, err := service.ChatMessage().NewOffline(ctx, admin)
	if err != nil {
		return err
	}
	if message != nil {
		message.UserId = msg.UserId
		message.SessionId = msg.SessionId
		message, err = service.ChatMessage().Insert(ctx, message)
		if err != nil {
			return err
		}
		err = userM.deliveryMessage(ctx, message)
		if err != nil {
			return err
		}
	}
	return nil
}

// 处理消息
func (m *adminManager) onMessage(ctx context.Context, arg eventArg) error {
	msg := arg.msg
	conn := arg.conn
	if msg.UserId > 0 {
		if !service.ChatRelation().IsUserValid(ctx, conn.getUserId(), msg.UserId) {
			conn.deliver(action.newErrorMessage("该用户已失效，无法发送消息"))
			return gerror.NewCode(gcode.CodeValidationFailed, "该用户已失效，无法发送消息")
		}
		session, err := service.ChatSession().FirstActive(ctx, msg.UserId, conn.getUserId(), nil)
		if err != nil {
			conn.deliver(action.newErrorMessage("无效的用户"))
			return gerror.NewCode(gcode.CodeValidationFailed, "无效的用户")
		}
		msg.AdminId = conn.getUserId()
		msg.Source = consts.MessageSourceAdmin
		msg.SessionId = session.Id
		msg, err := service.ChatMessage().Insert(ctx, msg)
		if err != nil {
			return err
		}
		_ = service.ChatRelation().UpdateUser(ctx, msg.AdminId, msg.UserId)
		// 服务器回执d
		conn.deliver(action.newReceipt(msg))
		return userM.deliveryMessage(ctx, msg)
	}
	return nil
}

func (m *adminManager) onRegister(ctx context.Context, arg eventArg) error {
	err := m.broadcastOnlineAdmins(ctx, arg.conn.getCustomerId())
	if err != nil {
		return err
	}
	err = m.broadcastWaitingUser(ctx, arg.conn.getCustomerId())
	if err != nil {
		return err
	}
	err = m.noticeUserTransfer(ctx, arg.conn.getCustomerId(), arg.conn.getUserId())
	if err != nil {
		return err
	}
	return nil
}

// conn断开连接后，更新admin的最后在线时间
func (m *adminManager) onUnRegister(ctx context.Context, arg eventArg) error {
	err := service.Admin().UpdateLastOnline(ctx, arg.conn.getUserId())
	if err != nil {
		return err
	}
	err = m.broadcastOnlineAdmins(ctx, arg.conn.getCustomerId())
	if err != nil {
		return err
	}
	return nil
}

func (m *adminManager) broadcastWaitingUser(ctx context.Context, customerId uint, forceLocal ...bool) (err error) {
	if m.isCallLocal(forceLocal...) {
		sessions, err := service.ChatSession().GetUnAccepts(ctx, customerId)
		if err != nil {
			return err
		}
		sessionIds := slice.Map(sessions, func(index int, item *model.CustomerChatSession) uint {
			return item.Id
		})
		userMap := make(map[uint]*v1.ChatWaitingUser)
		messages, err := service.ChatMessage().All(ctx, do.CustomerChatMessages{
			Source:    consts.MessageSourceUser,
			SessionId: sessionIds,
		}, nil, "id desc")
		if err != nil {
			return err
		}
		for _, session := range sessions {
			userMap[session.UserId] = &v1.ChatWaitingUser{
				Username:     session.User.Username,
				Avatar:       "",
				UserId:       session.User.Id,
				MessageCount: 0,
				Description:  "",
				Messages:     make([]v1.ChatSimpleMessage, 0),
				LastTime:     session.QueriedAt,
				SessionId:    session.Id,
			}
		}
		for _, m := range messages {
			userMap[m.UserId].Messages = append(userMap[m.UserId].Messages, v1.ChatSimpleMessage{
				Id:      m.Id,
				Type:    m.Type,
				Time:    m.ReceivedAt,
				Content: m.Content,
			})
			userMap[m.UserId].MessageCount += 1
		}

		waitingUser := maputil.Values(userMap)
		sort.Slice(waitingUser, func(i, j int) bool {
			return waitingUser[i].LastTime.Unix() > waitingUser[j].LastTime.Unix()
		})
		adminConns := m.getAllConn(customerId)
		act := action.newWaitingUsers(waitingUser)
		m.SendAction(act, adminConns...)
		return nil
	} else {
		err = service.Grpc().CallAll(ctx, func(client grpc.ChatClient) {
			_, err := client.BroadcastWaitingUser(ctx, &grpc.BroadcastWaitingUserRequest{
				CustomerId: uint32(customerId),
			})
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
			}
		})
		return nil
	}
}
func (m *adminManager) broadcastOnlineAdmins(ctx context.Context, customerId uint, forceLocal ...bool) error {
	if m.isCallLocal(forceLocal...) {
		admins, err := service.Admin().All(ctx, do.CustomerAdmins{
			CustomerId: customerId,
		}, nil, nil)
		if err != nil {
			return err
		}
		data := make([]v1.ChatCustomerAdmin, 0, len(admins))
		for _, c := range admins {
			conn, online := m.getConn(customerId, c.Id)
			platform := ""
			if online {
				platform = conn.getPlatform()
			}
			data = append(data, v1.ChatCustomerAdmin{
				Username:      c.Username,
				Online:        online,
				Id:            c.Id,
				AcceptedCount: service.ChatRelation().GetActiveCount(ctx, c.Id),
				Platform:      platform,
			})
		}
		conns := m.getAllConn(customerId)
		m.SendAction(action.newAdmins(data), conns...)
		return nil
	} else {
		err := service.Grpc().CallAll(ctx, func(client grpc.ChatClient) {
			_, err := client.BroadcastOnlineAdmins(ctx, &grpc.BroadcastOnlineAdminsRequest{
				CustomerId: uint32(customerId),
			})
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
			}
		})
		return err
	}
}

func (m *adminManager) noticeRate(message *model.CustomerChatMessage) {
	action := action.newRateAction(message)
	conn, exist := m.getConn(message.CustomerId, message.AdminId)
	if exist {
		conn.deliver(action)
	}
}

func (m *adminManager) noticeUserOffline(ctx context.Context, uid uint, forceLocal ...bool) (err error) {
	adminId, err := service.ChatRelation().GetUserValidAdmin(ctx, uid)
	if err != nil {
		return
	}
	if adminId > 0 {
		adminModel, err := service.Admin().Find(gctx.New(), adminId)
		if err != nil {
			return err
		}
		adminLocal, server, err := m.isUserLocal(ctx, adminModel.Id)
		if m.isCallLocal(forceLocal...) || adminLocal {
			conn, exist := m.getConn(adminModel.CustomerId, adminModel.Id)
			if exist {
				m.SendAction(action.newUserOffline(uid), conn)
			}
			return nil
		} else if server != "" {
			rpcClient := service.Grpc().Client(ctx, server)
			if rpcClient != nil {
				_, err = rpcClient.NoticeUserOffline(ctx, &grpc.NoticeUserOfflineRequest{
					UserId: uint32(uid),
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *adminManager) noticeUserOnline(ctx context.Context, uid uint, platform string, forceLocal ...bool) (err error) {
	adminId, err := service.ChatRelation().GetUserValidAdmin(ctx, uid)
	if adminId > 0 {
		adminModel, err := service.Admin().First(ctx, do.CustomerAdmins{Id: adminId})
		if err != nil {
			return err
		}
		adminLocal, server, err := m.isUserLocal(ctx, adminModel.Id)
		if err != nil {
			return err
		}
		if m.isCallLocal(forceLocal...) || adminLocal {
			conn, exist := m.getConn(adminModel.CustomerId, adminModel.Id)
			if exist {
				m.SendAction(action.newUserOnline(uid, platform), conn)
			}
			return nil
		} else if server != "" {
			rpcClient := service.Grpc().Client(ctx, server)
			if rpcClient != nil {
				_, err = rpcClient.NoticeUserOnline(ctx, &grpc.NoticeUserOnlineRequest{
					UserId:   uint32(uid),
					Platform: platform,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *adminManager) noticeUserTransfer(ctx context.Context, customerId, adminId uint, forceLocal ...bool) error {
	userLocal, server, err := m.isUserLocal(ctx, adminId)
	if err != nil {
		return err
	}
	if m.isCallLocal(forceLocal...) || userLocal {
		client, exist := m.getConn(customerId, adminId)
		if exist {
			transfers, err := service.ChatTransfer().All(ctx, g.Map{
				"to_admin_id":         adminId,
				"accepted_at is null": nil,
				"canceled_at is null": nil,
			}, g.Slice{
				model.CustomerChatTransfer{}.ToAdmin,
				model.CustomerChatTransfer{}.FormAdmin,
				model.CustomerChatTransfer{}.ToSession,
				model.CustomerChatTransfer{}.User,
			}, nil)
			if err != nil {
				return err
			}
			data := slice.Map(transfers, func(index int, item *model.CustomerChatTransfer) v1.ChatTransfer {
				return service.ChatTransfer().ToApi(item)
			})
			client.deliver(action.newUserTransfer(data))
		}
		return nil
	}
	if server != "" {
		rpcClient := service.Grpc().Client(ctx, server)
		if rpcClient != nil {
			_, err = rpcClient.NoticeTransfer(ctx, &grpc.NoticeTransferRequest{
				CustomerId: uint32(customerId),
				AdminId:    uint32(adminId),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateSetting 更新设置
func (m *adminManager) updateSetting(ctx context.Context, id uint, forceLocal ...bool) error {
	u, err := service.Admin().Find(ctx, id)
	if err != nil {
		return err
	}
	userLocal, server, _ := m.isUserLocal(ctx, id)
	if m.isCallLocal(forceLocal...) || userLocal {
		conn, exist := m.getConn(u.CustomerId, u.Id)
		if exist {
			iu, ok := conn.getUser().(*admin)
			if ok {
				setting, err := service.Admin().FindSetting(ctx, u.Id, true)
				if err != nil {
					return err
				}
				iu.Entity.Setting = setting
			}
		}
		return nil
	}
	if server != "" {
		rpcClient := service.Grpc().Client(ctx, server)
		if rpcClient != nil {
			_, err := rpcClient.UpdateAdminSetting(ctx, &grpc.UpdateAdminSettingRequest{
				Id: uint32(u.Id),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil

}
