package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

var CChat = &cChat{}

type cChat struct {
}

func (c cChat) AcceptUser(ctx context.Context, req *api.AcceptUserReq) (res *baseApi.NormalRes[api.AcceptRes], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	user, err := service.Chat().Accept(ctx, *admin, req.Sid)
	if err != nil {
		return nil, err
	}
	err = service.Chat().BroadcastWaitingUser(ctx, admin.CustomerId)
	if err != nil {
		return
	}
	res = baseApi.NewResp(api.AcceptRes{
		User: *user,
	})
	return
}

func (c cChat) Read(ctx context.Context, req *api.MessageReadReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	message, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:      req.MsgId,
		AdminId: admin.Id,
		Source:  consts.MessageSourceUser,
	})
	if err != nil {
		return
	}
	if message.SendAt == nil {
		_, err = service.ChatMessage().ToRead(ctx, message.Id)
		if err != nil {
			return
		}
		service.Chat().NoticeAdminRead(admin.CustomerId, message.UserId, []uint{message.Id})
	}
	return baseApi.NewNilResp(), nil
}

func (c cChat) CancelTransfer(ctx context.Context, req *api.CancelTransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	transfer, err := service.ChatTransfer().First(ctx, do.CustomerChatTransfers{
		ToAdminId:  admin.Id,
		CanceledAt: nil,
		AcceptedAt: nil,
	})
	if err != nil {
		return
	}
	err = service.ChatTransfer().Cancel(ctx, transfer)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cChat) Transfer(ctx context.Context, req *api.StoreTransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	err = service.Chat().Transfer(ctx, admin, req.ToId, req.UserId, req.Remark)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cChat) RemoveUser(ctx context.Context, req *api.RemoveUserReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")

	session, err := service.ChatSession().FirstActive(ctx, gconv.Uint(id), admin.Id, nil)
	if err != nil {
		return
	}
	err = service.ChatRelation().RemoveUser(ctx, admin.Id, gconv.Uint(id))
	if err != nil {
		return
	}
	if session.BrokenAt == nil {
		err = service.ChatSession().Close(ctx, session, true, false)
		if err != nil {
			return
		}
	}
	return baseApi.NewNilResp(), nil
}

func (c cChat) RemoveInvalidUser(ctx context.Context, req *api.RemoveAllUserReq) (res *baseApi.NormalRes[api.RemoveAllUserRes], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	ids := service.ChatRelation().GetInvalidUsers(ctx, admin.Id)
	resIds := make([]uint, 0, len(ids))
	for _, id := range ids {
		session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{
			AdminId: admin.Id,
			UserId:  id,
		})
		if err != nil {
			continue
		}
		err = service.ChatSession().Close(ctx, session, true, false)
		if err == nil {
			resIds = append(resIds, id)
		}
	}
	res = baseApi.NewResp(api.RemoveAllUserRes{Ids: resIds})
	return
}

func (c cChat) User(ctx context.Context, req *api.UserListReq) (res *baseApi.NormalRes[api.UserListRes], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	ids, times := service.ChatRelation().GetUsersWithLimitTime(ctx, gconv.Uint(admin.Id))
	users, err := service.User().All(ctx, do.Users{
		Id:         ids,
		CustomerId: admin.CustomerId,
	}, nil, nil)
	if err != nil {
		return
	}
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
	items := api.UserListRes{}
	for index, id := range ids {
		limitTime := times[index]
		disabled := limitTime <= time.Now().Unix()
		if count > 50 && disabled {
			go func() {
				_ = service.ChatRelation().RemoveUser(ctx, admin.Id, id)
			}()
			continue
		}
		count = count + 1
		user, exist := userMap[id]
		if exist {
			cu := api.ChatUser{
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
				lastMsg, err := service.ChatMessage().ToApi(ctx, *lastMsg)
				if err != nil {
					return nil, err
				}
				cu.LastMessage = &lastMsg
			}
			useUnread, exist := slice.Find(unreadCounts, func(index int, item Unread) bool {
				return item.UserId == cu.Id
			})
			if exist {
				cu.Unread = useUnread.Count
			}
			cu.LastChatTime = gtime.NewFromTimeStamp(int64(service.ChatRelation().GetLastChatTime(ctx, admin.Id, cu.Id)))
			items = append(items, cu)
		}
	}
	res = baseApi.NewResp(items)
	return
}
func (c cChat) UserInfo(ctx context.Context, req *api.GetUserChatInfoReq) (res *baseApi.NormalRes[api.UserInfoRes], err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	admin := service.AdminCtx().GetUser(ctx)
	user, err := service.User().First(ctx, do.Users{
		Id:         id,
		CustomerId: admin.CustomerId,
	})
	if err != nil {
		return
	}
	res = baseApi.NewResp(api.UserInfoRes{
		Phone: user.Username,
	})
	return
}

func (c cChat) Message(ctx context.Context, req *api.GetMessageReq) (res *baseApi.NormalRes[api.MessageRes], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	r := api.MessageRes{}
	models, err := service.ChatMessage().GetList(ctx, req.Mid, g.Map{
		"user_id":  req.Uid,
		"admin_id": admin.Id,
		"source":   []uint{consts.MessageSourceUser, consts.MessageSourceAdmin},
	}, 20)
	if err != nil {
		return
	}
	unReadIds := make([]uint, 0, len(models))
	for _, i := range models {
		if i.ReadAt == nil && i.Source == consts.MessageSourceUser {
			unReadIds = append(unReadIds, i.Id)
		}
		msg, err := service.ChatMessage().ToApi(ctx, *i)
		if err != nil {
			return nil, err
		}
		r = append(r, msg)
	}
	if len(unReadIds) > 0 {
		_, err = service.ChatMessage().ToRead(ctx, unReadIds)
		if err != nil {
			return
		}
		service.Chat().NoticeAdminRead(admin.CustomerId, req.Uid, unReadIds)
	}
	res = baseApi.NewResp(r)
	return
}

func (c cChat) ReqId(ctx context.Context, req *api.ReqIdReq) (res *baseApi.NormalRes[api.ReqIdRes], err error) {
	res = baseApi.NewResp(api.ReqIdRes{ReqId: service.ChatMessage().GenReqId()})
	return
}

func (c cChat) Sessions(ctx context.Context, req *api.GetUserSessionReq) (res *baseApi.NormalRes[api.UserSessionRes], err error) {
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
	r := api.UserSessionRes{}
	for _, s := range sessions {
		r = append(r, service.ChatSession().RelationToChat(s))
	}
	res = baseApi.NewResp(r)
	return
}
