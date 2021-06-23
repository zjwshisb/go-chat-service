package json

type Message struct {
	Id         uint64 `json:"id"`
	UserId     int64  `json:"user_id"`
	ServiceId  int64  `json:"service_id"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	ReceivedAT int64  `json:"received_at"`
	IsServer   bool   `json:"is_server"`
	ReqId      int64  `json:"req_id"`
	IsSuccess  bool   `json:"is_success"`
	IsRead     bool   `json:"is_read"`
	Avatar     string `json:"avatar"`
}
