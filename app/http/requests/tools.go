package requests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"ws/app/repositories"
)


func GetFilterWhere(c *gin.Context , fields map[string]interface{}) []*repositories.Where {
	wheres := make([]*repositories.Where, 0, len(fields))
	for field, handler := range fields {
		if val, exist := c.GetQuery(field); exist && val != "" {
			switch reflect.TypeOf(handler).Kind() {
			case reflect.Func:
				params := make([]reflect.Value, 1, 1)
				f := reflect.ValueOf(handler)
				params[0] = reflect.ValueOf(val)
				result := f.Call(params)
				for _, s := range result {
					switch s.Kind() {
					case reflect.Interface:
						switch reflect.TypeOf(s.Interface()).Kind() {
						case reflect.Slice:
							slice := reflect.ValueOf(s.Interface())
							for i:= 0; i < slice.Len(); i++ {
								v := slice.Index(i)
								where, ok := v.Interface().(*repositories.Where)
								if ok {
									wheres = append(wheres, where)
								}
							}
						case reflect.Ptr:
							i := s.Interface()
							where ,ok := i.(*repositories.Where)
							if ok {
								wheres = append(wheres, where)
							}
						}
					case reflect.Slice:
						slice := reflect.ValueOf(s.Interface())
						for i:= 0; i < slice.Len(); i++ {
							v := slice.Index(i)
							where, ok := v.Interface().(*repositories.Where)
							if ok {
								wheres = append(wheres, where)
							}
						}
					case reflect.Ptr:
						i := s.Interface()
						where ,ok := i.(*repositories.Where)
						if ok {
							wheres = append(wheres, where)
						}
					}
				}
			case reflect.String:
				operator := reflect.ValueOf(handler).String()
				if operator == "" {
					operator = "="
				}
				wheres = append(wheres, &repositories.Where{
					Filed: fmt.Sprintf("%s %s ?", field, operator),
					Value: val,
				})
			}
		}
	}
	return wheres
}
