package cron

import (
	"context"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtimer"
	"time"
)

func Run() {
	ctx := gctx.New()
	gtimer.Add(ctx, time.Minute, func(ctx context.Context) {
		admins, err := service.Admin().All(ctx, do.CustomerAdmins{}, nil, "id desc")
		if err != nil {
			g.Log().Errorf(ctx, "%+v", err)
		}
		for _, admin := range admins {
			invalidIds := service.ChatRelation().GetInvalidUsers(ctx, admin.Id)
			if len(invalidIds) > 0 {
				dao.CustomerChatSessions.Ctx(ctx).Where(g.Map{
					"user_id":  invalidIds,
					"admin_id": admin.Id,
				})
			}
		}
	})
}
