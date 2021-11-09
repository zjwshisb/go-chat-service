package databases

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"ws/configs"
)

var Db *gorm.DB

func init()  {
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configs.Mysql.Username,
		configs.Mysql.Password,
		configs.Mysql.Host,
		configs.Mysql.Port,
		configs.Mysql.Name,
	)
	db, err := gorm.Open(mysql.Open(dns),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger: logger.Default.LogMode(logger.Error),
		})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	Db = db
}
