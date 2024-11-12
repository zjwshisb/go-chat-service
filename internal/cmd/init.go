package cmd

import (
	"context"
	"gf-chat/internal/service"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gcmd"
)

var Init = &gcmd.Command{
	Name:        "init",
	Brief:       "init database",
	Description: "初始化客服数据库",
	Arguments: []gcmd.Argument{
		{
			Name:   "customerId",
			Short:  "c",
			Brief:  "客户id",
			Orphan: true,
		},
	},
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		customerId := parser.GetOpt("c", 1)
		service.Setup().Setup(ctx, customerId.Uint())
		return nil
	},
}
