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
	err := databases.Db.AutoMigrate(&models.ServiceUser{}, &models.User{}, &models.Message{})
	fmt.Println(err)
}
