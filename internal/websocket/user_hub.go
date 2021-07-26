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
	var session = chat.GetSession(uid, 0)
	if session == nil {
		session = &models.ChatSession{}
		session.UserId = uid
		session.QueriedAt = time.Now().Unix()
		session.ServiceId = 0
		databases.Db.Save(session)
	}
}