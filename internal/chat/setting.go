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
	AdminSessionDuration = "admin-session-duration"
	UserSessionDuration = "user-session-duration"
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
			field.val = cmd.Val()
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
	Settings[AdminSessionDuration] = &Field{
		Name: AdminSessionDuration,
		Title: "当用户给客服发消息时，客服多久没回复就断开会话",
		Options: map[string]string{
			"0.3333": "8小时",
			"0.1666": "4小时",
			"0.0833": "2小时",
			"0.0416": "1小时",
			"0.0208": "30分钟",
			"0.5": "12小时",
			"1": "1天",
		},
		defVal: "1",
	}
	Settings[UserSessionDuration] = &Field{
		Name: UserSessionDuration,
		Title: "当客服给用户发消息时，用户多久没回复就断开会话",
		Options: map[string]string{
			"0.3333": "8小时",
			"0.1666": "4小时",
			"0.0833": "2小时",
			"0.0416": "1小时",
			"0.0208": "30分钟",
			"0.5": "12小时",
			"1": "1天",
		},
		defVal: "0.0208",
	}
}

