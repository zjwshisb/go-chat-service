package chat

import (
	"time"
	"ws/app/models"
	"ws/app/repositories"
)
var SessionService = &sessionService{}
type sessionService struct {

}
// 关闭会话
func (sessionService *sessionService) Close(sessionId uint64, isRemoveUser bool, updateTime bool) {
	session := chatSessionRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: sessionId,
		},
	})
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
func (sessionService *sessionService) Create(uid int64, ty int) *models.ChatSession  {
	session := &models.ChatSession{}
	session.UserId = uid
	session.QueriedAt = time.Now().Unix()
	session.AdminId = 0
	session.Type = ty
	_ = chatSessionRepo.Save(session)
	return session
}
func (sessionService *sessionService) Get(uid int64, adminId int64) *models.ChatSession {
	session := chatSessionRepo.First([]*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: uid,
		},
		{
			Filed: "admin_id = ?",
			Value: adminId,
		},
		{
			Filed: "broke_at = ? ",
			Value: 0,
		},
		{
			Filed: "canceled_at = ?",
			Value: 0,
		},
	}, "id desc")
	return session
}

