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
func (session *ChatSessionRepo) First(where []*Where, orders []string) *models.ChatSession  {
	model := &models.ChatSession{}
	result := databases.Db.Scopes(AddWhere(where)).Scopes(AddOrder(orders)).Find(model)
	if result.RowsAffected == 0 {
		return nil
	}
	return model
}

func (session *ChatSessionRepo) GetWaitHandles() []*models.ChatSession  {
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
	})).
		Preload("Messages", "source = ?" , models.SourceUser).
		Preload("User").
		Find(&sessions)
	return sessions
}
func (session *ChatSessionRepo) Get(wheres []*Where, limit int, load []string, orders []string) []*models.ChatSession {
	sessions := make([]*models.ChatSession, 0)
	databases.Db.
		Limit(limit).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Find(&sessions)
	return sessions
}
func (session *ChatSessionRepo) Paginate(c *gin.Context, wheres []*Where, load []string, orders []string) *Pagination {
	sessions := make([]*models.ChatSession, 0)
	databases.Db.Order("id desc").
		Scopes(Paginate(c)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres)).
		Find(&sessions)
	var total int64
	databases.Db.Model(&models.ChatSession{}).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return NewPagination(sessions, total)
}

