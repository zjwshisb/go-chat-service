package migrate

import (
	"ws/db"
	"ws/models"
)

func Run()  {
	db.Db.AutoMigrate(&models.Message{},
		&models.ServerUser{},
		&models.User{})
}
