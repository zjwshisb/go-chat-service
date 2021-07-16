package chat

import (
	"context"
	"errors"
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
	defVal string
	Validator func(val string, field *Field) error
}

func (field *Field) GetValue() string  {
	if field.val == "" {
		ctx := context.Background()
		cmd := databases.Redis.Get(ctx, fmt.Sprintf(Key, field.Name))
		if cmd.Err() == redis.Nil {
			field.val = field.defVal
		} else {
			field.val = field.defVal
		}
	}
	return field.val
}

func (field *Field) SetValue(val string) error {
	for v := range field.Options {
		if v == val {
			field.val = val
			ctx := context.Background()
			cmd := databases.Redis.Set(ctx, fmt.Sprintf(Key, field.Name), val, 0)
			return cmd.Err()
		}
	}
	return errors.New("validated failed")
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
		defVal: "1",
	}
}

