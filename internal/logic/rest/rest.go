package rest

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
)

type RDao interface {
	DB() gdb.DB
	Table() string
	Group() string
	Ctx(ctx context.Context) *gdb.Model
	Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error)
}

type Rest[C any] struct {
	RDao RDao
}

func (s Rest[C]) Get(ctx context.Context, where any) (res []C) {
	query := s.RDao.Ctx(ctx)
	if where != nil {
		query = query.Where(where)
	}
	query.Scan(&res)
	return
}
