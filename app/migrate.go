package app

import (
	"ws/app/databases"
	"ws/app/models"
)

func Migrate()  {

	_ = databases.Db.Migrator().CreateTable(&models.ChatSession{})
	_ = databases.Db.Migrator().CreateTable(&models.Message{})
	_ = databases.Db.Migrator().CreateTable(&models.AutoMessage{})
	_ = databases.Db.Migrator().CreateTable(&models.AdminChatSetting{})
	_ = databases.Db.Migrator().CreateTable(&models.ChatTransfer{})
	_ = databases.Db.Migrator().CreateTable(&models.AutoRule{})
	_ = databases.Db.Migrator().CreateTable(&models.AutoRuleScene{})
	_ = databases.Db.Migrator().CreateTable(&models.Admin{})
	_ = databases.Db.Migrator().CreateTable(&models.User{})
}
