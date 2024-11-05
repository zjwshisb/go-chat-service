package model

import "time"

type AutoMessageListItem struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	Url        string    `json:"url"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RulesCount int       `json:"rules_count"`
}
