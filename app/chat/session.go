package chat

import (
	"log"
	"strconv"
	"time"
	"ws/app/models"
	"ws/app/repositories"
)
// 创建会话
func CreateSession(uid int64, t int) *models.ChatSession {
	session := &models.ChatSession{}
	session.UserId = uid
	session.QueriedAt = time.Now().Unix()
	session.AdminId = 0
	session.Type = t
	_ = chatSessionRepo.Save(session)
	return session
}
// 获取会话
func GetSession(uid int64, adminId int64) *models.ChatSession {
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
