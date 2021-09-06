package repositories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/models"
)

type AdminRepo struct {

}

func (repo *AdminRepo) First(wheres []*Where, orders ...string) *models.Admin {
	rule := &models.Admin{}
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(wheres)).First(rule)
	if result.Error != nil {
		return nil
	}
	return rule
}
func (repo *AdminRepo) SaveSetting(setting *models.AdminChatSetting)  {
	databases.Db.Omit(clause.Associations).Save(setting)
}
func (repo *AdminRepo) Save(admin *models.Admin)  {
	databases.Db.Omit(clause.Associations).Save(admin)
}
func (repo *AdminRepo) Get(wheres []*Where, limit int, load []string, orders ...string) []*models.Admin {
	admins := make([]*models.Admin, 0 )
	databases.Db.Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(orders)).Limit(limit).Find(&admins)
	return admins
}
func (repo *AdminRepo) Paginate(c *gin.Context, wheres []*Where, load []string, order ...string) *Pagination {
	rules := make([]*models.Admin, 0)
	databases.Db.
		Scopes(Paginate(c)).
		Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(order)).
		Find(&rules)
	var total int64
	databases.Db.Model(&models.Admin{}).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return NewPagination(rules, total)
}