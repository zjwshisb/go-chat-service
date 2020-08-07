package database

import "ws/model"

func Migrate()  {
	Db.AutoMigrate(&model.User{})
}
