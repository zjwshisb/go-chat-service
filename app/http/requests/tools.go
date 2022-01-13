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
			handlerType := reflect.TypeOf(handler)
			switch handlerType.Kind() {
			case reflect.Func:
				numIn := handlerType.NumIn()
				if numIn == 1 {
					params := make([]reflect.Value, 1, 1)
					f := reflect.ValueOf(handler)
					params[0] = reflect.ValueOf(val)
					result := f.Call(params)
					for _, s := range result {
						if !s.IsNil() {
							switch s.Kind() {
							case reflect.Interface:
								ii := s.Interface()
								switch ii.(type) {
								case *repositories.Where:
									w ,ok := ii.(*repositories.Where)
									if ok {
										wheres = append(wheres, w)
									}
								case []*repositories.Where:
									ws ,ok := ii.([]*repositories.Where)
									if ok {
										wheres = append(wheres, ws...)
									}
								}
							case reflect.Slice:
								for i:= 0; i < s.Len(); i++ {
									v := s.Index(i)
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
