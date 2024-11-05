package backend

import (
	"context"
	"database/sql"
	baseApi "gf-chat/api"
	"gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

var CChat = &cChat{}

type cChat struct {
}

func (c cChat) SmsNotice(ctx context.Context, req *backend.ChatSmsNoticeReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	code := service.ChatSetting().GetSmsCode(admin.CustomerId)
	if code == "" {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "没有配置短信模板")
	}
	if !service.ChatRelation().IsUserValid(gconv.Int(admin.Id), req.Uid) {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "用户已失效")
	}
	err = service.Sms().Send(ctx, code, req.Uid)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, err.Error())
	}
	return &baseApi.NilRes{}, err
}

func (c cChat) AcceptUser(ctx context.Context, req *backend.ChatAcceptReq) (res *backend.ChatAcceptRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	user, err := service.Chat().Accept(*admin, req.Sid)
	if err != nil {
		return nil, err
	}
	service.Chat().BroadcastWaitingUser(admin.CustomerId)
	return &backend.ChatAcceptRes{
		User: *user,
	}, nil
}

func (c cChat) Read(ctx context.Context, req *backend.ChatReadReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	message := service.ChatMessage().First(do.CustomerChatMessages{
		Id:      req.MsgId,
		AdminId: admin.Id,
		Source:  consts.MessageSourceUser,
	})
	if message == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	if message.SendAt == 0 {
		service.ChatMessage().ChangeToRead([]int64{message.Id})
		service.Chat().NoticeAdminRead(admin.CustomerId, message.UserId, []int64{message.Id})
	}
	return &baseApi.NilRes{}, nil
}

func (c cChat) CancelTransfer(ctx context.Context, req *backend.ChatCancelTransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	transfer := service.ChatTransfer().FirstRelation(do.CustomerChatTransfers{
		ToAdminId:  admin.Id,
		CanceledAt: 0,
		AcceptedAt: 0,
	})
	if transfer != nil {
		_ = service.ChatTransfer().Cancel(transfer)
	}
	return &baseApi.NilRes{}, nil
}

func (c cChat) Transfer(ctx context.Context, req *backend.ChatTransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	err = service.Chat().Transfer(admin, req.ToId, req.UserId, req.Remark)
	if err != nil {
		return nil, err
	}
	return &baseApi.NilRes{}, err
}

func (c cChat) RemoveUser(ctx context.Context, req *backend.ChatRemoveReq) (res *baseApi.NilRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	admin := service.AdminCtx().GetAdmin(ctx)
	session := service.ChatSession().ActiveOne(gconv.Int(id), admin.Id, nil)
	if session != nil {
		if session.BrokeAt == 0 {
			service.ChatSession().Close(session, true, false)
		}
	} else {
		service.ChatRelation().RemoveUser(gconv.Int(admin.Id), gconv.Int(id))
	}
	return &baseApi.NilRes{}, err
}

func (c cChat) RemoveInvalidUser(ctx context.Context, req *backend.ChatRemoveAllReq) (res *backend.ChatRemoveAllRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	ids := service.ChatRelation().GetInvalidUsers(gconv.Int(admin.Id))
	resIds := make([]int, 0, len(ids))
	for _, id := range ids {
		session := &entity.CustomerChatSessions{}
		err := dao.CustomerChatSessions.Ctx(ctx).
			Where("admin_id", admin.Id).
			Where("user_id", id).
			OrderDesc("id").Scan(session)
		if err == sql.ErrNoRows {
			continue
		}
		service.ChatSession().Close(session, true, false)
		resIds = append(resIds, id)
	}
	return &backend.ChatRemoveAllRes{Ids: resIds}, nil
}

