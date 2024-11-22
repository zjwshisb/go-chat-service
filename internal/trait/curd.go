package trait

import (
	"context"
	"database/sql"
	"gf-chat/api"
	"gf-chat/internal/model"

	"github.com/gogf/gf/v2/database/gdb"
)

type ctx = context.Context

type ICurd[R any] interface {
	Save(ctx ctx, data *R) (id int64, err error)
	Find(ctx ctx, primaryKey any) (model *R, err error)
	All(ctx ctx, where any, with []any, order any) (items []*R, err error)
	First(ctx ctx, where any, order ...string) (model *R, err error)
	Paginate(ctx ctx, where any, p api.Paginate, with []any, order any) (paginator *model.Paginator[R], err error)
	Insert(ctx ctx, data *R) (id int64, err error)
	Update(ctx ctx, where any, data any) (count int64, err error)
	Exists(ctx ctx, where any) (exists bool, err error)
	Delete(ctx ctx, primaryKey any) error
}

type IDao interface {
	DB() gdb.DB
	Table() string
	Group() string
	Ctx(ctx context.Context) *gdb.Model
	Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error)
}

type Curd[R any] struct {
	Dao IDao
}

func (c Curd[R]) Delete(ctx ctx, primaryKey any) error {
	_, err := c.Dao.Ctx(ctx).WherePri(primaryKey).Delete()
	return err
}

func (c Curd[R]) Find(ctx ctx, primaryKey any) (model *R, err error) {
	err = c.Dao.Ctx(ctx).WherePri(primaryKey).Scan(&model)
	if err != nil {
		return
	}
	if model == nil {
		err = sql.ErrNoRows
	}
	return
}

func (c Curd[R]) First(ctx ctx, where any, order ...string) (model *R, err error) {
	err = c.Dao.Ctx(ctx).Where(where).Order(order).Scan(&model)
	if err != nil {
		return
	}
	if model == nil {
		err = sql.ErrNoRows
	}
	return
}

func (c Curd[R]) Exists(ctx ctx, where any) (exists bool, err error) {
	return c.Dao.Ctx(ctx).Where(where).Exist()

}

func (c Curd[R]) All(ctx ctx, where any, with []any, order any) (items []*R, err error) {
	err = c.Dao.Ctx(ctx).Where(where).With(with...).Order(order).Scan(&items)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]*R, 0)
	}
	return
}

func (c Curd[R]) Save(ctx ctx, data *R) (id int64, err error) {
	result, err := c.Dao.Ctx(ctx).Save(data)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

func (c Curd[R]) Insert(ctx ctx, data *R) (id int64, err error) {
	result, err := c.Dao.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

func (c Curd[R]) Update(ctx ctx, where any, data any) (count int64, err error) {
	result, err := c.Dao.Ctx(ctx).Where(where).Data(data).Update()
	if err != nil {
		return
	}
	count, err = result.RowsAffected()
	return
}

func (c Curd[R]) SimplePaginate(ctx context.Context, where any, p api.Paginate, with []any, order any) (paginator *model.Paginator[R], err error) {
	query := c.Dao.Ctx(ctx)
	if where != nil {
		query = query.Where(where)
	}
	items := make([]R, 0)
	total := 0
	query = query.Page(p.Current, p.PageSize)
	err = query.With(with...).Order(order).Scan(&items)
	if err != nil {
		return
	}
	return &model.Paginator[R]{
		Items:    items,
		Total:    total,
		IsSimple: true,
	}, nil
}

func (c Curd[R]) Paginate(ctx context.Context, where any, p api.Paginate, with []any, order any) (paginator *model.Paginator[R], err error) {
	query := c.Dao.Ctx(ctx)
	if where != nil {
		query = query.Where(where)
	}
	items := make([]R, 0)
	total := 0
	query = query.Page(p.Current, p.PageSize)
	err = query.With(with...).Order(order).ScanAndCount(&items, &total, true)
	if err != nil {
		return
	}
	return &model.Paginator[R]{
		Items:    items,
		Total:    total,
		IsSimple: false,
	}, nil
}
