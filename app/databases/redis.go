package databases

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"strconv"
)


var Redis *redis.Client

func RedisSetup() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("Redis.Addr"),
		Password: viper.GetString("Redis.Auth"),
		DB:       viper.GetInt("Redis.Db"),
	})
	cmd := Redis.Ping(context.Background())
	if cmd.Err() != nil {
		panic(fmt.Errorf("redis error: %w \n", cmd.Err()))
	}
}

func GetSystemReqId() string {
	key := "system:req-id"
	ctx := context.Background()
	cmd := Redis.Incr(ctx, key)
	return "s" + strconv.FormatInt(cmd.Val(), 10)
}