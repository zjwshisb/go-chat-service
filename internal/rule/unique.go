package rule

import (
	"context"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gvalid"
)

func init() {
	name := "unique"
	gvalid.RegisterRule(name, UniqueRule)
}

// UniqueRule 唯一验证规则
func UniqueRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	name := in.Rule
	paramsStr := name[7:]
	params := gstr.Explode(",", paramsStr)
	length := len(params)
	if !(length == 2 || length == 3) {
		return gerror.New("unique rule need two/three params")
	}
	tableName := params[0]
	field := params[1]
	query := g.Model(tableName).Where(field, in.Value.String())
	if length == 3 {
		primaryField := params[2]
		request := ghttp.RequestFromCtx(ctx)
		v := request.GetRouter(primaryField).Val()
		query = query.WhereNot(primaryField, v)
	}
	customerId := service.AdminCtx().GetCustomerId(ctx)
	if customerId > 0 {
		query = query.Where("customer_id", service.AdminCtx().GetCustomerId(ctx))
	}
	count, err := query.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	if in.Message != "" {
		return gerror.New(in.Message)
	}
	return gerror.New("name already token by others")
}
