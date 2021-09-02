package repositories

import (
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/models"
)
type ChatSessionRepo struct {

}

func (session *ChatSessionRepo) Save(model *models.ChatSession) int64 {
	result := databases.Db.Omit(clause.Associations).Save(model)
	return result.RowsAffected
}
func (session *ChatSessionRepo) Delete(where []*Where) int64  {
	result := databases.Db.Scopes(AddWhere(where)).Delete(&models.ChatSession{})
	return result.RowsAffected
}
func (session *ChatSessionRepo) First(where []*Where, orders ...string) *models.ChatSession  {
	model := &models.ChatSession{}
	result := databases.Db.Scopes(AddWhere(where)).Scopes(AddOrder(orders)).First(model)
	if result.Error != nil {
		return nil
	}
	return model
}
func (session *ChatSessionRepo) Get(wheres []*Where, limit int, load []string, orders ...string) []*models.ChatSession {
	messages := make([]*models.ChatSession, 0)
	query := databases.Db.
		Limit(limit).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres))
	for _, relate := range load {
		query = query.Preload(relate)
	}
	query.Find(&messages)
	return messages
}
func (session *ChatSessionRepo) Update()  {

}
