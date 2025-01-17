package rule

import (
	"context"
	"fmt"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"
	"strings"
)

func init() {
	gvalid.RegisterRule("unique", uniqueRule)
}

// UniqueRule unique:users[,name][,id] 唯一验证规则
func uniqueRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	request := ghttp.RequestFromCtx(ctx)
	params := parseRuleParams(in.Rule)
	length := len(params)
	if length < 1 {
		panic("unsupported used for unique rule")
	}
	tableName := params[0]
	field := strings.ToLower(in.Field)
	if length >= 2 {
		field = params[1]
	}
	primaryKey := "id"
	if length >= 3 {
		primaryKey = params[2]
	}
	query := g.Model(tableName).Where(field, in.Value.String())
	if request.Method == "PUT" {
		v := request.GetRouter(primaryKey).Val()
		query = query.WhereNot(primaryKey, v)
	}
	customerId := service.AdminCtx().GetCustomerId(ctx)
	if customerId > 0 {
		query = query.Where("customer_id", customerId)
	}
	count, err := query.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	message := fmt.Sprintf("%s(%s) in %s is duplicate", field, in.Value.String(), tableName)
	if in.Message != "" {
		message = in.Message
	}
	return gerror.NewCode(gcode.CodeValidationFailed, message)
}
