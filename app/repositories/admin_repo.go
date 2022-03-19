package repositories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/models"
)

type adminRepo struct {
}

func (repo *adminRepo) First(wheres []*Where, orders []string) *models.Admin {
	admin := &models.Admin{}
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(wheres)).First(admin)
	if result.Error != nil {
		return nil
	}
	return admin
}

func (repo *adminRepo) SaveSetting(setting *models.AdminChatSetting) {
	databases.Db.Omit(clause.Associations).Save(setting)
}
func (repo *adminRepo) UpdateSetting(setting *models.AdminChatSetting, column string, value interface{}) {
	databases.Db.Model(setting).Update(column, value)
}

func (repo *adminRepo) Save(admin *models.Admin) {
	databases.Db.Omit(clause.Associations).Save(admin)
}

func (repo *adminRepo) Get(wheres []*Where, limit int, load []string, orders []string) []*models.Admin {
	admins := make([]*models.Admin, 0)
	databases.Db.Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(orders)).Limit(limit).Find(&admins)
	return admins
}

func (repo *adminRepo) Paginate(c *gin.Context, wheres []*Where, load []string, order []string) *Pagination {
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
