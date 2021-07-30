package databases

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"ws/configs"
)

var Db *gorm.DB
var Redis *redis.Client

func Setup() {
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
		})
	if err != nil {
		log.Fatal(err)
	}
	Db = db
	Redis = redis.NewClient(&redis.Options{
		Addr:     configs.Redis.Addr,
		Password: configs.Redis.Auth, // no password set
		DB:       0,                  // use default DB
	})
	cmd := Redis.Ping(context.Background())
	if cmd.Err() != nil {
		log.Fatal(cmd.Err().Error())
	}
}
