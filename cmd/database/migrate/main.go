package main

import (
	"fmt"
	"ws/internal/databases"
	"ws/internal/models"
)

func init()  {
	databases.Setup()
}
func main()  {
	err := databases.Db.AutoMigrate(&models.BackendUser{}, &models.User{}, &models.Message{},
	&models.Setting{}, &models.AutoMessage{}, &models.AutoRule{})
	fmt.Println(err)
}