func (c cChat) User(ctx context.Context, req *backend.ChatUserListReq) (res *backend.ChatUserListRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	ids, times := service.ChatRelation().GetUsersWithLimitTime(gconv.Int(admin.Id))
	users := service.User().GetUsers(ctx, map[string]any{
		"id": ids,
	})
	resp := backend.ChatUserListRes{}
	count := 0
	userMap := make(map[int]*entity.Users)
	for _, user := range users {
		userMap[user.Id] = user
	}
	type maxChatIds struct {
		Id     int
		UserId int
	}
	lastRes := make([]maxChatIds, 0)
	sources := []int{consts.MessageSourceAdmin, consts.MessageSourceUser}
	err = dao.CustomerChatMessages.Ctx(ctx).
		FieldMax("id", "id").
		Fields("user_id").
		Where("user_id in (?)", ids).
		Where("source in (?)", sources).
		Where("admin_id", admin.Id).
		Group("user_id").Scan(&lastRes)
	lastMessages := make([]relation.CustomerChatMessages, 0)
	lastMessageIds := slice.Map(lastRes, func(index int, item maxChatIds) int {
		return item.Id
	})
	_ = dao.CustomerChatMessages.Ctx(ctx).Where("id in (?)", lastMessageIds).
		Where("source in (?)", sources).
		WithAll().
		Scan(&lastMessages)
	type Unread struct {
		Count  int
		UserId int
	}
	unreadCounts := make([]Unread, 0)
	_ = dao.CustomerChatMessages.Ctx(ctx).Where("user_id in (?)", ids).
		FieldCount("*", "Count").
		Fields("user_id").
		Group("user_id").
		Where("admin_id", admin.Id).
		Where("read_at", 0).
		Where("source", consts.MessageSourceUser).
		Scan(&unreadCounts)
	for index, id := range ids {
		limitTime := times[index]
		disabled := limitTime <= time.Now().Unix()
		if count > 50 && disabled {
			go func() {
				_ = service.ChatRelation().RemoveUser(gconv.Int(admin.Id), id)
			}()
			continue
		}
		count = count + 1
		user, exist := userMap[id]
		if exist {
			cu := chat.User{
				Id:          user.Id,
				Username:    user.Phone,
				Disabled:    disabled,
				LastMessage: nil,
				Unread:      0,
				Avatar:      "",
				Online:      service.Chat().IsOnline(user.CustomerId, user.Id, "user"),
				Platform:    service.Chat().GetPlatform(user.CustomerId, user.Id, "user"),
			}
			lastMsg, exist := slice.Find(lastMessages, func(index int, item relation.CustomerChatMessages) bool {
				return item.UserId == user.Id
			})
			if exist {
				lastMsg := service.ChatMessage().RelationToChat(*lastMsg)
				cu.LastMessage = &lastMsg
			}
			useUnread, exist := slice.Find(unreadCounts, func(index int, item Unread) bool {
				return item.UserId == cu.Id
			})
			if exist {
				cu.Unread = useUnread.Count
			}
			cu.LastChatTime = service.ChatRelation().GetLastChatTime(gconv.Int(admin.Id), cu.Id)
			resp = append(resp, cu)
		}

	}
	return &resp, nil
}

func (c cChat) UserInfo(ctx context.Context, req *backend.ChatUserInfoReq) (res *backend.ChatUserInfoRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	admin := service.AdminCtx().GetAdmin(ctx)
	user := service.User().First(do.Users{
		Id:         id,
		CustomerId: admin.CustomerId,
	})
	if user == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}

	return &backend.ChatUserInfoRes{
		Phone: user.Phone,
	}, nil
}

func (c cChat) Message(ctx context.Context, req *backend.ChatMessageReq) (res *backend.ChatMessageRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	r := backend.ChatMessageRes{}
	models := service.ChatMessage().GetModels(req.Mid, g.Map{
		"user_id":  req.Uid,
		"admin_id": admin.Id,
		"source":   []int{consts.MessageSourceUser, consts.MessageSourceAdmin},
	}, 20)
	unReadIds := make([]int64, 0, len(models))
	for _, i := range models {
		if i.ReadAt == 0 && i.Source == consts.MessageSourceUser {
			unReadIds = append(unReadIds, i.Id)
		}
		r = append(r, service.ChatMessage().RelationToChat(*i))
	}
	if len(unReadIds) > 0 {
		_, _ = service.ChatMessage().ChangeToRead(unReadIds)
		service.Chat().NoticeAdminRead(admin.CustomerId, req.Uid, unReadIds)
	}
	return &r, nil
}

func (c cChat) ReqId(ctx context.Context, req *backend.ChatReqIdReq) (res *backend.ChatReqIdRes, err error) {
	return &backend.ChatReqIdRes{ReqId: service.ChatMessage().GenReqId()}, nil
}

func (c cChat) Sessions(ctx context.Context, req *backend.ChatUserSessionReq) (res *backend.ChatUserSessionRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	sessions := service.ChatSession().Get(ctx, map[string]any{
		"user_id":     id,
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
		"admin_id>":   0,
	})
	r := backend.ChatUserSessionRes{}
	for _, s := range sessions {
		r = append(r, service.ChatSession().RelationToChat(s))
	}
	return &r, nil
}
