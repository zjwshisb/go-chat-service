package backend

import (
	"context"
	v1 "gf-chat/api"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
)

var CChatFile = &cChatFile{}

type cChatFile struct {
}

func (f cChatFile) Index(ctx context.Context, req *api.ChatFileListReq) (res *v1.ListRes[api.ChatFile], err error) {
	paginator, err := service.File().Paginate(ctx, do.CustomerChatFiles{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		IsResource: 0,
	}, req.Paginate, nil, "id desc")
	if err != nil {
		return
	}
	adminIds := make([]uint, 0, len(paginator.Items))
	userIds := make([]uint, 0, len(paginator.Items))
	for _, f := range paginator.Items {
		if f.FromModel == "admin" {
			adminIds = append(adminIds, f.FromId)
		} else {
			userIds = append(userIds, f.FromId)
		}
	}
	admins, err := service.Admin().All(ctx, do.CustomerAdmins{
		Id: adminIds,
	}, nil, nil)
	if err != nil {
		return
	}
	adminsMap := slice.KeyBy(admins, func(item *model.CustomerAdmin) uint {
		return item.Id
	})
	users, err := service.User().All(ctx, do.Users{
		Id: userIds,
	}, nil, nil)
	usersMap := slice.KeyBy(users, func(item *model.User) uint {
		return item.Id
	})
	files := slice.Map(paginator.Items, func(index int, item *model.CustomerChatFile) api.ChatFile {
		var adminName string
		var userName string
		if item.FromModel == "admin" {
			admin, ok := adminsMap[item.FromId]
			if ok {
				adminName = admin.Username
			}
		} else {
			user, ok := usersMap[item.FromId]
			if ok {
				userName = user.Username
			}
		}
		return api.ChatFile{
			File:      service.File().ToApi(item),
			AdminName: adminName,
			CreatedAt: item.CreatedAt,
			UserName:  userName,
		}
	})
	res = v1.NewListResp(files, paginator.Total)
	return
}

func (f cChatFile) Delete(ctx context.Context, _ *api.ChatFileDeleteReq) (res *v1.NilRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	file, err := service.File().First(ctx, do.CustomerChatFiles{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if err != nil {
		return
	}
	err = service.File().RemoveFile(ctx, file)
	if err != nil {
		return
	}
	res = v1.NewNilResp()
	return
}
