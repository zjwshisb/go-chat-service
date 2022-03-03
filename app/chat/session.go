package chat

import (
	"time"
	"ws/app/repositories"
)

var SessionService = &sessionService{}
type sessionService struct {

}

// Close 关闭会话
func (sessionService *sessionService) Close(sessionId uint64, isRemoveUser bool, updateTime bool) {
	session := chatSessionRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: sessionId,
		},
	}, []string{})
	if session != nil {
		session.BrokeAt = time.Now().Unix()
		chatSessionRepo.Save(session)
		if isRemoveUser {
			_ = AdminService.RemoveUser(session.AdminId, session.UserId)
		}
		if updateTime {
			_ = AdminService.UpdateLimitTime(session.AdminId, session.UserId, 0)
		}
	}
}




