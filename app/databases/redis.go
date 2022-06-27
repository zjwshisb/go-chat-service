package databases

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Redis *redis.Client

func RedisSetup() {
	options := &redis.Options{
		Addr: viper.GetString("Redis.Addr"),
		DB:   viper.GetInt("Redis.Db"),
	}
	if viper.GetString("Redis.Auth") != "" {
		options.Password = ""
	}
	Redis = redis.NewClient(options)
	cmd := Redis.Ping(context.Background())
	if cmd.Err() != nil {
		panic(fmt.Errorf("redis error: %w \n", cmd.Err()))
	}
}
