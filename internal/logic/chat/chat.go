package chat

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
)

var adminM *adminManager
var userM *userManager

func init() {
	service.RegisterChat(newChat())
}

func newChat() *sChat {
	config, _ := g.Config().Get(gctx.New(), "grpc.open", false)
	cluster := config.Bool()
	m := &sChat{
		admin:   newAdminManager(cluster),
		user:    newUserManager(cluster),
		cluster: cluster,
	}
	m.run()
	return m
}

type sChat struct {
	admin   *adminManager
	user    *userManager
	cluster bool
}

func (s sChat) run() {
	s.admin.run()
	s.user.run()
}

// GetUsersWithLimitTime 获取所有user以及对应的有效期
func (s sChat) GetUsersWithLimitTime(ctx gctx.Ctx, adminId uint) (uids []uint, times []int64, err error) {
	return relation.getUsersWithLimitTime(ctx, adminId)
}

// GetUsersByAdmin 获取客服的所有用户
func (s sChat) GetUsersByAdmin(ctx gctx.Ctx, adminId uint) ([]uint, error) {
	uids, _, err := relation.getUsersWithLimitTime(ctx, adminId)
	return uids, err
}

func (s sChat) GetInvalidUsers(ctx gctx.Ctx, adminId uint) ([]uint, error) {
	return relation.getInvalidUsers(ctx, adminId)
}
func (s sChat) UpdateUserLimitTime(ctx gctx.Ctx, adminId uint, uid uint, duration int64) error {
	return relation.updateLimitTime(ctx, adminId, uid, duration)
}

func (s sChat) IsUserValid(ctx gctx.Ctx, adminId uint, uid uint) (bool, error) {
	b, err := relation.getLimitTime(ctx, adminId, uid)
	if err != nil {
		return false, err
	}
	return b > time.Now().Unix(), nil
}
func (s sChat) GetActiveUserCount(ctx gctx.Ctx, adminId uint) (uint, error) {
	return relation.getActiveCount(ctx, adminId)
}
func (s sChat) GetUserLimitTime(ctx gctx.Ctx, adminId uint, uid uint) (int64, error) {
	return relation.getLimitTime(ctx, adminId, uid)
}
func (s sChat) GetUserLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) (uint, error) {
	return relation.getLastChatTime(ctx, adminId, uid)
}

func (s sChat) RemoveUser(ctx gctx.Ctx, adminId uint, uid uint) (err error) {
	return relation.removeUser(ctx, adminId, uid)
}

func (s sChat) UpdateAdminSetting(ctx context.Context, id uint, forceLocal ...bool) error {
	return s.admin.updateSetting(ctx, id, forceLocal...)
}

// NoticeTransfer 发送转接通知
func (s sChat) NoticeTransfer(ctx context.Context, customer, admin uint, forceLocal ...bool) error {
	return s.admin.noticeUserTransfer(ctx, customer, admin, forceLocal...)
}

func (s sChat) NoticeUserOnline(ctx context.Context, uid uint, platform string, forceLocal ...bool) error {
	return s.admin.noticeUserOnline(ctx, uid, platform, forceLocal...)
}

func (s sChat) NoticeUserOffline(ctx context.Context, uid uint, forceLocal ...bool) error {
	return s.admin.noticeUserOffline(ctx, uid, forceLocal...)
}

