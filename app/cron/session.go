package cron

import (
	"time"
	"ws/app/chat"
	"ws/app/http/websocket"
	"ws/app/log"
	"ws/app/repositories"
)

func closeSessions() {
	log.Log.WithField("type", "cron").Info("<start-job:close-sessions>")
	admins := repositories.AdminRepo.Get([]*repositories.Where{}, -1, []string{}, []string{})
	for _, admin := range admins {
		uids, limits := chat.AdminService.GetUsersWithLimitTime(admin.GetPrimaryKey())
		length := len(uids)
		for i := 0; i <= length-1; i++ {
			uid := uids[i]
			limit := limits[i]
			if limit <= time.Now().Unix() {
				session := repositories.ChatSessionRepo.First([]*repositories.Where{
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
				}, []string{"id desc"})
				if session != nil {
					chat.SessionService.Close(session.Id, false, false)
					noticeMessage := admin.GetBreakMessage(uid, session.Id)
					websocket.UserManager.DeliveryMessage(noticeMessage, false)
					repositories.MessageRepo.Save(noticeMessage)
				}
			}
		}
	}
	log.Log.WithField("type", "cron").Info("<end-job:close-sessions>")
}
