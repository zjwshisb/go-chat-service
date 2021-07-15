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
	var record = &models.QueryRecord{}
	databases.Db.Table("query_records").
		Where("user_id = ?" , uid).
		Where("service_id = ?", 0).Find(record)
	if record.Id == 0 {
		record.UserId = uid
		record.QueriedAt = time.Now().Unix()
		record.ServiceId = 0
		databases.Db.Save(record)
	}
}