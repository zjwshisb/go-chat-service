package models

type ChatSession struct {
	Id uint64 `gorm:"primaryKey" json:"id"`
	UserId int64 `gorm:"index"`
	QueriedAt int64 
	AcceptedAt int64 
	BrokeAt int64
	AdminId int64 `gorm:"index"`
	Admin  *Admin `gorm:"foreignKey:admin_id"`
	User       *User        `gorm:"foreignKey:user_id"`
}

func (chatSession *ChatSession) ToJson() *ChatSessionJson {
	var userName, adminName string
	if chatSession.Admin != nil {
		adminName = chatSession.Admin.Username
	}
	if chatSession.User != nil {
		userName = chatSession.User.Username
	}
	return &ChatSessionJson{
		Id:          chatSession.Id,
		UserId:      chatSession.UserId,
		QueriedAt:   chatSession.QueriedAt * 1000,
		AcceptedAt:  chatSession.AcceptedAt * 1000,
		BrokeAt:     chatSession.BrokeAt * 1000,
		AdminId:   chatSession.AdminId,
		UserName:    userName,
		AdminName: adminName,
	}
}

type ChatSessionJson struct {
	Id uint64 `json:"id"`
	UserId int64 `json:"-"`
	QueriedAt int64 `json:"queried_at"`
	AcceptedAt int64 `json:"accepted_at"`
	BrokeAt int64	`json:"broke_at"`
	AdminId int64  `json:"Admin_id"`
	UserName string `json:"user_name"`
	AdminName string `json:"admin_name"`
}

