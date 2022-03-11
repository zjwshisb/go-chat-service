package repositories

import (
	"ws/app/databases"
	"ws/app/models"
)

type userRepo struct {
	
}

func (repo *userRepo) Get(wheres []*Where, limit int, load []string, orders []string) []*models.User {
	users := make([]*models.User, 0)
	databases.Db.
		Limit(limit).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres)).Scopes(AddLoad(load)).
		Find(&users)
	return users
}

func (repo *userRepo) First(where []*Where, orders []string) *models.User {
	user := &models.User{}
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(where)).First(user)
	if result.Error != nil {
		return nil
	}
	return user
}
