package chat

type Session struct {
	Id          uint64 `json:"id"`
	UserId      int    `json:"-"`
	QueriedAt   int64  `json:"queried_at"`
	AcceptedAt  int64  `json:"accepted_at"`
	BrokeAt     int64  `json:"broke_at"`
	CanceledAt  int64  `json:"canceled_at"`
	AdminId     int    `json:"admin_id"`
	UserName    string `json:"username"`
	AdminName   string `json:"admin_name"`
	TypeLabel   string `json:"type_label"`
	Status      string `json:"status"`
	StatusLabel string `json:"status_label"`
	Rate        int    `json:"rate"`
}
