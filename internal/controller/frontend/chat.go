package frontend

import (
	"context"
	baseApi "gf-chat/api/v1"
	api "gf-chat/api/v1/frontend"
	"gf-chat/internal/consts"
	"gf-chat/internal/library/storage"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

var CChat = &cChat{}

type cChat struct {
}

func (c cChat) Setting(ctx context.Context, _ *api.SettingReq) (res *baseApi.NormalRes[*api.SettingRes], err error) {
	customerId := service.UserCtx().GetCustomerId(ctx)
	isShowQueue, err := service.ChatSetting().GetIsUserShowQueue(ctx, customerId)
	if err != nil {
		return
	}
	isShowRead, err := service.ChatSetting().GetIsUserShowRead(ctx, customerId)
	if err != nil {
		return
	}
	return baseApi.NewResp(&api.SettingRes{
		IsShowQueue: isShowQueue,
		IsShowRead:  isShowRead,
	}), nil
}

func (c cChat) File(ctx context.Context, req *api.FileStoreReq) (res *baseApi.NormalRes[*baseApi.File], err error) {
	var parent *model.CustomerChatFile
	request := g.RequestFromCtx(ctx)
	dirVal := request.GetCtxVar("file-dir")
	if p, ok := dirVal.Interface().(*model.CustomerChatFile); ok {
		parent = p
	}
	file, err := req.File.Open()
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()
	fileType, _ := storage.FileType(file)
	relativePath := ""
	if parent != nil {
		relativePath = parent.Path
	}
	fileModel, err := storage.Disk().SaveUpload(ctx, req.File, relativePath)
	if err != nil {
		return nil, err
	}
	fileModel.Type = fileType
	fileModel.ParentId = 0
	fileModel.CustomerId = service.UserCtx().GetCustomerId(ctx)
	fileModel.FromId = service.UserCtx().GetId(ctx)
	fileModel.IsResource = 0
	fileModel.FromModel = "user"
	_, err = service.File().Insert(ctx, fileModel)
	if err != nil {
		return
	}
	return baseApi.NewResp(service.File().ToApi(fileModel)), nil
}

func (c cChat) Message(ctx context.Context, req *api.ChatMessageReq) (res *baseApi.NormalRes[[]*baseApi.ChatMessage], err error) {
	uid := service.UserCtx().GetUser(ctx).Id
	w := g.Map{
		"user_id": uid,
	}
	if req.Id > 0 {
		w["id <"] = req.Id
	}
	user := service.UserCtx().GetUser(ctx)
	messages, err := service.ChatMessage().All(ctx, w, nil, "id desc", req.PageSize)
	if err != nil {
		return
	}
	adminIds := slice.Unique(slice.Map(messages, func(index int, item *model.CustomerChatMessage) uint {
		return item.AdminId
	}))
	admins, err := service.Admin().GetAdminsWithSetting(ctx, do.CustomerAdmins{Id: adminIds})
	adminToMessageId := make(map[uint][]uint)
	r := make([]*baseApi.ChatMessage, 0, len(messages))
	for _, item := range messages {
		item.User = user
		if item.AdminId > 0 {
			item.Admin, _ = slice.FindBy(admins, func(index int, a *model.CustomerAdmin) bool {
				return a.Id == item.AdminId
			})
			if item.ReadAt == nil && item.Source == consts.MessageSourceAdmin {
				ids, exist := adminToMessageId[item.AdminId]
				if exist {
					adminToMessageId[item.AdminId] = append(ids, item.Id)
				} else {
					adminToMessageId[item.AdminId] = []uint{item.Id}
				}
			}
		}
		msg, err := service.ChatMessage().ToApi(ctx, item)
		if err != nil {
			return nil, err
		}
		r = append(r, msg)
	}
	customer := service.UserCtx().GetCustomerId(ctx)

	go func() {
		ctx := gctx.New()
		for adminId, ids := range adminToMessageId {
			_, err := service.ChatMessage().ToRead(ctx, ids)
			if err != nil {
				g.Log().Error(ctx, err)
			}
			if adminId > 0 {
				service.Chat().NoticeUserRead(customer, adminId, ids)
			}
		}

	}()
	return baseApi.NewResp(r), nil
}

func (c cChat) Read(ctx context.Context, req *api.ChatReadReq) (res *baseApi.NilRes, err error) {
	user := service.UserCtx().GetUser(ctx)
	message, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:     req.MsgId,
		UserId: user.Id,
		Source: []int{consts.MessageSourceAdmin, consts.MessageSourceSystem},
	})
	if err != nil {
		return
	}
	if message.ReadAt == nil {
		_, err = service.ChatMessage().ToRead(ctx, message.Id)
		msgIds := []uint{req.MsgId}
		if err != nil {
			return
		}
		service.Chat().NoticeUserRead(user.CustomerId, message.AdminId, msgIds)
	}
	return baseApi.NewNilResp(), nil
}

func (c cChat) Rate(ctx context.Context, req *api.ChatRateReq) (res *baseApi.NilRes, err error) {
	msg, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:     ghttp.RequestFromCtx(ctx).GetRouter("id"),
		UserId: service.UserCtx().GetUser(ctx).Id,
		Type:   consts.MessageTypeRate,
	})
	if err != nil {
		return
	}
	session, err := service.ChatSession().Find(ctx, msg.SessionId)
	if err != nil {
		return
	}
	msg.Content = gconv.String(req.Rate)
	_, err = service.ChatMessage().Save(ctx, msg)
	if err != nil {
		return
	}
	session.Rate = req.Rate
	_, err = service.ChatSession().Save(ctx, session)
	if err != nil {
		return
	}
	service.Chat().NoticeRate(msg)
	return baseApi.NewNilResp(), nil
}

func (c cChat) ReqId(_ context.Context, _ *api.ChatReqIdReq) (res *baseApi.NormalRes[api.ChatReqId], err error) {
	return baseApi.NewResp(api.ChatReqId{ReqId: service.ChatMessage().GenReqId()}), nil
}
