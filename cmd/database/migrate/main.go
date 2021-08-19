package main

import (
	"log"
	"ws/app/databases"
	"ws/app/models"
)

func main() {
	err := databases.Db.Migrator().CreateTable(&models.ChatSession{})
	err = databases.Db.Migrator().CreateTable(&models.Message{})
	err = databases.Db.Migrator().CreateTable(&models.AutoMessage{})
	err = databases.Db.Migrator().CreateTable(&models.AdminChatSetting{})
	err = databases.Db.Migrator().CreateTable(&models.ChatTransfer{})
	err = databases.Db.Migrator().CreateTable(&models.AutoRule{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("migrate success")
	}
}
