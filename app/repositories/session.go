package repositories

import (
	"time"
	"ws/app/databases"
	"ws/app/models"
)

type chatSessionRepo struct {
	Repository[models.ChatSession]
}

func (session *chatSessionRepo) Create(uid int64, groupId int64, ty int8) *models.ChatSession {
	s := &models.ChatSession{}
	s.UserId = uid
	s.QueriedAt = time.Now().Unix()
	s.AdminId = 0
	s.Type = ty
	s.GroupId = groupId
	session.Save(s)
	return s
}

func (session *chatSessionRepo) GetWaitHandles() []*models.ChatSession {
	sessions := make([]*models.ChatSession, 0)
	databases.Db.
		Limit(-1).
		Scopes(AddWhere([]*Where{
			{
				Filed: "admin_id = ?",
				Value: 0,
			},
			{
				Filed: "canceled_at = ?",
				Value: 0,
			},
			{
				Filed: "type = ?",
				Value: models.ChatSessionTypeNormal,
			},
		})).
		Preload("Messages", "source = ?", models.SourceUser).
		Preload("User").
		Find(&sessions)
	return sessions
}

// FirstActiveByUser 获取有效会话
func (session *chatSessionRepo) FirstActiveByUser(uid int64, adminId int64) *models.ChatSession {
	s := session.Repository.First([]*Where{
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
	}, []string{"id desc"})
	return s
}
