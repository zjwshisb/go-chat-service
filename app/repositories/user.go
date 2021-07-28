package repositories

import (
	"gorm.io/gorm"
	"ws/app/auth"
	"ws/app/databases"
	"ws/app/models"
)


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