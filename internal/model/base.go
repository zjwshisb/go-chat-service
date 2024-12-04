package model

type Paginator[I any] struct {
	Items    []I
	Total    int
	IsSimple bool
}

type UnreadCount struct {
	Count  uint
	UserId uint
}
