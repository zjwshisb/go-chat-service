package databases

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type Pagination struct {
	Data interface{} `json:"data"`
	Total int64 `json:"total"`
	Success bool `json:"success"`
}


func NewPagination(data interface{}, total int64) *Pagination {
	return &Pagination{
		Data: data,
		Total: total,
		Success: true,
	}
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