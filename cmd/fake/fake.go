package fake

import (
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"ws/app/databases"
	"ws/app/models"
)

var defaultGroupId int64 = 1


func NewFakeCommand() *cobra.Command {
	return &cobra.Command{
		Use:                        "fake",
		Short:                      "Create fake data",
		Run: func(cmd *cobra.Command, args []string) {
			databases.MysqlSetup()
			for i:=1; i<=20; i++ {
				adminName :=  "admin" + strconv.Itoa(i)
				var count int64
				databases.Db.Model(&models.Admin{}).Where("username = ?" , adminName).Count(&count)
				if count == 0 {
					password, _ :=  bcrypt.GenerateFromPassword([]byte(adminName),bcrypt.DefaultCost)
					admins := &models.Admin{
						Username:  adminName,
						Password: string(password),
						GroupId: defaultGroupId,
					}
					databases.Db.Save(admins)
				}
				username :=  "user" + strconv.Itoa(i)
				databases.Db.Model(&models.User{}).Where("username = ?" , username).Count(&count)
				if count <= 0 {
					password, _ :=  bcrypt.GenerateFromPassword([]byte(username),bcrypt.DefaultCost)
					admins := &models.User{
						Username:  username,
						Password: string(password),
						GroupId: defaultGroupId,
					}
					databases.Db.Save(admins)
				}
			}
		},
	}
}



