package model

import "github.com/gogf/gf/v2/os/gtime"

type AutoMessageListItem struct {
	Id         uint        `json:"id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Content    string      `json:"content"`
	Url        string      `json:"url"`
	Title      string      `json:"title"`
	CreatedAt  *gtime.Time `json:"created_at"`
	UpdatedAt  *gtime.Time `json:"updated_at"`
	RulesCount uint        `json:"rules_count"`
}
