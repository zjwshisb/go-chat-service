package repositories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/models"
)

type AutoMessageRepo struct {

}


func (repo *AutoMessageRepo) Get(wheres []*Where, limit int, loads []string, orders ...string) []*models.AutoMessage {
	messages := make([]*models.AutoMessage, 0)
	query := databases.Db.
		Limit(limit).
		Scopes(AddWhere(wheres)).Scopes(AddLoad(loads)).Scopes(AddOrder(orders))
	query.Find(&messages)
	return messages
}

func (repo *AutoMessageRepo) First(wheres []*Where, orders ...string) *models.AutoMessage {
	message := &models.AutoMessage{}
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(wheres)).First(message)
	if result.Error != nil {
		return nil
	}
	return message
}

func (repo *AutoMessageRepo) Delete(model *models.AutoMessage) int64 {
	result := databases.Db.Delete(model)
	return result.RowsAffected
}

func (repo *AutoMessageRepo) Save(model *models.AutoMessage)  {
	databases.Db.Omit(clause.Associations).Save(model)
}

func (repo *AutoMessageRepo) Paginate(c *gin.Context, wheres []*Where, load []string, order ...string) *Pagination {
	messages := make([]*models.AutoMessage, 0)
	databases.Db.Order("id desc").
		Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(order)).
		Scopes(Paginate(c)).
		Find(&messages)
	var total int64
	databases.Db.Model(&models.AutoMessage{}).
		Scopes().
		Count(&total)
	return NewPagination(messages, total)
}

