package mq

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"ws/app/databases"
)

func newRedisMq() MessageQueue {
	return &RedisMq{}
}

type RedisSubscribe struct {
	channel <-chan *redis.Message
}

func (s *RedisSubscribe) Close()  {
}

func (s *RedisSubscribe) ReceiveMessage() (*Payload, error) {
	message := <-s.channel
	payload := &Payload{}
	err := json.Unmarshal([]byte(message.Payload), payload)
	return payload, err
}

type RedisMq struct {

}
func (m *RedisMq) Publish(channel string, p *Payload) error  {
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