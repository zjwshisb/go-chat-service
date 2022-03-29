package migrate

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"ws/app/databases"
	"ws/app/models"
)

var defaultGroupId int64 = 1

func printErr(err error) {
	if err != nil {
		log.Fatalf("migrate error: %v \n", err)
	}
}
func getSettings() []*models.ChatSetting {
	s := make([]*models.ChatSetting, 0, 0)
	options1, _ := json.Marshal([]map[string]string{
		{
			"label": "是",
			"value": "1",
		},
		{
			"label": "否",
			"value": "0",
		},
	})
	s = append(s, &models.ChatSetting{
		Name:      models.IsAutoTransfer,
		Title:     "是否自动转接人工客服",
		GroupId:   defaultGroupId,
		Value:     "1",
		Options:   string(options1),
		CreatedAt: nil,
		UpdatedAt: nil,
		Type:      "select",
	})
	options2, _ := json.Marshal([]map[string]string{
		{
			"label": "5分钟",
			"value": "5",
		},
		{
			"label": "10分钟",
			"value": "10",
		},
		{
			"label": "30分钟",
			"value": "30",
		},
		{
			"label": "60分钟",
			"value": "60",
		},
	})
	s = append(s, &models.ChatSetting{
		Name:      models.MinuteToBreak,
		Title:     "客服离线多少分钟(用户发送消息时)自动断开会话",
		GroupId:   defaultGroupId,
		Value:     "30",
		Options:   string(options2),
		CreatedAt: nil,
		UpdatedAt: nil,
		Type:      "select",
	})
	s = append(s, &models.ChatSetting{
		Name:      models.SystemAvatar,
		Title:     "系统头像",
		GroupId:   defaultGroupId,
		Value:     "",
		Options:   "",
		Type:      "image",
		CreatedAt: nil,
		UpdatedAt: nil,
	})
	s = append(s, &models.ChatSetting{
		Name:      models.SystemName,
		Title:     "系统名称",
		GroupId:   defaultGroupId,
		Value:     "",
		Options:   "",
		Type:      "text",
		CreatedAt: nil,
		UpdatedAt: nil,
	})
	return s
}

func NewMigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "create table",
		Run: func(cmd *cobra.Command, args []string) {
			databases.MysqlSetup()
			err := databases.Db.Migrator().CreateTable(&models.ChatSession{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.Message{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.AutoMessage{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.AdminChatSetting{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.ChatTransfer{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.AutoRule{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.AutoRuleScene{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.Admin{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.User{})
			printErr(err)
			err = databases.Db.Migrator().CreateTable(&models.ChatSetting{})
			rules := []models.AutoRule{
				{
					Name:      "用户进入客服系统时",
					MatchType: models.MatchTypeAll,
					Match:     models.MatchEnter,
					ReplyType: models.ReplyTypeMessage,
					IsSystem:  1,
					GroupId:   defaultGroupId,
				},
				{
					Name:      "当转接到人工客服而没有客服在线时(如不设置则继续转接到人工客服)",
					MatchType: models.MatchTypeAll,
					Match:     models.MatchAdminAllOffLine,
					ReplyType: models.ReplyTypeMessage,
					IsSystem:  1,
					GroupId:   defaultGroupId,
				},
			}
			for _, rule := range rules {
				var exist int64
				databases.Db.Model(&models.AutoRule{}).
					Where("group_id", defaultGroupId).
					Where("`match`=?", rule.Match).Count(&exist)
				if exist == 0 {
					databases.Db.Save(&rule)
				}
			}
			for _, setting := range getSettings() {
				databases.Db.Save(setting)
			}

		},
	}
}
