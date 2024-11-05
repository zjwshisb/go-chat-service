package model

import "time"

type AutoMessageListItem struct {
	Id         uint      `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	Url        string    `json:"url"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RulesCount uint      `json:"rules_count"`
}
