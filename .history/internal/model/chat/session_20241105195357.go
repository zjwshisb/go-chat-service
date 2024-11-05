package chat

type Session struct {
	Id          uint        `json:"id"`
	UserId      uint        `json:"-"`
	QueriedAt   *gtime.TIme `json:"queried_at"`
	AcceptedAt  *gtime.TIme `json:"accepted_at"`
	BrokeAt     *gtime.TIme `json:"broke_at"`
	CanceledAt  *gtime.TIme `json:"canceled_at"`
	AdminId     uint        `json:"admin_id"`
	UserName    string      `json:"username"`
	AdminName   string      `json:"admin_name"`
	TypeLabel   string      `json:"type_label"`
	Status      string      `json:"status"`
	StatusLabel string      `json:"status_label"`
	Rate        int         `json:"rate"`
}
