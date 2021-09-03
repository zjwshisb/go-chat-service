package repositories

import (
	"errors"
	"reflect"
)

type Pagination struct {
	Data interface{} `json:"data"`
	Total int64 `json:"total"`
	Success bool `json:"success"`
}

func (p *Pagination) DataFormat(f func(i interface{}) interface{}) error {
	if reflect.TypeOf(p.Data).Kind() != reflect.Slice {
		return errors.New("data is not slice")
	}
	s := reflect.ValueOf(p.Data)
	data := make([]interface{}, s.Len())
	for i:=0; i< s.Len(); i ++ {
		item := s.Index(i)
		data[i] = f(item.Interface())
	}
	p.Data = data
	return nil
}


func NewPagination(data interface{}, total int64) *Pagination {
	return &Pagination{
		Data: data,
		Total: total,
		Success: true,
	}
}
