package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"ws/app/databases"
	"ws/app/json"
)

const (
	IsAutoTransfer = "is-auto-transfer"
	AdminSessionDuration = "admin-session-duration"
	UserSessionDuration = "user-session-duration"
	MinuteToBreak = "minute-to-break"
)

const Key = "chat:%s:setting"

var SettingService = &settingService{
	Values: map[string]*SettingField{},
}

type settingService struct {
	Values map[string]*SettingField
}




type SettingField struct {
	Name string
	Title string
	val string
	Options map[string]string
	defVal string
	Validator func(val string, field *SettingField) error
}

func (field *SettingField) ToJson() *json.SettingField  {
	return &json.SettingField{
		Name:    field.Name,
		Title:   field.Title,
		Value:   field.GetValue(),
		Options: field.Options,
	}
}
func (field *SettingField) GetValue() string  {
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

func (field *SettingField) SetValue(val string) error {
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


// 离线时超过多久就自动断开会话
func (settingService *settingService) GetOfflineDuration() int64 {
	setting := settingService.Values[MinuteToBreak]
	minuteStr := setting.GetValue()
	minute, err := strconv.ParseInt(minuteStr, 10,64)
	if err != nil {
		log.Fatal(err)
	}
	return minute * 60
}
// 客服给用户发消息后的会话有效期, 既用户在这时间内可以回复客服
func (settingService *settingService) GetUserSessionSecond() int64 {
	setting := settingService.Values[UserSessionDuration]
	dayFloat, err := strconv.ParseFloat(setting.GetValue(), 64)
	if err != nil {
		log.Fatal(err)
	}
	second := int64(dayFloat* 24 * 60 * 60)
	return second
}
// 用户给客服发消息后的会话有效期, 既客服在这时间内可以回复用户
func (settingService *settingService) GetServiceSessionSecond() int64 {
	setting := settingService.Values[AdminSessionDuration]
	dayFloat, err := strconv.ParseFloat(setting.GetValue(), 64)
	if err != nil {
		log.Fatal(err)
	}
	second := int64(dayFloat * 24 * 60 * 60)
	return second
}

func (settingService *settingService) GetIsAutoTransferManual() bool {
	field, exist := settingService.Values[IsAutoTransfer]
	if !exist {
		return true
	}
	return field.GetValue() == "1"
}

func init() {
	SettingService.Values[IsAutoTransfer] = &SettingField{
		Name: IsAutoTransfer,
		Title: "是否自动转接人工客服",
		Options: map[string]string{
			"0": "否",
			"1": "是",
		},
		defVal: "1",
	}
	SettingService.Values[AdminSessionDuration] = &SettingField{
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
	SettingService.Values[UserSessionDuration] = &SettingField{
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
	SettingService.Values[MinuteToBreak] = &SettingField{
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

