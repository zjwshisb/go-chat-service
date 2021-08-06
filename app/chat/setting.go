package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"ws/app/databases"
)

const (
	IsAutoTransfer = "is-auto-transfer"
	AdminSessionDuration = "admin-session-duration"
	UserSessionDuration = "user-session-duration"
	MinuteToBreak = "minute-to-break"
)

const Key = "chat:%s:setting"

type FieldJson struct {
	Name string `json:"name"`
	Title string `json:"title"`
	Value string `json:"value"`
	Options map[string]string `json:"options"`
}

type Field struct {
	Name string
	Title string
	val string
	Options map[string]string
	defVal string
	Validator func(val string, field *Field) error
}

func (field *Field) ToJson() *FieldJson  {
	return &FieldJson{
		Name:    field.Name,
		Title:   field.Title,
		Value:   field.GetValue(),
		Options: field.Options,
	}
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
	Settings[MinuteToBreak] = &Field{
		Name:      MinuteToBreak,
		Title:     "客服离线多少分钟(用户发送消息时)自动断开会话",
		Options:   map[string]string{
			"5": "5分钟",
			"10": "10分钟",
			"15": "15分钟",
			"20": "20分钟",
			"30": "30分钟",
			"60": "60分钟",
		},
		defVal:    "10",
		Validator: nil,
	}
}

