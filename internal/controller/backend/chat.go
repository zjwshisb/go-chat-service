package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
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

func (c cChat) CancelTransfer(ctx context.Context, _ *api.CancelTransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	transfer, err := service.ChatTransfer().First(ctx, g.Map{
		"to_admin_id":         admin.Id,
		"canceled_at is null": nil,
		"accepted_at is null": nil,
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

func (c cChat) RemoveUser(ctx context.Context, _ *api.RemoveUserReq) (res *baseApi.NilRes, err error) {
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

func (c cChat) RemoveInvalidUser(ctx context.Context, _ *api.RemoveAllUserReq) (res *baseApi.NormalRes[api.RemoveAllUserRes], err error) {
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

func (c cChat) User(ctx context.Context, _ *api.UserListReq) (res *baseApi.NormalRes[[]api.ChatUser], err error) {
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
	userMap := make(map[uint]*model.User)
	for _, user := range users {
		userMap[user.Id] = user
	}
	lastMessages, err := service.ChatMessage().GetLastGroupByUsers(ctx, admin.Id, ids)
	if err != nil {
		return
	}
	unreadCounts, err := service.ChatMessage().GetUnreadCountGroupByUsers(ctx, ids, do.CustomerChatMessages{
		AdminId: admin.Id,
		Source:  consts.MessageSourceUser,
	})
	if err != nil {
		return
	}
	items := make([]api.ChatUser, 0)
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
				Messages:    make([]api.ChatMessage, 0),
			}
			lastMsg, exist := slice.FindBy(lastMessages, func(index int, item *model.CustomerChatMessage) bool {
				return item.UserId == user.Id
			})
			if exist {
				lastMsg, err := service.ChatMessage().ToApi(ctx, lastMsg)
				if err != nil {
					return nil, err
				}
				cu.LastMessage = lastMsg
			}
			useUnread, exist := slice.FindBy(unreadCounts, func(index int, item model.UnreadCount) bool {
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
func (c cChat) UserInfo(ctx context.Context, _ *api.GetUserChatInfoReq) (res *baseApi.NormalRes[[]api.UserInfoItem], err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	admin := service.AdminCtx().GetUser(ctx)
	user, err := service.User().First(ctx, do.Users{
		Id:         id,
		CustomerId: admin.CustomerId,
	})
	if err != nil {
		return
	}
	info, err := service.User().GetInfo(ctx, user)
	if err != nil {
		return
	}
	return baseApi.NewResp(info), nil
}

func (c cChat) Message(ctx context.Context, req *api.GetMessageReq) (res *baseApi.NormalRes[[]*api.ChatMessage], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	settings, err := service.Admin().FindSetting(ctx, admin.Id, true)
	if err != nil {
		return
	}
	admin.Setting = settings
	r := make([]*api.ChatMessage, 0)
	w := g.Map{
		"user_id":  req.Uid,
		"admin_id": admin.Id,
		"source":   []uint{consts.MessageSourceUser, consts.MessageSourceAdmin},
	}
	if req.Mid > 0 {
		w["id < ?"] = req.Mid
	}
	messages, err := service.ChatMessage().All(ctx, w, g.Slice{
		model.CustomerChatMessage{}.User,
	}, "id desc")
	if err != nil {
		return
	}
	unReadIds := make([]uint, 0, len(messages))
	for _, i := range messages {
		i.Admin = admin
		if i.ReadAt == nil && i.Source == consts.MessageSourceUser {
			unReadIds = append(unReadIds, i.Id)
		}
		msg, err := service.ChatMessage().ToApi(ctx, i)
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
	return baseApi.NewResp(r), nil
}

func (c cChat) ReqId(_ context.Context, _ *api.ReqIdReq) (res *baseApi.NormalRes[api.ReqIdRes], err error) {
	res = baseApi.NewResp(api.ReqIdRes{ReqId: service.ChatMessage().GenReqId()})
	return
}

func (c cChat) Sessions(ctx context.Context, _ *api.GetUserSessionReq) (res *baseApi.NormalRes[api.UserSessionRes], err error) {
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
		r = append(r, service.ChatSession().ToApi(s))
	}
	res = baseApi.NewResp(r)
	return
}
