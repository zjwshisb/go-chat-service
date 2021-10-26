package chat

import (
	"log"
	"strconv"
	"time"
	"ws/app/models"
	"ws/app/repositories"
)
var SessionService = &sessionService{}
type sessionService struct {

}
// 关闭会话
func (sessionService *sessionService) Close(session *models.ChatSession, isRemoveUser bool, updateTime bool) {
	session.BrokeAt = time.Now().Unix()
	chatSessionRepo.Save(session)
	if isRemoveUser {
		_ = AdminService.RemoveUser(session.AdminId, session.UserId)
	}
	if updateTime {
		_ = AdminService.UpdateLimitTime(session.UserId, session.AdminId, 0)
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

// 客服给用户发消息后的会话有效期, 既用户在这时间内可以回复客服
func GetUserSessionSecond() int64 {
	setting := Settings[UserSessionDuration]
	dayFloat, err := strconv.ParseFloat(setting.GetValue(), 64)
	if err != nil {
		log.Fatal(err)
	}
	second := int64(dayFloat* 24 * 60 * 60)
	return second
}
// 用户给客服发消息后的会话有效期, 既客服在这时间内可以回复用户
func GetServiceSessionSecond() int64 {
	setting := Settings[AdminSessionDuration]
	dayFloat, err := strconv.ParseFloat(setting.GetValue(), 64)
	if err != nil {
		log.Fatal(err)
	}
	second := int64(dayFloat * 24 * 60 * 60)
	return second
}
