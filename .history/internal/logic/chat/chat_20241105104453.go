package chat

import (
	"context"
	"database/sql"
	"errors"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

var adminM *adminManager
var userM *userManager

func init() {
	adminM = &adminManager{
		&manager{
			ShardCount:   10,
			ConnMessages: make(chan *chatConnMessage, 100),
			Types:        TypeAdmin,
		},
	}
	adminM.OnRegister = adminM.registerHook
	adminM.OnUnRegister = adminM.unregisterHook
	adminM.run()

	userM = &userManager{
		&manager{
			ShardCount:   10,
			ConnMessages: make(chan *chatConnMessage, 100),
			Types:        "user",
		},
	}
	userM.OnRegister = userM.registerHook
	userM.OnUnRegister = userM.unRegisterHook
	userM.run()
	service.RegisterChat(&sChat{
		admin: adminM,
		user:  userM,
	})
}

type sChat struct {
	admin *adminManager
	user  *userManager
}

func (s sChat) UpdateAdminSetting(customerId int, setting *entity.CustomerAdminChatSettings) {
	s.admin.updateSetting(customerId, setting)
}

func (s sChat) NoticeTransfer(customer, admin int) {
	s.admin.noticeUserTransfer(customer, admin)
}

func (s sChat) Accept(admin entity.CustomerAdmins, sessionId uint64) (*chat.User, error) {
	session := &relation.CustomerChatSessions{}
	ctx := gctx.New()
	err := dao.CustomerChatSessions.Ctx(ctx).
		Where("customer_id", admin.CustomerId).WithAll().
		WherePri(sessionId).
		Scan(session)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	if session.CanceledAt > 0 {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该用户已被取消")
	}
	if session.AcceptedAt > 0 {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该用户已接入")
	}
	if session.Type == consts.ChatSessionTypeTransfer {
		transfer := service.ChatTransfer().FirstEntity(do.CustomerChatTransfers{
			ToSessionId: session.Id,
			AcceptedAt:  0,
			CanceledAt:  0,
		})
		if transfer == nil {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该转接已被接入")
		}
		service.ChatTransfer().Accept(transfer)
	}
	session.AcceptedAt = time.Now().Unix()
	session.AdminId = gconv.Int(admin.Id)
	dao.CustomerChatSessions.Ctx(ctx).Save(session)
	unRead, _ := dao.CustomerChatMessages.Ctx(ctx).
		Where("session_id", session.Id).
		Where("admin_id", 0).
		Where("source", consts.MessageSourceUser).Count()
	// 更新未发送的消息
	dao.CustomerChatMessages.Ctx(ctx).Where("session_id", session.Id).
		Where("admin_id", 0).
		Where("source", consts.MessageSourceUser).Data(g.Map{
		"admin_id": admin.Id,
	}).Update()

	userConn, exist := s.user.GetConn(session.CustomerId, session.UserId)
	platform := ""
	if exist {
		// 服务提醒
		platform = userConn.GetPlatform()
		notice := service.ChatMessage().NewNotice(&session.CustomerChatSessions,
			service.Admin().GetChatName(&admin)+"为您服务")
		service.ChatMessage().SaveOne(notice)
		relationNotice := service.ChatMessage().EntityToRelation(notice)
		s.user.SendAction(service.Action().NewReceiveAction(relationNotice), userConn)
		// 欢迎语
		adminRelation := service.Admin().EntityToRelation(&admin)
		welcomeMsg := service.ChatMessage().NewWelcome(adminRelation)
		if welcomeMsg != nil {
			welcomeMsg.UserId = session.UserId
			welcomeMsg.SessionId = session.Id
			service.ChatMessage().SaveRelationOne(welcomeMsg)
			action := service.Action().NewReceiveAction(welcomeMsg)
			s.user.SendAction(action, userConn)
		}
	}
	lastMessage := &relation.CustomerChatMessages{}
	err = dao.CustomerChatMessages.Ctx(ctx).
		Where("user_id", session.UserId).
		OrderDesc("id").
		WithAll().
		Where("source in (?)", []int{consts.MessageSourceUser, consts.MessageSourceAdmin}).
		Scan(lastMessage)
	var lastMsg *chat.Message
	if err != sql.ErrNoRows {
		v := service.ChatMessage().RelationToChat(*lastMessage)
		lastMsg = &v
	}
	service.ChatRelation().AddUser(gconv.Int(admin.Id), session.UserId)
	service.ChatManual().Remove(session.UserId, session.CustomerId)
	user := &chat.User{
		Id:           session.User.Id,
		Username:     session.User.Phone,
		LastChatTime: time.Now().Unix(),
		Disabled:     false,
		Online:       userM.IsOnline(session.CustomerId, session.UserId),
		LastMessage:  lastMsg,
		Unread:       unRead,
		Avatar:       "",
		Platform:     platform,
	}
	return user, nil

}

func (s sChat) Register(ctx context.Context, u any, conn *ghttp.WebSocket) error {
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
	case *relation.CustomerAdmins:
		uu, _ := u.(*relation.CustomerAdmins)
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

func (s sChat) IsOnline(customerId int, uid int, t string) bool {
	if t == "user" {
		return s.user.IsOnline(customerId, uid)
	}
	if t == "admin" {
		return s.admin.IsOnline(customerId, uid)
	}
	return false
}

func (s sChat) BroadcastWaitingUser(customerId int) {
	s.admin.broadcastWaitingUser(customerId)
}

func (s sChat) GetOnlineCount(customerId int) chat.OnlineCount {
	return chat.OnlineCount{
		Admin:   s.admin.GetOnlineTotal(customerId),
		User:    s.user.GetOnlineTotal(customerId),
		Waiting: service.ChatManual().GetTotalCount(customerId),
	}
}

func (s sChat) GetPlatform(customerId, uid int, t string) string {
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

func (s sChat) NoticeRate(msg *entity.CustomerChatMessages) {
	s.admin.noticeRate(msg)
}

func (s sChat) NoticeUserRead(customerId, uid int, msgIds []int64) {
	s.admin.NoticeRead(customerId, uid, msgIds)
}

func (s sChat) NoticeAdminRead(customerId, uid int, msgIds []int64) {
	s.user.NoticeRead(customerId, uid, msgIds)
}

func (s sChat) Transfer(fromAdmin *entity.CustomerAdmins, toId int, userId int, remark string) error {
	user := &entity.Users{}
	ctx := gctx.New()
	err := dao.Users.Ctx(ctx).
		Where("customer_id", fromAdmin.CustomerId).
		Where("id", userId).Scan(user)
	if err == sql.ErrNoRows {
		return gerror.NewCode(gcode.CodeNotFound)
	}
	toAdmin := &entity.CustomerAdmins{}
	err = dao.CustomerAdmins.Ctx(ctx).Where("customer_id", fromAdmin.CustomerId).
		Where("id", toId).
		Where("is_chat", 1).
		Scan(toAdmin)
	if err == sql.ErrNoRows {
		return gerror.NewCode(gcode.CodeNotFound)
	}
	isValid := service.ChatRelation().IsUserValid(gconv.Int(fromAdmin.Id), user.Id)
	if !isValid {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "用户已失效，无法转接")
	}
	return service.ChatTransfer().Create(gconv.Int(fromAdmin.Id), toId, userId, remark)
}

func (s sChat) GetOnlineAdmin(customerId int) []chat.SimpleUser {
	conns := s.admin.GetAllConn(customerId)
	res := make([]chat.SimpleUser, len(conns), len(conns))
	for index, c := range conns {
		res[index] = chat.SimpleUser{
			Id:       c.GetUserId(),
			Username: c.GetUser().GetUsername(),
		}
	}
	return res
}

func (s sChat) GetOnlineUser(customerId int) []chat.SimpleUser {
	conns := s.user.GetAllConn(customerId)
	res := make([]chat.SimpleUser, len(conns), len(conns))
	for index, c := range conns {
		res[index] = chat.SimpleUser{
			Id:       c.GetUserId(),
			Username: c.GetUser().GetUsername(),
		}
	}
	return res
}
