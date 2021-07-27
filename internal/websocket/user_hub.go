package websocket

import (
	"time"
	"ws/internal/action"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/models"
)

type userHub struct {
	BaseHub
}

func (userHub *userHub) Run() {
	userHub.Register(UserLogin, func(i ...interface{}) {
		length := len(i)
		if length >= 1 {
			ai := i[0]
			conn, ok := ai.(Conn)
			if ok {
				uid := conn.GetUserId()
				sid := chat.GetUserLastServerId(uid)
				if sid > 0 {
					adminConn, exist := AdminHub.GetConn(sid)
					if exist {
						adminConn.Deliver(action.NewUserOnline(uid))
					}
				}
			}
		}
	})
	userHub.Register(UserLogout, func(i ...interface{}) {
		length := len(i)
		if length >= 1 {
			ai := i[0]
			conn, ok := ai.(Conn)
			if ok {
				uid := conn.GetUserId()
				sid := chat.GetUserLastServerId(uid)
				if sid > 0 {
					adminConn, exist := AdminHub.GetConn(sid)
					if exist {
						adminConn.Deliver(action.NewUserOffline(uid))
					}
				}
			}
		}
	})
	userHub.BaseHub.Run()
}
func (userHub *userHub) addToManual(uid int64)  {
	_ = chat.AddToManual(uid)
	AdminHub.BroadcastWaitingUser()
	session := chat.GetSession(uid, 0)
	if session == nil {
		session = &models.ChatSession{}
		session.UserId = uid
		session.QueriedAt = time.Now().Unix()
		session.AdminId = 0
		databases.Db.Save(session)
	}
}