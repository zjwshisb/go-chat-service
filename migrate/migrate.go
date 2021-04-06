package migrate

import (
	"ws/db"
	"ws/models"
)

func Migrate()  {
	db.Db.AutoMigrate(&models.Message{},
		&models.ServerUser{},
		&models.User{})

}
