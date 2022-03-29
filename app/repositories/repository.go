package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"ws/app/databases"
	"ws/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[T any] struct {
}

// DeleteAll delete models by wheres
func (Repository *Repository[T]) DeleteAll(wheres []*Where) int64 {
	result := databases.Db.Model(new(T)).Scopes(AddWhere(wheres)).Delete(new(T))
	return result.RowsAffected
}

// Delete delete model
func (Repository *Repository[T]) Delete(model *T) int64 {
	result := databases.Db.Delete(model)
	return result.RowsAffected
}

// Update update models
func (Repository *Repository[T]) Update(wheres []*Where, values map[string]interface{}) int64 {
	result := databases.Db.Model(new(T)).Scopes(AddWhere(wheres)).Updates(values)
	return result.RowsAffected
}

// Save create or update model
func (Repository *Repository[T]) Save(model *T) error {
	query := databases.Db.Omit(clause.Associations).Save(model)
	return query.Error
}

// Get get models
func (Repository *Repository[T]) Get(wheres []*Where, limit int, load []string, orders []string) []*T {
	lists := make([]*T, 0)
	databases.Db.Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(orders)).Limit(limit).Find(&lists)
	return lists
}

// First get first model
func (Repository *Repository[T]) First(wheres []*Where, orders []string) *T {
	model := new(T)
	result := databases.Db.Scopes(AddOrder(orders)).Scopes(AddWhere(wheres)).First(model)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return model
}

// FirstById get first model by primary key
func (Repository *Repository[T]) FirstById(primaryKey interface{}) *T {
	model := new(T)
	result := databases.Db.First(model, primaryKey)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return model
}

// Paginate get modes pagination
func (Repository *Repository[T]) Paginate(c *gin.Context, wheres []*Where, load []string, order []string) *Pagination[*T] {
	items := make([]*T, 0)
	databases.Db.
		Scopes(Paginate(c)).
		Scopes(AddWhere(wheres)).
		Scopes(AddLoad(load)).
		Scopes(AddOrder(order)).
		Find(&items)
	var total int64
	databases.Db.Model(&models.Admin{}).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return NewPagination(items, total)
}

type Where struct {
	Filed string
	Value interface{}
}

func Filter(c *gin.Context, fields []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, field := range fields {
			if value, exist := c.GetQuery(field); exist {
				db.Where(fmt.Sprintf("%s = ?", field), value)
			}
		}
		return db
	}
}

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var offset int
		limit := 20
		page := 1
		limitStr, ok := c.GetQuery("pageSize")
		if ok {
			i, err := strconv.Atoi(limitStr)
			if err == nil {
				limit = i
			}
		}
		pageStr, ok := c.GetQuery("current")
		if ok {
			i, err := strconv.Atoi(pageStr)
			if err == nil {
				page = i
			}
		}
		offset = (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func AddLoad(relations []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, load := range relations {
			db = db.Preload(load)
		}
		return db
	}
}
func AddOrder(orders []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, load := range orders {
			db = db.Order(load)
		}
		return db
	}
}
func AddWhere(wheres []*Where) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, where := range wheres {
			db = db.Where(where.Filed, where.Value)
		}
		return db
	}
}
