package chat

import (
	"context"
	"gf-chat/api/v1"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
)

var adminM *adminManager
var userM *userManager

func init() {
	service.RegisterChat(newChat())
}

func newChat() *sChat {
	m := &sChat{
		admin: newAdminManager(),
		user:  newUserManager(),
	}
	m.run()
	return m
}

type sChat struct {
	admin *adminManager
	user  *userManager
}

func (s sChat) run() {
	s.admin.run()
	s.user.run()
}

func (s sChat) UpdateAdminSetting(ctx context.Context, admin *model.CustomerAdmin) {
	s.admin.updateSetting(ctx, admin)
}

// NoticeTransfer 发送转接通知
func (s sChat) NoticeTransfer(ctx context.Context, customer, admin uint) error {
	return s.admin.noticeUserTransfer(ctx, customer, admin)
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
	userConn, online := s.user.GetConn(session.CustomerId, session.UserId)
	platform := ""
	if online {
		// 服务提醒
		platform = userConn.GetPlatform()
		chatName, _ := service.Admin().GetChatName(ctx, &admin)
		notice := service.ChatMessage().NewNotice(session,
			chatName+"为您服务")
		_, err = service.ChatMessage().Insert(ctx, notice)
		if err != nil {
			return
		}
		s.user.SendAction(action.newReceive(notice), userConn)
		// 欢迎语
		welcomeMsg, err := service.ChatMessage().NewWelcome(ctx, &admin)
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
			action := action.newReceive(welcomeMsg)
			s.user.SendAction(action, userConn)
		}
	}
	messagesLength := len(messages)
	var lastMsg *v1.ChatMessage
	if messagesLength > 0 {
		lastMessage := messages[0]
		v, err := service.ChatMessage().ToApi(ctx, lastMessage)
		if err != nil {
			return nil, err
		}
		lastMsg = v
	}
	err = service.ChatRelation().AddUser(ctx, admin.Id, session.UserId)
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
		Online:       userM.IsOnline(session.CustomerId, session.UserId),
		LastMessage:  lastMsg,
		Unread:       uint(unRead),
		Avatar:       "",
		Platform:     platform,
	}
	return user, nil

}

func (s sChat) Register(ctx context.Context, u any, conn *websocket.Conn, platform string) error {
	switch u.(type) {
	case *model.CustomerAdmin:
		uu, _ := u.(*model.CustomerAdmin)
		e := &admin{
			uu,
		}
		return s.admin.Register(ctx, conn, e, platform)
	case *model.User:
		uu, _ := u.(*model.User)
		e := &user{
			uu,
		}
		return s.user.Register(ctx, conn, e, platform)
	}
	return gerror.New("无效的用户模型")
}

func (s sChat) IsOnline(customerId uint, uid uint, t string) bool {
	if t == "user" {
		return s.user.IsOnline(customerId, uid)
	}
	if t == "admin" {
		return s.admin.IsOnline(customerId, uid)
	}
	return false
}

func (s sChat) BroadcastWaitingUser(ctx context.Context, customerId uint) error {
	return s.admin.broadcastWaitingUser(ctx, customerId)
}

func (s sChat) GetPlatform(customerId, uid uint, t string) string {
	var conn iWsConn
	var online bool
	if t == "admin" {
		conn, online = s.admin.GetConn(customerId, uid)
	}
	if t == "user" {
		conn, online = s.user.GetConn(customerId, uid)
	}
	if online {
		return conn.GetPlatform()
	}
	return ""
}

func (s sChat) NoticeRate(msg *model.CustomerChatMessage) {
	s.admin.noticeRate(msg)
}

func (s sChat) NoticeUserRead(customerId, uid uint, msgIds []uint) {
	s.admin.NoticeRead(customerId, uid, msgIds)
}

func (s sChat) NoticeAdminRead(customerId, uid uint, msgIds []uint) {
	s.user.NoticeRead(customerId, uid, msgIds)
}

func (s sChat) Transfer(ctx context.Context, fromAdmin *model.CustomerAdmin, toId uint, userId uint, remark string) (err error) {
	user, err := service.User().First(ctx, do.Users{
		CustomerId: fromAdmin.CustomerId,
		Id:         userId,
	})
	if err != nil {
		return
	}
	admin, err := service.Admin().First(ctx, do.CustomerAdmins{
		CustomerId: fromAdmin.CustomerId,
		Id:         toId,
	})
	if err != nil {
		return err
	}
	isValid := service.ChatRelation().IsUserValid(ctx, fromAdmin.Id, user.Id)
	if !isValid {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "用户已失效，无法转接")
	}
	return service.ChatTransfer().Create(ctx, fromAdmin.Id, admin.Id, userId, remark)
}

func (s sChat) GetOnlineAdmins(customerId uint) []api.ChatSimpleUser {
	conns := s.admin.GetAllConn(customerId)
	res := make([]api.ChatSimpleUser, len(conns))
	for index, c := range conns {
		res[index] = api.ChatSimpleUser{
			Id:       c.GetUserId(),
			Username: c.GetUser().GetUsername(),
		}
	}
	return res
}

func (s sChat) GetOnlineUsers(customerId uint) []api.ChatSimpleUser {
	conns := s.user.GetAllConn(customerId)
	res := make([]api.ChatSimpleUser, len(conns))
	for index, c := range conns {
		res[index] = api.ChatSimpleUser{
			Id:       c.GetUserId(),
			Username: c.GetUser().GetUsername(),
		}
	}
	return res
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
