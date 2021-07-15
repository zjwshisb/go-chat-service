package chat

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"ws/internal/databases"
)

const (
	IsAutoTransfer = "is-auto-transfer"
)

const Key = "chat:%s:setting"

type Field struct {
	Name string
	Title string
	val string
	Options map[string]string
}

func (field *Field) GetValue() string  {
	if field.val == "" {
		ctx := context.Background()
		cmd := databases.Redis.Get(ctx, fmt.Sprintf(Key, field.Name))
		if cmd.Err() == redis.Nil {
			field.val =  "0"
		} else {
			field.val = "1"
		}
	}
	return field.val
}

func (field *Field) SetValue(val string) error {
	field.val = val
	ctx := context.Background()
	cmd := databases.Redis.Set(ctx, fmt.Sprintf(Key, field.Name), val, 0)
	return cmd.Err()
}

var Settings map[string]*Field

func init() {
	Settings = make(map[string]*Field)
	Settings[IsAutoTransfer] = &Field{
		Name: IsAutoTransfer,
		Title: "是否自动转接人工客服",
		Options: map[string]string{
			"0": "否",
			"1": "是",
		},
	}
}

