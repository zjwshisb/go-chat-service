package models

import "time"

type ShortcutReply struct {
	Id int `json:"id" grom:"primaryKey,autoIncrement"`
	UserId int64 `json:"user_id"`
	Content string `json:"content"`
	CreatedAt time.Time 
	UpdatedAt time.Time
	Count int `json:"count"`
}