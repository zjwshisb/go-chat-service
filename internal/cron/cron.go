package cron

func Run() {
	//ctx := gctx.New()
	//gtimer.Add(ctx, time.Minute, func(ctx context.Context) {
	//	customers := make([]entity.Customers, 0)
	//	dao.Customers.Ctx(ctx).Scan(&customers)
	//	for _, c := range customers {
	//		admins := service.Admin().GetChatAll(c.Id)
	//		for _, admin := range admins {
	//			invalidIds := service.ChatRelation().GetInvalidUsers(gconv.Int(admin.Id))
	//			if len(invalidIds) > 0 {
	//				dao.CustomerChatSessions.Ctx(ctx).Where(do.CustomerChatSessions{
	//					UserId:     invalidIds,
	//					BrokeAt:    0,
	//					AdminId:    admin.Id,
	//					CanceledAt: 0,
	//				}).Where("accepted_at > ?", 0).Data("broke_at", time.Now().Unix()).Update()
	//			}
	//		}
	//	}
	//})
}
