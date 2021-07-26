package models

type ChatSession struct {
	Id uint64 `gorm:"primaryKey" json:"id"`
	UserId int64 `gorm:"index"`
	QueriedAt int64 
	AcceptedAt int64 
	BrokeAt int64
	ServiceId int64 `gorm:"index"`
	BackendUser *BackendUser `gorm:"foreignKey:service_id"`
	User       *User        `gorm:"foreignKey:user_id"`
}

func (chatSession *ChatSession) ToJson() *ChatSessionJson {
	var userName, serviceName string
	if chatSession.BackendUser != nil {
		serviceName = chatSession.BackendUser.Username
	}
	if chatSession.User != nil {
		userName = chatSession.User.Username
	}
	return &ChatSessionJson{
		Id:          chatSession.Id,
		UserId:      chatSession.ServiceId,
		QueriedAt:   chatSession.QueriedAt * 1000,
		AcceptedAt:  chatSession.AcceptedAt * 1000,
		BrokeAt:     chatSession.BrokeAt * 1000,
		ServiceId:   chatSession.ServiceId,
		UserName:    userName,
		ServiceName: serviceName,
	}
}

type ChatSessionJson struct {
	Id uint64 `json:"id"`
	UserId int64 `json:"-"`
	QueriedAt int64 `json:"queried_at"`
	AcceptedAt int64 `json:"accepted_at"`
	BrokeAt int64	`json:"broke_at"`
	ServiceId int64  `json:"service_id"`
	UserName string `json:"user_name"`
	ServiceName string `json:"service_name"`
}

