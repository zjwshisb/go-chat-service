package chat

import (
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