// Accept 接入用户
func (s sChat) Accept(ctx context.Context, admin model.CustomerAdmin, sessionId uint) (u *api.ChatUser, err error) {
	session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{
		Id:         sessionId,
		CustomerId: admin.CustomerId,
	})
	if err != nil {
		return
	}
	session.User, err = service.User().Find(ctx, session.UserId)
	if err != nil {
		return
	}
	if session.CanceledAt != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该用户已被取消")
	}
	if session.AcceptedAt != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该用户已接入")
	}
	// 如果是转接
	if session.Type == consts.ChatSessionTypeTransfer {
		transfer, err := service.ChatTransfer().First(ctx, do.CustomerChatTransfers{
			ToSessionId: session.Id,
		})
		if err != nil {
			return nil, err
		}
		err = service.ChatTransfer().Accept(ctx, transfer)
		if err != nil {
			return nil, err
		}
	}
	session.AcceptedAt = gtime.Now()
	session.AdminId = admin.Id
	_, err = service.ChatSession().UpdatePri(ctx, session.Id, do.CustomerChatSessions{
		AcceptedAt: session.AcceptedAt,
		AdminId:    session.AdminId,
	})
	if err != nil {
		return
	}
	messages, err := service.ChatMessage().All(ctx, do.CustomerChatMessages{
		SessionId: session.Id,
		AdminId:   0,
		Source:    consts.MessageSourceUser,
	}, g.Slice{model.CustomerChatMessage{}.User}, "id desc")
	if err != nil {
		return
	}
	unRead := len(messages)
	// 更新未发送的消息
	updateIds := slice.Map(messages, func(index int, item *model.CustomerChatMessage) uint {
		return item.Id
	})
	_, err = service.ChatMessage().Update(ctx, do.CustomerChatMessages{
		Id: updateIds,
	}, do.CustomerChatMessages{
		AdminId: admin.Id,
	})
	if err != nil {
		return
	}
	online, platform, err := s.user.getConnInfo(ctx, session.CustomerId, session.UserId)
	if err != nil {
		return
	}
	if online {
		// 服务提醒
		chatName, _ := service.Admin().GetChatName(ctx, &admin)
		notice := service.ChatMessage().NewNotice(session,
			chatName+"为您服务")
		_, err = service.ChatMessage().Insert(ctx, notice)
		if err != nil {
			return
		}
		err = s.user.deliveryMessage(ctx, notice)
		if err != nil {
			return
		}
		var welcomeMsg *model.CustomerChatMessage
		// 欢迎语
		welcomeMsg, err = service.ChatMessage().NewWelcome(ctx, &admin)
		if err != nil {
			return nil, err
		}
		if welcomeMsg != nil {
			welcomeMsg.UserId = session.UserId
			welcomeMsg.SessionId = session.Id
			_, err = service.ChatMessage().Insert(ctx, welcomeMsg)
			if err != nil {
				return nil, err
			}
			err = s.user.deliveryMessage(ctx, welcomeMsg)
			if err != nil {
				return
			}
		}
	}
	messagesLength := len(messages)
	var lastMsg *baseApi.ChatMessage
	if messagesLength > 0 {
		lastMessage := messages[0]
		lastMsg, err = service.ChatMessage().ToApi(ctx, lastMessage)
		if err != nil {
			return nil, err
		}
	}
	err = relation.addUser(ctx, admin.Id, session.UserId)
	if err != nil {
		return
	}
	err = s.RemoveManual(ctx, session.UserId, session.CustomerId)
	if err != nil {
		return
	}
	user := &api.ChatUser{
		Id:           session.User.Id,
		Username:     session.User.Username,
		LastChatTime: gtime.Now(),
		Disabled:     false,
		Online:       online,
		LastMessage:  lastMsg,
		Unread:       uint(unRead),
		Avatar:       "",
		Platform:     platform,
	}
	return user, nil

}

func (s sChat) Register(ctx context.Context, u any, conn *websocket.Conn, platform string) error {
	switch uu := u.(type) {
	case *model.CustomerAdmin:
		e := &admin{uu}
		return s.admin.register(ctx, conn, e, platform)
	case *model.User:
		e := &user{uu}
		return s.user.register(ctx, conn, e, platform)
	}
	return gerror.NewCode(gcode.CodeBusinessValidationFailed, "无效的用户模型")
}
func (s sChat) GetConnInfo(ctx context.Context, customerId, uid uint, t string, forceLocal ...bool) (exist bool, platform string, err error) {
	if t == consts.WsTypeAdmin {
		return s.admin.getConnInfo(ctx, customerId, uid, forceLocal...)
	}
	if t == consts.WsTypeUser {
		return s.user.getConnInfo(ctx, customerId, uid, forceLocal...)
	}
	return false, "", nil
}

