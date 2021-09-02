package repositories

import (
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

