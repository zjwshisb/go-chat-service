package migrate

import (
	"ws/db"
	"ws/internal/models"
)

func Migrate()  {
	db.Db.AutoMigrate(&models.Message{},
		&models.ServerUser{},
		&models.User{})

}
