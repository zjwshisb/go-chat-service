package mq

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
	"ws/app/databases"
)

func newRedisMq() MessageQueue {
	return &RedisMq{}
}

type RedisSubscribe struct {
	channel <-chan *redis.Message
}

func (s *RedisSubscribe) Close() {
}

func (s *RedisSubscribe) ReceiveMessage() gjson.Result {
	message := <-s.channel
	result := gjson.Parse(message.Payload)
	return result
}

type RedisMq struct {
}

func (m *RedisMq) Publish(channel string, p *Payload) error {
	ctx := context.Background()
	cmd := databases.Redis.Publish(ctx, channel, p)
	return cmd.Err()
}

func (m *RedisMq) Subscribe(channel string) SubScribeChannel {
	ctx := context.Background()
	sub := databases.Redis.Subscribe(ctx, channel)
	return &RedisSubscribe{
		channel: sub.Channel(),
	}
}
