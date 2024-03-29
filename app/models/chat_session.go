package models

import (
	"ws/app/databases"
	"ws/app/resource"
)

const ChatSessionTypeNormal = 0
const ChatSessionTypeTransfer = 1

type ChatSession struct {
	Id         uint64 `gorm:"primaryKey" json:"id"`
	UserId     int64  `gorm:"index"`
	QueriedAt  int64
	AcceptedAt int64
	CanceledAt   int64
	BrokeAt    int64
	GroupId   int64  `gorm:"group_id"`
	AdminId    int64  `gorm:"index"`
	Admin      *Admin `gorm:"foreignKey:admin_id"`
	Type       int8    `gorm:"default:0"`
	User       *User  `gorm:"foreignKey:user_id"`
	Messages []*Message `gorm:"foreignKey:session_id"`
}

func (chatSession *ChatSession) GetAdmin() *Admin  {
	if chatSession.Admin == nil {
		admin := &Admin{}
		databases.Db.Model(chatSession).Association("Admin").Find(admin)
		chatSession.Admin = admin
	}
	return chatSession.Admin
}
func (chatSession *ChatSession) GetUser() *User  {
	if chatSession.User == nil {
		user := &User{}
		_ = databases.Db.Model(chatSession).Association("User").Find(user)
		chatSession.User = user
	}
	return chatSession.User
}
func (chatSession *ChatSession) getTypeLabel() string {
	switch chatSession.Type {
	case ChatSessionTypeTransfer:
		return "转接"
	case ChatSessionTypeNormal:
		return "普通"
	default:
		return ""
	}
}
func (chatSession *ChatSession) getStatus() string  {
	if chatSession.CanceledAt > 0 {
		return "cancel"
	}
	if chatSession.AcceptedAt > 0 {
		return "accept"
	}
	return "wait"
}

func (chatSession *ChatSession) ToJson() *resource.ChatSession {
	return &resource.ChatSession{
		Id:         chatSession.Id,
		UserId:     chatSession.UserId,
		QueriedAt:  chatSession.QueriedAt * 1000,
		AcceptedAt: chatSession.AcceptedAt * 1000,
		BrokeAt:    chatSession.BrokeAt * 1000,
		CanceledAt: chatSession.CanceledAt * 1000,
		AdminId:    chatSession.AdminId,
		TypeLabel:  chatSession.getTypeLabel(),
		Status: chatSession.getStatus(),
		UserName:   chatSession.GetUser().Username,
		AdminName:  chatSession.GetAdmin().Username,
	}
}

