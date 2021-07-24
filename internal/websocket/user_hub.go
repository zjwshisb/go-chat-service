package websocket

import (
	"time"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/models"
)

type userHub struct {
	BaseHub
}

func (userHub *userHub) addToManual(uid int64)  {
	_ = chat.AddToManual(uid)
	ServiceHub.BroadcastWaitingUser()
	var session = &models.ChatSession{}
	databases.Db.Table("query_records").
		Where("user_id = ?" , uid).
		Where("service_id = ?", 0).
		Find(session)
	if session.Id == 0 {
		session.UserId = uid
		session.QueriedAt = time.Now().Unix()
		session.ServiceId = 0
		databases.Db.Save(session)
	}
}