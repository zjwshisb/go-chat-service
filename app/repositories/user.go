package repositories

import (
	"gorm.io/gorm"
	"ws/app/auth"
	"ws/app/databases"
	"ws/app/models"
)

type UserRepo struct {
	
}

func (session *UserRepo) Get(wheres []*Where, limit int, load []string, orders ...string) []*models.User {
	users := make([]*models.User, 0)
	databases.Db.
		Limit(limit).
		Scopes(AddOrder(orders)).
		Scopes(AddWhere(wheres)).Scopes(AddLoad(load)).
		Find(&users)
	return users
}

func GetUserById(id int64) (auth.User,  bool) {
	var user models.User
	var exist bool
	query := databases.Db.Where("id = ?", id).First(&user)
	if query.Error == gorm.ErrRecordNotFound {
		return &user, exist
	}
	exist = true
	return &user, exist
}

func GetUserByIds(ids []int64) (users []auth.User) {
	var ms []*models.User
	if len(ids) <= 0 {
		return
	}
	databases.Db.Find(&ms, ids)
	for _, u := range ms {
		users = append(users,u )
	}
	return users
}