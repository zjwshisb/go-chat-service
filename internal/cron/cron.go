package cron

import (
	"context"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtimer"
)

func Run() {
	ctx := gctx.New()
	gtimer.Add(ctx, time.Minute, func(ctx context.Context) {
		admins, err := service.Admin().All(ctx, do.CustomerAdmins{}, nil, "id desc")
		if err != nil {
			g.Log().Errorf(ctx, "%+v", err)
		}
		for _, admin := range admins {
			invalidIds, err := service.Chat().GetInvalidUsers(ctx, admin.Id)
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
				return
			}
			if len(invalidIds) > 0 {
				sessions, err := service.ChatSession().All(ctx, g.Map{
					"user_id":           invalidIds,
					"admin_id":          admin.Id,
					"broken_at is null": nil,
				}, nil, nil)
				if err != nil {
					g.Log().Errorf(ctx, "%+v", err)
				}
				for _, session := range sessions {
					if session.BrokenAt == nil {
						err = service.ChatSession().Close(ctx, session, false, false)
						if err != nil {
							g.Log().Errorf(ctx, "%+v", err)
						}
					}
				}
			}
		}
	})
}
