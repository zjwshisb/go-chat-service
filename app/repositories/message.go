package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"ws/app/databases"
	"ws/app/models"
)

type messageRepo struct {
}

func (repo *messageRepo) Save(message *models.Message)  {
	databases.Db.Omit(clause.Associations).Save(message)
}
func (repo *messageRepo) First(id interface{}) *models.Message {
	message := &models.Message{}
	query := databases.Db.Find(message, id)
	if query.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return message
}
func (repo *messageRepo) Get(wheres []*Where, limit int, loads []string, orders []string) []*models.Message {
	messages := make([]*models.Message, 0)
	query := databases.Db.
		Limit(limit).
		Scopes(AddWhere(wheres)).Scopes(AddLoad(loads)).Scopes(AddOrder(orders))
	query.Find(&messages)
	return messages
}

func (repo *messageRepo) Update(wheres []*Where, values map[string]interface{}) int64 {
	query := databases.Db.Table("messages").Scopes(AddWhere(wheres))
	query.Updates(values)
	return query.RowsAffected
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
		Content:    content,
		ReceivedAT: time.Now().Unix(),
		Source:     models.SourceSystem,
		SessionId:  session.Id,
		ReqId:      databases.GetSystemReqId(),
	}
}