package databases

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"ws/configs"
)


var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     configs.Redis.Addr,
		Password: configs.Redis.Auth,
		DB:       configs.Redis.Db,
	})
	cmd := Redis.Ping(context.Background())
	if cmd.Err() != nil {
		log.Fatal(cmd.Err().Error())
	}
}
