package rule

import (
	"context"
	"fmt"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gvalid"
)

func init() {
	name := "exists"
	gvalid.RegisterRule(name, existsRule)
}

// ExistsRule 数据存在验证规则
func existsRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	params := parseRuleParams(in.Rule)
	length := len(params)
	if length < 1 {
		panic("exists rule must set table name params")
	}
	tableName := params[0]
	field := "id"
	if length >= 2 {
		field = params[1]
	}
	query := g.Model(tableName).Where(field, in.Value.String())
	customerId := service.AdminCtx().GetCustomerId(ctx)
	if customerId > 0 {
		query = query.Where("customer_id", customerId)
	}
	exists, err := query.Exist()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	message := fmt.Sprintf("%s had no row for %s:%s", tableName, field, in.Value.String())
	if in.Message != "" {
		message = in.Message
		return gerror.New(in.Message)
	}
	return gerror.Newf(message)
}
