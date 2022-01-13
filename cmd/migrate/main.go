package migrate

import (
	"github.com/spf13/cobra"
	"log"
	"ws/app/databases"
	"ws/app/models"
)

func printErr(err error)  {
	if err != nil {
		log.Fatalf("migrate error: %v \n", err)
	}
}
func NewMigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:                        "migrate",
		Short:                      "create table",
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
			rules := []models.AutoRule{
				{
					Name: "用户进入客服系统时",
					MatchType: models.MatchTypeAll,
					Match: models.MatchEnter,
					ReplyType: models.ReplyTypeMessage,
					IsSystem: 1,
					GroupId: 1,
				},
				{
					Name: "当转接到人工客服而没有客服在线时(如不设置则继续转接到人工客服)",
					MatchType: models.MatchTypeAll,
					Match: models.MatchAdminAllOffLine,
					ReplyType: models.ReplyTypeMessage,
					IsSystem: 1,
					GroupId: 1,
				},
			}
			for _, rule := range rules {
				var exist int64
				databases.Db.Model(&models.AutoRule{}).
					Where("group_id", 1).
					Where("`match`=?" , rule.Match).Count(&exist)
				if exist == 0 {
					databases.Db.Save(&rule)
				}
			}
		},
	}
}

