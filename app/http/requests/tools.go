package requests

import (
	"github.com/gin-gonic/gin"
	"ws/app/repositories"
)

func GetFilterWheres(c *gin.Context, fields []string) []*repositories.Where {
	wheres := make([]*repositories.Where, 0, len(fields))
	for _, field := range fields {
		if value, exist := c.GetQuery(field); exist {
			if value != "" {
				wheres = append(wheres, &repositories.Where{
					Filed: field + " = ?",
					Value: value,
				})
			}
		}
	}
	return wheres
}
