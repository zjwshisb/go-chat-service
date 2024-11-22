package chat

import (
	"context"
	"errors"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
	"strings"
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

func (s sChat) UpdateAdminSetting(customerId uint, setting *entity.CustomerAdminChatSettings) {
	s.admin.updateSetting(customerId, setting)
}

func (s sChat) NoticeTransfer(ctx context.Context, customer, admin uint) error {
	return s.admin.noticeUserTransfer(ctx, customer, admin)
}

func (s sChat) Accept(ctx context.Context, admin model.CustomerAdmin, sessionId uint) (u *api.ChatUser, err error) {
	session := &model.CustomerChatSession{}
	err = dao.CustomerChatSessions.Ctx(ctx).
		Where("customer_id", admin.CustomerId).WithAll().
		WherePri(sessionId).
		Scan(session)
	if err != nil {
		return
	}
	if session.CanceledAt != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该用户已被取消")
	}
	if session.AcceptedAt != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该用户已接入")
	}
	if session.Type == consts.ChatSessionTypeTransfer {
		transfer, _ := service.ChatTransfer().First(ctx, do.CustomerChatTransfers{
			ToSessionId: session.Id,
			AcceptedAt:  nil,
			CanceledAt:  nil,
		})
		if transfer == nil {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该转接已被接入")
		}
		err = service.ChatTransfer().Accept(ctx, transfer)
		if err != nil {
			return
		}
	}
	session.AcceptedAt = gtime.New()
	session.AdminId = admin.Id
	_, err = service.ChatSession().Save(ctx, session)
	if err != nil {
		return
	}
	unRead, _ := dao.CustomerChatMessages.Ctx(ctx).
		Where("session_id", session.Id).
		Where("admin_id", 0).
		Where("source", consts.MessageSourceUser).Count()
	// 更新未发送的消息
	_, err = dao.CustomerChatMessages.Ctx(ctx).Where("session_id", session.Id).
		Where("admin_id", 0).
		Where("source", consts.MessageSourceUser).Data(g.Map{
		"admin_id": admin.Id,
	}).Update()
	if err != nil {
		return
	}
	userConn, exist := s.user.GetConn(session.CustomerId, session.UserId)
	platform := ""
	if exist {
		// 服务提醒
		platform = userConn.GetPlatform()
		chatName, _ := service.Admin().GetChatName(ctx, &admin)
		notice := service.ChatMessage().NewNotice(session,
			chatName+"为您服务")
		_ = service.ChatMessage().SaveWithUpdate(ctx, notice)
		s.user.SendAction(newReceiveAction(notice), userConn)
		// 欢迎语
		welcomeMsg := service.ChatMessage().NewWelcome(&admin)
		if welcomeMsg != nil {
			welcomeMsg.UserId = session.UserId
			welcomeMsg.SessionId = session.Id
			_ = service.ChatMessage().SaveWithUpdate(ctx, welcomeMsg)
			action := newReceiveAction(welcomeMsg)
			s.user.SendAction(action, userConn)
		}
	}
	lastMessage, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		UserId: session.UserId,
		Source: []int{consts.MessageSourceUser, consts.MessageSourceAdmin},
	}, "order desc")
	var lastMsg *api.ChatMessage
	if err == nil {
		v, err := service.ChatMessage().ToApi(ctx, *lastMessage)
		if err != nil {
			return nil, err
		}
		lastMsg = &v
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

func (s sChat) Register(ctx context.Context, u any, conn *websocket.Conn) error {
	request := ghttp.RequestFromCtx(ctx)
	userAgent := strings.ToLower(request.UserAgent())
	wechatAgent := []string{"micromessenger", "wechatdevtools"}
	isWeapp := false
	for _, s := range wechatAgent {
		if strings.Contains(userAgent, s) {
			isWeapp = true
			break
		}
	}
	platform := ""
	switch u.(type) {
	case *model.CustomerAdmin:
		uu, _ := u.(*model.CustomerAdmin)
		e := &admin{
			uu,
		}
		if isWeapp {
			platform = consts.WebsocketPlatformMp
		} else {
			platform = consts.WebsocketPlatformWeb
		}
		s.admin.Register(conn, e, platform)
		return nil
	case *entity.Users:
		uu, _ := u.(*entity.Users)
		e := &user{
			uu,
		}
		if isWeapp {
			platform = consts.WebsocketPlatformMp
		} else {
			platform = consts.WebsocketPlatformH5
		}
		s.user.Register(conn, e, platform)
		return nil
	}
	return errors.New("无效的用户模型")
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

func (s sChat) GetOnlineCount(ctx context.Context, customerId uint) (res api.ChatOnlineCount, err error) {
	waiting, err := manual.getCount(ctx, customerId)
	if err != nil {
		return
	}
	res = api.ChatOnlineCount{
		Admin:   s.admin.GetOnlineTotal(customerId),
		User:    s.user.GetOnlineTotal(customerId),
		Waiting: waiting,
	}
	return
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

func (s sChat) GetOnlineAdmin(customerId uint) []api.ChatSimpleUser {
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

func (s sChat) GetOnlineUser(customerId uint) []api.ChatSimpleUser {
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

func (s sChat) RemoveManual(ctx context.Context, uid uint, customerId uint) error {
	return manual.removeFromSet(ctx, uid, customerId)
}
