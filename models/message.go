package models

import "google.golang.org/genproto/googleapis/type/datetime"

type message struct {
	Id int64
	UserId int64
	ServiceId int64
	Type string
	content string
	CreatedAT datetime.DateTime
	ReqId int64
	IsServer bool
}
