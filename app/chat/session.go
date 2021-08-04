package chat

import (
	"log"
	"strconv"
	"time"
	"ws/app/databases"
	"ws/app/models"
)

func CreateSession(uid int64, t int) *models.ChatSession {
	session := &models.ChatSession{}
	session.UserId = uid
	session.QueriedAt = time.Now().Unix()
	session.AdminId = 0
	session.Type = t
	databases.Db.Save(session)
	return session
}
// 获取会话
func GetSession(uid int64, adminId int64) *models.ChatSession {
	session := &models.ChatSession{}
	databases.Db.Where("user_id = ?" , uid).
		Where("admin_id = ?", adminId).
		Order("id desc").First(session)
	if session.Id <= 0 {
		return nil
	}
	return session
}

// 客服给用户发消息的会话有效期, 既用户在这时间内可以回复客服
func GetUserSessionSecond() int64 {
	setting := Settings[UserSessionDuration]
	dayFloat, err := strconv.ParseFloat(setting.GetValue(), 64)
	if err != nil {
		log.Fatal(err)
	}
	second := int64(dayFloat* 24 * 60 * 60)
	return second
}
// 用户给客服发消息的会话有效期, 既客服在这时间内可以回复用户
func GetServiceSessionSecond() int64 {
	setting := Settings[AdminSessionDuration]
	dayFloat, err := strconv.ParseFloat(setting.GetValue(), 64)
	if err != nil {
		log.Fatal(err)
	}
	second := int64(dayFloat * 24 * 60 * 60)
	return second
}
