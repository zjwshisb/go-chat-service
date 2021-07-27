package json

type Message struct {
	Id         uint64 `json:"id"`
	UserId     int64  `json:"user_id"`
	AdminId  int64  `json:"admin_id"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	ReceivedAT int64  `json:"received_at"`
	Source   int8   `json:"source"`
	ReqId      int64  `json:"req_id"`
	IsSuccess  bool   `json:"is_success"`
	IsRead     bool   `json:"is_read"`
	Avatar     string `json:"avatar"`
}


