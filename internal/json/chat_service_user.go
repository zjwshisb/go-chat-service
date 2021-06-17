package json


type ChatServiceUser struct {
	Avatar string `json:"avatar"`
	Username string `json:"username"`
	Online bool `json:"online"`
	Id int64 `json:"id"`
	TodayAcceptCount int64 `json:"today_accept_count"`
}