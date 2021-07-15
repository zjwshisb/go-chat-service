package websocket

import "ws/internal/chat"

type userHub struct {
	BaseHub
}

func (userHub *userHub) addToManual(uid int64)  {
	_ = chat.AddToManual(uid)
	ServiceHub.BroadcastWaitingUser()
}