func (s sChat) BroadcastWaitingUser(ctx context.Context, customerId uint, forceLocal ...bool) error {
	return s.admin.broadcastWaitingUser(ctx, customerId, forceLocal...)
}
func (s sChat) BroadcastOnlineAdmins(ctx context.Context, customerId uint, forceLocal ...bool) error {
	return s.admin.broadcastOnlineAdmins(ctx, customerId, forceLocal...)
}
func (s sChat) BroadcastQueueLocation(ctx context.Context, customerId uint, forceLocal ...bool) error {
	return s.user.broadcastQueueLocation(ctx, customerId, forceLocal...)

}
func (s sChat) NoticeRate(msg *model.CustomerChatMessage) {
	s.admin.noticeRate(msg)
}
func (s sChat) NoticeRepeatConnect(ctx context.Context, uid, customerId uint, newUuid string, t string, forceLocal ...bool) error {
	if t == consts.WsTypeUser {
		return s.user.noticeRepeatConnect(ctx, uid, customerId, newUuid, forceLocal...)
	} else {
		return s.admin.noticeRepeatConnect(ctx, uid, customerId, newUuid, forceLocal...)
	}
}
func (s sChat) NoticeRead(ctx context.Context, customerId, uid uint, msgIds []uint, t string, forceLocal ...bool) error {
	if t == consts.WsTypeAdmin {
		return s.admin.noticeRead(ctx, customerId, uid, msgIds, forceLocal...)
	} else if t == consts.WsTypeUser {
		return s.user.noticeRead(ctx, customerId, uid, msgIds, forceLocal...)
	}
	return nil
}

func (s sChat) GetOnlineAdmins(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error) {
	ids, err := s.admin.getOnlineUserIds(ctx, customerId)
	if err != nil {
		return nil, err
	}
	users, err := service.Admin().All(ctx, do.Users{
		Id: ids,
	}, nil, nil)
	if err != nil {
		return nil, err
	}
	res := make([]api.ChatSimpleUser, len(users))
	for index, u := range users {
		res[index] = api.ChatSimpleUser{
			Id:       u.Id,
			Username: u.Username,
		}
	}
	return res, nil
}

func (s sChat) GetOnlineUserIds(ctx context.Context, customerId uint, types string, forceLocal ...bool) ([]uint, error) {
	if types == consts.WsTypeUser {
		return s.user.getOnlineUserIds(ctx, customerId, forceLocal...)
	} else {
		return s.admin.getOnlineUserIds(ctx, customerId, forceLocal...)
	}
}

func (s sChat) GetOnlineUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error) {
	ids, err := s.GetOnlineUserIds(ctx, customerId, consts.WsTypeUser)
	if err != nil {
		return nil, err
	}
	users, err := service.User().All(ctx, do.Users{
		Id: ids,
	}, nil, nil)
	if err != nil {
		return nil, err
	}
	res := make([]api.ChatSimpleUser, len(users))
	for index, u := range users {
		res[index] = api.ChatSimpleUser{
			Id:       u.Id,
			Username: u.Username,
		}
	}
	return res, nil
}
func (s sChat) GetWaitingUsers(ctx context.Context, customerId uint) (res []api.ChatSimpleUser, err error) {
	ids, err := manual.getAllList(ctx, customerId)
	if err != nil {
		return
	}
	users, err := service.User().All(ctx, g.Map{
		"id": ids,
	}, nil, nil)

	res = slice.Map(users, func(index int, item *model.User) api.ChatSimpleUser {
		return api.ChatSimpleUser{
			Id:       item.Id,
			Username: item.Username,
		}
	})
	return
}

func (s sChat) RemoveManual(ctx context.Context, uid uint, customerId uint) error {
	return manual.removeFromSet(ctx, uid, customerId)
}

func (s sChat) DeliveryMessage(ctx context.Context, msgId uint, types string, forceLocal ...bool) error {
	msg, err := service.ChatMessage().Find(ctx, msgId)
	if err != nil {
		return err
	}
	if types == consts.WsTypeAdmin {
		return s.admin.deliveryMessage(ctx, msg, forceLocal...)
	} else {
		return s.user.deliveryMessage(ctx, msg, forceLocal...)
	}
}
