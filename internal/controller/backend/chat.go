package backend

import (
	"context"
	baseApi "gf-chat/api"
	"gf-chat/api/backend/v1"
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

func (c cChat) AcceptUser(ctx context.Context, req *v1.AcceptUserReq) (res *baseApi.NormalRes[v1.ChatUser], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	user, err := service.Chat().Accept(ctx, *admin, req.Sid)
	if err != nil {
		return nil, err
	}
	err = service.Chat().BroadcastWaitingUser(ctx, admin.CustomerId)
	if err != nil {
		return
	}
	res = baseApi.NewResp(
		*user)
	return
}

func (c cChat) Read(ctx context.Context, req *v1.MessageReadReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	message, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:      req.MsgId,
		AdminId: admin.Id,
		Source:  consts.MessageSourceUser,
	})
	if err != nil {
		return
	}
	if message.ReadAt == nil {
		_, err = service.ChatMessage().ToRead(ctx, message.Id)
		if err != nil {
			return
		}
		err = service.Chat().NoticeRead(ctx, admin.CustomerId, message.UserId, []uint{message.Id}, consts.WsTypeUser)
		if err != nil {
			return
		}
	}
	return baseApi.NewNilResp(), nil
}

func (c cChat) CancelTransfer(ctx context.Context, _ *v1.CancelTransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	transfer, err := service.ChatTransfer().First(ctx, g.Map{
		"to_admin_id": admin.Id,
		"id":          ghttp.RequestFromCtx(ctx).GetRouter("id").Int(),
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

func (c cChat) Transfer(ctx context.Context, req *v1.StoreTransferReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	err = service.ChatTransfer().Create(ctx, admin, req.ToId, req.UserId, req.Remark)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cChat) RemoveUser(ctx context.Context, _ *v1.RemoveUserReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{
		UserId:  id,
		AdminId: admin.Id,
	}, "id desc")
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

func (c cChat) RemoveInvalidUser(ctx context.Context, _ *v1.RemoveAllUserReq) (res *baseApi.NormalRes[v1.RemoveAllUserRes], err error) {
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
	res = baseApi.NewResp(v1.RemoveAllUserRes{Ids: resIds})
	return
}

func (c cChat) User(ctx context.Context, _ *v1.UserListReq) (res *baseApi.NormalRes[[]v1.ChatUser], err error) {
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
	items := make([]v1.ChatUser, 0)
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
			online, platform := service.Chat().GetConnInfo(ctx, user.CustomerId, user.Id, consts.WsTypeUser)
			cu := v1.ChatUser{
				Id:          user.Id,
				Username:    user.Username,
				Disabled:    disabled,
				LastMessage: nil,
				Unread:      0,
				Avatar:      "",
				Online:      online,
				Platform:    platform,
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
func (c cChat) UserInfo(ctx context.Context, _ *v1.GetUserChatInfoReq) (res *baseApi.NormalRes[[]v1.UserInfoItem], err error) {
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

func (c cChat) TransferMessage(ctx context.Context, _ *v1.TransferMessageReq) (res *baseApi.NormalRes[[]*baseApi.ChatMessage], err error) {
	transfer, err := service.ChatTransfer().First(ctx, do.CustomerChatTransfers{
		Id:        ghttp.RequestFromCtx(ctx).GetRouter("id"),
		ToAdminId: service.AdminCtx().GetId(ctx),
	}, nil, nil)
	if err != nil {
		return
	}
	messages, err := service.ChatMessage().All(ctx, do.CustomerChatMessages{
		SessionId: transfer.FromSessionId,
		Source:    g.Slice{consts.MessageSourceUser, consts.MessageSourceAdmin},
	}, g.Slice{
		model.CustomerChatMessage{}.User,
		model.CustomerChatMessage{}.Admin,
	}, "id desc")
	if err != nil {
		return
	}
	apiMessage := slice.Map(messages, func(index int, item *model.CustomerChatMessage) *baseApi.ChatMessage {
		i, _ := service.ChatMessage().ToApi(ctx, item)
		return i
	})
	return baseApi.NewResp(apiMessage), nil
}

func (c cChat) Message(ctx context.Context, req *v1.GetMessageReq) (res *baseApi.NormalRes[[]*baseApi.ChatMessage], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	settings, err := service.Admin().FindSetting(ctx, admin.Id, true)
	if err != nil {
		return
	}
	admin.Setting = settings
	r := make([]*baseApi.ChatMessage, 0)
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
	}, "id desc", 50)
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
		err = service.Chat().NoticeRead(ctx, admin.CustomerId, req.Uid, unReadIds, consts.WsTypeUser)
		if err != nil {
			return
		}
	}
	return baseApi.NewResp(r), nil
}

func (c cChat) ReqId(_ context.Context, _ *v1.ReqIdReq) (res *baseApi.NormalRes[v1.ReqIdRes], err error) {
	res = baseApi.NewResp(v1.ReqIdRes{ReqId: service.ChatMessage().GenReqId()})
	return
}

func (c cChat) Sessions(ctx context.Context, _ *v1.GetUserSessionReq) (res *baseApi.NormalRes[[]v1.ChatSession], err error) {
	sessions, err := service.ChatSession().All(ctx, map[string]any{
		"user_id":     ghttp.RequestFromCtx(ctx).GetRouter("id"),
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
		"admin_id >":  0,
	}, g.Array{
		model.CustomerChatSession{}.User,
		model.CustomerChatSession{}.Admin,
	}, "id desc")
	if err != nil {
		return
	}
	apiSessions := slice.Map(sessions, func(index int, item *model.CustomerChatSession) v1.ChatSession {
		return service.ChatSession().ToApi(item)
	})
	return baseApi.NewResp(apiSessions), nil
}
