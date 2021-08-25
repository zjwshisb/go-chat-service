package models

import "time"

type ChatTransfer struct {
	Id          int64
	UserId      int64  `gorm:"index"`
	SessionId   uint64 `gorm:"index"`
	FromAdminId int64  `gorm:"index"`
	ToAdminId   int64  `gorm:"index"`
	Remark      string `gorm:"size:255"`
	IsAccepted  bool
	IsCanceled  bool
	CreatedAt   *time.Time
	AcceptedAt  *time.Time
	CanceledAt  *time.Time
	Session     *ChatSession `gorm:"foreignKey:session_id"`
	User        *User        `gorm:"foreignKey:user_id"`
	FromAdmin   *Admin       `gorm:"foreignKey:from_admin_id"`
	ToAdmin     *Admin       `gorm:"foreignKey:to_admin_id"`
}

func (transfer *ChatTransfer) ToJson() *ChatTransferJson {
	json := &ChatTransferJson{
		Id:         transfer.Id,
		SessionId:  transfer.SessionId,
		UserId:     transfer.UserId,
		Remark:     transfer.Remark,
		CreatedAt:  transfer.CreatedAt,
		AcceptedAt: transfer.AcceptedAt,
		CanceledAt: transfer.CanceledAt,
	}
	if transfer.FromAdmin != nil {
		json.FromAdminName = transfer.FromAdmin.GetUsername()
	}
	if transfer.User != nil {
		json.Username = transfer.User.GetUsername()
	}
	if transfer.ToAdmin != nil {
		json.ToAdminName = transfer.ToAdmin.GetUsername()
	}

	return json
}

type ChatTransferJson struct {
	Id            int64      `json:"id"`
	SessionId     uint64     `json:"session_id"`
	UserId        int64      `json:"user_id"`
	Remark        string     `json:"remark"`
	FromAdminName string     `json:"from_admin_name"`
	ToAdminName   string     `json:"to_admin_name"`
	Username      string     `json:"username"`
	CreatedAt     *time.Time `json:"created_at"`
	AcceptedAt    *time.Time `json:"accepted_at"`
	CanceledAt    *time.Time `json:"canceled_at"`
}
