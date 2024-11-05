// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// CustomerChatMessages is the golang structure for table customer_chat_messages.
type CustomerChatMessages struct {
	Id         int64  `json:"id"         ` //
	UserId     int    `json:"userId"     ` //
	AdminId    int    `json:"adminId"    ` //
	Type       string `json:"type"       ` //
	Content    string `json:"content"    ` //
	ReceivedAt int64  `json:"receivedAt" ` //
	CustomerId int    `json:"customerId" ` //
	SendAt     int64  `json:"sendAt"     ` //
	Source     int    `json:"source"     ` //
	SessionId  uint64 `json:"sessionId"  ` //
	ReqId      string `json:"reqId"      ` //
	ReadAt     int64  `json:"readAt"     ` //
}
