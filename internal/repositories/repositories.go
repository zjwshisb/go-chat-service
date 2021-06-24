package repositories

import "gorm.io/gorm"

type Where struct {
	Filed string
	Value interface{}
}

type Pagination struct {
	Data interface{} `json:"data"`
	Total int64 `json:"total"`
	Success bool `json:"success"`
}

func Filter(wheres []*Where) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, where := range wheres {
			db = db.Where(where.Filed, where.Value)
		}
		return db
	}
}
func Paginate(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var offset int
		if limit <= 0 {
			limit = 20
		}
		if page <= 0 {
			page = 1
		}
		offset = (page - 1) * limit
		return db.Offset(offset)
	}
}

func NewPagination(data interface{}, total int64) *Pagination {
	return &Pagination{
		 Data: data,
		 Total: total,
		 Success: true,
	}
}

