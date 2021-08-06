package models

import (
	"strings"
	"time"
	"ws/app/json"
	"ws/app/util"
)

const (
	MatchTypeAll  = "all"
	MatchTypePart = "part"

	MatchEnter             = "enter"
	MatchServiceAllOffLine = "u-offline"

	ReplyTypeMessage  = "message"
	ReplyTypeTransfer = "transfer"
	ReplyTypeEvent = "event"

	SceneNotAccepted = "not-accepted"
	SceneAdminOnline = "admin-online"
	SceneAdminOffline = "admin-offline"

	EventBreak = "break"
)
var ScenesOptions = []*json.Options{
	{
		Value: SceneNotAccepted,
		Label: "用户未被客服接入",
	},
	{
		Value: SceneAdminOnline,
		Label: "用户已接入且客服在线",
	},
	{
		Value: SceneAdminOffline,
		Label: "用户已接入且客服离线",
	},
}
var EventOptions = []*json.Options{
	{
		Value: EventBreak,
		Label: "断开当前会话",
	},
}
type AutoRuleScene struct {
	ID uint `json:"-"`
	Name string
	RuleId uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
type AutoRule struct {
	ID        uint
	Name      string       `gorm:"size:255" `
	Match     string       `gorm:"size:32"`
	MatchType string       `gorm:"size:20"`
	ReplyType string       `gorm:"size:20" `
	MessageId uint         `gorm:"index"`
	Key string `gorm:"key" json:"key"`
	IsSystem  uint8        `gorm:"is_system"`
	Sort      uint8        `gorm:"sort"`
	IsOpen    bool         `gorm:"is_open"`
	Count     uint         `gorm:"not null;default:0"`
	Scenes  []*AutoRuleScene `gorm:"foreignKey:RuleId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Message   *AutoMessage `json:"message" gorm:"foreignKey:MessageId"`
}

type AutoRuleJson struct {
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	Match     string       `json:"match"`
	MatchType string       `json:"match_type"`
	ReplyType string       `json:"reply_type"`
	MessageId uint         `json:"message_id"`
	Key string `gorm:"key" json:"key"`
	Sort      uint8        `json:"sort"`
	IsOpen    bool         `json:"is_open"`
	Count     uint         `json:"count"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	EventLabel string `json:"event_label"`
	Message   *AutoMessage `json:"message"`
	Scenes []string `json:"scenes"`
	ScenesLabel string `json:"scenes_label"`
}
// 是否匹配
func (rule *AutoRule) IsMatch(str string) bool  {
	switch rule.MatchType {
	case MatchTypeAll:
		return rule.Match == str
	case MatchTypePart:
		return strings.Contains(str, rule.Match)
	}
	return false
}
// 场景
func (rule *AutoRule) SceneInclude(str string) bool {
	for _, s := range rule.Scenes {
		if s.Name == str {
			return true
		}
	}
	return false
}
// 事件名称
func (rule *AutoRule) GetEventLabel() string {
	if rule.ReplyType == ReplyTypeEvent {
		for _, o := range EventOptions {
			if o.Value == rule.Key {
				return o.Label
			}
		}
	}
	return ""
}

func (rule *AutoRule) ToJson()  *AutoRuleJson  {
	scenesSli := make([]string, 0)
	scenesLabel := ""
	for _, scene := range rule.Scenes {
		scenesSli = append(scenesSli, scene.Name)
		for _, s := range ScenesOptions {
			if s.Value == scene.Name {
				if scenesLabel == "" {
					scenesLabel = s.Label
				} else {
					scenesLabel += "|" + s.Label
				}
			}
		}
	}
	return &AutoRuleJson{
		ID:          rule.ID,
		Name:        rule.Name,
		Match:       rule.Match,
		MatchType:   rule.MatchType,
		ReplyType:   rule.ReplyType,
		MessageId:   rule.MessageId,
		Key:         rule.Key,
		Sort:        rule.Sort,
		IsOpen:      rule.IsOpen,
		Count:       rule.Count,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
		Message:     rule.Message,
		EventLabel:  rule.GetEventLabel(),
		Scenes:      scenesSli,
		ScenesLabel: scenesLabel,
	}
}

func (rule *AutoRule) GetReplyMessage(uid int64) (message *Message) {
	if rule.Message != nil {
		message = &Message{
			UserId:     uid,
			AdminId:  0,
			Type:       rule.Message.Type,
			Content:    rule.Message.Content,
			ReceivedAT: time.Now().Unix(),
			SendAt:     0,
			Source:     SourceSystem,
			ReqId:      util.CreateReqId(),
			IsRead:     true,
			Avatar:     util.SystemAvatar(),
		}
	}
	return
}
