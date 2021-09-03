package repositories

import (
	"github.com/gin-gonic/gin"
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
	databases.Db.
		Limit(limit).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Find(&messages)
	return messages
}
func (session *ChatSessionRepo) Paginate(c *gin.Context, wheres []*Where, load []string, orders ...string) *Pagination {
	rules := make([]*models.ChatSession, 0)
	databases.Db.Order("id desc").
		Scopes(Paginate(c)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres)).
		Find(&rules)
	var total int64
	databases.Db.Model(&models.AutoRule{}).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return NewPagination(rules, total)
}

