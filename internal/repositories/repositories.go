package repositories

import (
	"gorm.io/gorm"
)

type Where struct {
	Filed string
	Value interface{}
}

func AddWhere(wheres []*Where) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, where := range wheres {
			db = db.Where(where.Filed, where.Value)
		}
		return db
	}
}



