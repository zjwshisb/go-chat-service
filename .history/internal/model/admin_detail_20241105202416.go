package model

type AdminDetailInfo struct {
	Avatar        string `json:"avatar"`
	Username      string `json:"username"`
	Online        bool   `json:"online"`
	Id            uint   `json:"id"`
	AcceptedCount uint   `json:"accepted_count"`
}
