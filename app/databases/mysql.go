package databases

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"ws/config"
)

var Db *gorm.DB

func MysqlSetup() {
	env := config.GetEnv()

	var level logger.LogLevel
	switch strings.ToLower(env) {
	case "production":
		level = logger.Silent
	case "local":
		level = logger.Info
	case "test":
		level = logger.Info
	}
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("Mysql.Username"),
		viper.GetString("Mysql.Password"),
		viper.GetString("Mysql.Host"),
		viper.GetString("Mysql.Port"),
		viper.GetString("Mysql.Database"),
	)
	db, err := gorm.Open(mysql.Open(dns),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   logger.Default.LogMode(level),
		})
	if err != nil {
		panic(fmt.Errorf("mysql err: %w \n", err))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	Db = db
}
