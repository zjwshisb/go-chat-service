package backend

import (
	"context"
	"database/sql"
	baseApi "gf-chat/api"
	chatapi "gf-chat/api/v1/backend/chat"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

var CChat = &cChat{}

type cChat struct {
}

func (c cChat) AcceptUser(ctx context.Context, req *chatapi.AcceptReq) (res *chatapi.AcceptRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	user, err := service.Chat().Accept(ctx, *admin, req.Sid)
	if err != nil {
		return nil, err
	}
	service.Chat().BroadcastWaitingUser(admin.CustomerId)
	return &chatapi.AcceptRes{
		User: *user,
	}, nil
}

func (c cChat) Read(ctx context.Context, req *chatapi.ReadReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	message, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:      req.MsgId,
		AdminId: admin.Id,
		Source:  consts.MessageSourceUser,
	})
	if err != nil {
		return
	}
	if message.SendAt == nil {
		service.ChatMessage().ChangeToRead([]uint{message.Id})
		service.Chat().NoticeAdminRead(admin.CustomerId, message.UserId, []uint{message.Id})
	}
	return &baseApi.NilRes{}, nil
}

func (c cChat) CancelTransfer(ctx context.Context, req *chatapi.CancelTransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	transfer, err := service.ChatTransfer().First(do.CustomerChatTransfers{
		ToAdminId:  admin.Id,
		CanceledAt: nil,
		AcceptedAt: nil,
	})
	if err != nil {
		return
	}
	if transfer != nil {
		_ = service.ChatTransfer().Cancel(transfer)
	}
	return &baseApi.NilRes{}, nil
}

func (c cChat) Transfer(ctx context.Context, req *chatapi.TransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	err = service.Chat().Transfer(admin, req.ToId, req.UserId, req.Remark)
	if err != nil {
		return nil, err
	}
	return &baseApi.NilRes{}, err
}

func (c cChat) RemoveUser(ctx context.Context, req *chatapi.RemoveReq) (res *baseApi.NilRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	admin := service.AdminCtx().GetAdmin(ctx)
	session, err := service.ChatSession().ActiveOne(ctx, gconv.Uint(id), admin.Id, nil)
	if err != nil {
		if err != sql.ErrNoRows {
			return

		}
		service.ChatRelation().RemoveUser(admin.Id, gconv.Uint(id))

	}
	if session.BrokenAt == nil {
		service.ChatSession().Close(ctx, session, true, false)
	}
	return &baseApi.NilRes{}, err
}

func (c cChat) RemoveInvalidUser(ctx context.Context, req *chatapi.RemoveAllReq) (res *chatapi.RemoveAllRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	ids := service.ChatRelation().GetInvalidUsers(admin.Id)
	resIds := make([]uint, 0, len(ids))
	for _, id := range ids {
		session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{
			AdminId: admin.Id,
			UserId:  id,
		})
		if err != nil {
			continue
		}
		service.ChatSession().Close(ctx, session, true, false)
		resIds = append(resIds, id)
	}
	return &chatapi.RemoveAllRes{Ids: resIds}, nil
}

func (c cChat) User(ctx context.Context, req *chatapi.UserListReq) (res *chatapi.UserListRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	ids, times := service.ChatRelation().GetUsersWithLimitTime(gconv.Uint(admin.Id))
	users := service.User().GetUsers(ctx, map[string]any{
		"id": ids,
	})
	resp := chatapi.UserListRes{}
	count := 0
	userMap := make(map[uint]*entity.Users)
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
		FieldMax("id").
		Fields("user_id").
		Where("user_id in (?)", ids).
		Where("source in (?)", sources).
		Where("admin_id", admin.Id).
		Group("user_id").Scan(&lastRes)
	if err != nil {
		return
	}
	lastMessages := make([]model.CustomerChatMessage, 0)
	lastMessageIds := slice.Map(lastRes, func(index int, item maxChatIds) int {
		return item.Id
	})
	_ = dao.CustomerChatMessages.Ctx(ctx).Where("id in (?)", lastMessageIds).
		Where("source in (?)", sources).
		WithAll().
		Scan(&lastMessages)
	type Unread struct {
		Count  uint
		UserId uint
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
				_ = service.ChatRelation().RemoveUser(admin.Id, id)
			}()
			continue
		}
		count = count + 1
		user, exist := userMap[id]
		if exist {
			cu := model.ChatUser{
				Id:          user.Id,
				Username:    user.Username,
				Disabled:    disabled,
				LastMessage: nil,
				Unread:      0,
				Avatar:      "",
				Online:      service.Chat().IsOnline(user.CustomerId, user.Id, "user"),
				Platform:    service.Chat().GetPlatform(user.CustomerId, user.Id, "user"),
			}
			lastMsg, exist := slice.Find(lastMessages, func(index int, item model.CustomerChatMessage) bool {
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
			cu.LastChatTime = gtime.NewFromTimeStamp(int64(service.ChatRelation().GetLastChatTime(admin.Id, cu.Id)))
			resp = append(resp, cu)
		}

	}
	return &resp, nil
}

func (c cChat) UserInfo(ctx context.Context, req *chatapi.UserInfoReq) (res *chatapi.UserInfoRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	admin := service.AdminCtx().GetAdmin(ctx)
	user := service.User().First(do.Users{
		Id:         id,
		CustomerId: admin.CustomerId,
	})
	if user == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	return &chatapi.UserInfoRes{
		Phone: user.Username,
	}, nil
}

func (c cChat) Message(ctx context.Context, req *chatapi.MessageReq) (res *chatapi.MessageRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	r := chatapi.MessageRes{}
	models := service.ChatMessage().GetModels(req.Mid, g.Map{
		"user_id":  req.Uid,
		"admin_id": admin.Id,
		"source":   []uint{consts.MessageSourceUser, consts.MessageSourceAdmin},
	}, 20)
	unReadIds := make([]uint, 0, len(models))
	for _, i := range models {
		if i.ReadAt == nil && i.Source == consts.MessageSourceUser {
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

func (c cChat) ReqId(ctx context.Context, req *chatapi.ReqIdReq) (res *chatapi.ReqIdRes, err error) {
	return &chatapi.ReqIdRes{ReqId: service.ChatMessage().GenReqId()}, nil
}

func (c cChat) Sessions(ctx context.Context, req *chatapi.UserSessionReq) (res *chatapi.UserSessionRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	sessions, err := service.ChatSession().All(ctx, map[string]any{
		"user_id":     id,
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
		"admin_id>":   0,
	}, g.Array{
		model.CustomerChatSession{}.User,
		model.CustomerChatSession{}.Admin,
	}, nil)
	if err != nil {
		return
	}
	r := chatapi.UserSessionRes{}
	for _, s := range sessions {
		r = append(r, service.ChatSession().RelationToChat(s))
	}
	return &r, nil
}
