package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/models"
)

type MessageRepo struct {
}

func (repo *MessageRepo) Save(message *models.Message)  {
	databases.Db.Omit(clause.Associations).Save(message)
}
func (repo *MessageRepo) First(id interface{}) *models.Message {
	message := &models.Message{}
	query := databases.Db.Find(message, id)
	if query.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return message
}
func (repo *MessageRepo) Get(wheres []*Where, limit int, loads []string, orders []string) []*models.Message {
	messages := make([]*models.Message, 0)
	query := databases.Db.
		Limit(limit).
		Scopes(AddWhere(wheres)).Scopes(AddLoad(loads)).Scopes(AddOrder(orders))
	query.Find(&messages)
	return messages
}

func (repo *MessageRepo) Update(wheres []*Where, values map[string]interface{}) int64 {
	query := databases.Db.Table("messages").Scopes(AddWhere(wheres))
	query.Updates(values)
	return query.RowsAffected
}

func (repo *MessageRepo) GetUnSend(wheres []*Where) []*models.Message {
	wheres = append(wheres, &Where{
		Filed: "admin_id = ?",
		Value: 0,
	}, &Where{
		Filed: "source = ?",
		Value: models.SourceUser,
	})
	return repo.Get(wheres, -1, []string{}, []string{"id desc"})
}
