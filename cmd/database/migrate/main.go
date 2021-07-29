package main

import (
	"log"
	"ws/app/databases"
	"ws/app/models"
)

func init() {
	databases.Setup()
}
func main() {
	err := databases.Db.AutoMigrate(
		&models.ChatSession{},
		&models.Message{},
		&models.AutoMessage{},
		&models.AdminChatSetting{},
		&models.ChatTransfer{},
		&models.AutoRule{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("migrate success")
	}
}
