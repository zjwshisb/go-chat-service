package repositories

import (
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/models"
)

type adminRepo struct {
	Repository[models.Admin]
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
