package cron

import (
	"time"
	"ws/app/chat"
	"ws/app/log"
	"ws/app/repositories"
	"ws/app/websocket"
)

func closeSessions()  {
	admins := adminRepo.Get([]*repositories.Where{},-1, []string{})
	for _, admin := range admins {
		uids, limits := chat.AdminService.GetUsersWithLimitTime(admin.GetPrimaryKey())
		length := len(uids)
		for i := 0; i <= length-1; i++ {
			uid := uids[i]
			limit := limits[i]
			if limit <= time.Now().Unix() {
				session := sessionRepo.First([]*repositories.Where{
					{
						Filed: "admin_id = ?",
						Value: admin.GetPrimaryKey(),
					},
					{
						Filed: "user_id = ?",
						Value: uid,
					},
					{
						Filed: "broke_at = ?",
						Value: 0,
					},
				}, "id desc")
				if session != nil {
					chat.SessionService.Close(session.Id, false, false)
					noticeMessage := admin.GetBreakMessage(uid, session.Id)
					websocket.UserManager.DeliveryMessage(noticeMessage, false)
					messageRepo.Save(noticeMessage)
				}
			}
		}
	}
	log.Log.Info("close sessions done")
}
