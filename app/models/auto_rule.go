package models

import (
	"gorm.io/gorm"
	"strings"
	"time"
	"ws/app/databases"
	"ws/app/resource"
)

const (
	MatchTypeAll  = "all"
	MatchTypePart = "part"

	MatchEnter             = "enter"
	MatchAdminAllOffLine = "u-offline"

	ReplyTypeMessage  = "message"
	ReplyTypeTransfer = "transfer"
	ReplyTypeEvent    = "event"

	SceneNotAccepted  = "not-accepted"
	SceneAdminOnline  = "admin-online"
	SceneAdminOffline = "admin-offline"

	EventBreak = "break"
)

var ScenesOptions = []*resource.Options{
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
var EventOptions = []*resource.Options{
	{
		Value: EventBreak,
		Label: "断开当前会话",
	},
}

type AutoRuleScene struct {
	ID        uint `json:"-"`
	Name      string
	RuleId    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
type AutoRule struct {
	ID        uint
	Name      string           `gorm:"size:255" `
	Match     string           `gorm:"size:32"`
	MatchType string           `gorm:"size:20"`
	ReplyType string           `gorm:"size:20" `
	MessageId uint             `gorm:"index"`
	Key       string           `gorm:"key" json:"key"`
	IsSystem  uint8            `gorm:"is_system"`
	Sort      uint8            `gorm:"sort"`
	IsOpen    bool             `gorm:"is_open"`
	GroupId   int64            `gorm:"group_id"`
	Count     uint             `gorm:"not null;default:0"`
	Scenes    []*AutoRuleScene `gorm:"foreignKey:RuleId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Message   *AutoMessage `json:"message" gorm:"foreignKey:MessageId"`
}



func (rule *AutoRule) AddCount() {
	databases.Db.Model(rule).Update("count", gorm.Expr("count + 1"))
}

// IsMatch 是否匹配
func (rule *AutoRule) IsMatch(str string) bool {
	switch rule.MatchType {
	case MatchTypeAll:
		return rule.Match == str
	case MatchTypePart:
		return strings.Contains(str, rule.Match)
	}
	return false
}

// SceneInclude 场景
func (rule *AutoRule) SceneInclude(str string) bool {
	for _, s := range rule.Scenes {
		if s.Name == str {
			return true
		}
	}
	return false
}

// GetEventLabel 事件名称
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

func (rule *AutoRule) ToJson() *resource.AutoRule {
	scenesSli := make([]string, 0)
	scenesLabel := ""
	if rule.Scenes == nil {
		scenes := make([]*AutoRuleScene, 0, 0)
		databases.Db.Model(rule).Association("Scenes").Find(&scenes)
		rule.Scenes = scenes
	}
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
	var messageJson *resource.AutoMessage
	if rule.Message != nil {
		messageJson = rule.Message.ToJson()
	}
	return &resource.AutoRule{
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
		Message:     messageJson,
		EventLabel:  rule.GetEventLabel(),
		Scenes:      scenesSli,
		ScenesLabel: scenesLabel,
	}
}

func (rule *AutoRule) GetReplyMessage(uid int64) (message *Message) {
	if rule.Message == nil {
		autoMessage := &AutoMessage{}
		databases.Db.Model(rule).Association("Message").Find(autoMessage)
		rule.Message = autoMessage
	}
	if rule.Message.ID > 0{
		message = &Message{
			UserId:     uid,
			AdminId:    0,
			Type:       rule.Message.Type,
			Content:    rule.Message.Content,
			ReceivedAT: time.Now().Unix(),
			SendAt:     0,
			Source:     SourceSystem,
			ReqId:      databases.GetSystemReqId(),
			IsRead:     true,
		}
	}
	return
}
