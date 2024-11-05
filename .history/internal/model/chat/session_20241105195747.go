package chat

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type Session struct {
	Id          uint        `json:"id"`
	UserId      uint        `json:"-"`
	QueriedAt   *gtime.Time `json:"queried_at"`
	AcceptedAt  *gtime.Time `json:"accepted_at"`
	BrokenAt    *gtime.Time `json:"broken_at"`
	CanceledAt  *gtime.Time `json:"canceled_at"`
	AdminId     uint        `json:"admin_id"`
	UserName    string      `json:"username"`
	AdminName   string      `json:"admin_name"`
	TypeLabel   string      `json:"type_label"`
	Status      string      `json:"status"`
	StatusLabel string      `json:"status_label"`
	Rate        uint        `json:"rate"`
}
