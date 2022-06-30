package models

import (
	"github.com/duke-git/lancet/v2/random"
	"time"
	"ws/app/contract"
	"ws/app/databases"
)

type Admin struct {
	ID        int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	Username  string            `gorm:"string;size:255" `
	Password  string            `gorm:"string;size:255" `
	ApiToken  string            `gorm:"string;size:255" `
	Avatar    string            `gorm:"string;size:512"`
	GroupId   int64             `gorm:"group_id"`
	Setting   *AdminChatSetting `json:"setting" gorm:"foreignKey:admin_id"`
	IsSuper   bool              `gorm:"is_super"`
}

func (admin *Admin) GetGroupId() int64 {
	return admin.GroupId
}

// AccessTo 是否有user的权限
func (admin *Admin) AccessTo(user contract.User) bool {
	return admin.GetGroupId() == user.GetGroupId()
}

func (admin *Admin) GetIsSuper() bool {
	return admin.IsSuper
}

func (admin *Admin) GetPrimaryKey() int64 {
	return admin.ID
}

func (admin *Admin) GetAvatarUrl() string {
	return admin.GetSetting().Avatar
}

func (admin *Admin) GetUsername() string {
	return admin.Username
}

func (admin *Admin) GetSetting() *AdminChatSetting {
	if admin.Setting == nil {
		setting := &AdminChatSetting{}
		databases.Db.Model(admin).Association("Setting").Find(setting)
		if setting.Id == 0 {
			setting = &AdminChatSetting{
				AdminId:   admin.GetPrimaryKey(),
				Name:      admin.GetUsername(),
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			}
			databases.Db.Save(setting)
		}
		admin.Setting = setting
	}
	return admin.Setting
}

// GetChatName 客服名称
func (admin *Admin) GetChatName() string {
	setting := admin.GetSetting()
	if setting != nil {
		if setting.Name != "" {
			return setting.Name
		}
	}
	return admin.GetUsername()
}

// GetBreakMessage 断开消息
func (admin *Admin) GetBreakMessage(uid int64, sessionId uint64) *Message {
	return &Message{
		UserId:     uid,
		AdminId:    admin.GetPrimaryKey(),
		Type:       TypeNotice,
		Content:    admin.GetChatName() + "已断开服务",
		ReceivedAT: time.Now().Unix(),
		Source:     SourceSystem,
		SessionId:  sessionId,
		GroupId:    admin.GetGroupId(),
		ReqId:      random.RandString(20),
	}
}

func (admin *Admin) RefreshSetting() {
	setting := &AdminChatSetting{}
	_ = databases.Db.Model(admin).Association("Setting").Find(setting)
	admin.Setting = setting
}
