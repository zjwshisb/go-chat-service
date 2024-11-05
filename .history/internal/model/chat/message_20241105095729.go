package chat

type Message struct {
	Id         int64  `json:"id"`
	UserId     int    `json:"user_id"`
	AdminId    int    `json:"admin_id"`
	AdminName  string `json:"admin_name"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	ReceivedAT int64  `json:"received_at"`
	Source     int    `json:"source"`
	ReqId      string `json:"req_id"`
	IsSuccess  bool   `json:"is_success"`
	IsRead     bool   `json:"is_read"`
	Avatar     string `json:"avatar"`
	Username   string `json:"username"`
}
