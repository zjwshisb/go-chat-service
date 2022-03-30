package repositories

import (
	"time"
	"ws/app/databases"
	"ws/app/models"
)

type messageRepo struct {
	Repository[models.Message]
}

func (repo *messageRepo) GetUnSend(wheres []*Where) []*models.Message {
	wheres = append(wheres, &Where{
		Filed: "admin_id = ?",
		Value: 0,
	}, &Where{
		Filed: "source = ?",
		Value: models.SourceUser,
	})
	return repo.Get(wheres, -1, []string{}, []string{"id desc"})
}

func (repo *messageRepo) NewNotice(session *models.ChatSession, content string) *models.Message {
	return &models.Message{
		UserId:     session.UserId,
		AdminId:    session.AdminId,
		Type:       models.TypeNotice,
		GroupId:    session.GroupId,
		Content:    content,
		ReceivedAT: time.Now().Unix(),
		Source:     models.SourceSystem,
		SessionId:  session.Id,
		ReqId:      databases.GetSystemReqId(),
	}
}
