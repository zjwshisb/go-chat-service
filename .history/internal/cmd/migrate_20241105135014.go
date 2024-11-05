package cmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gcmd"
)

var Migrate = &gcmd.Command{
	Name:        "migrate",
	Brief:       "migrate database",
	Description: "migrate database",
	Func: func(ctx context.Context, parser *gcmd.Parser) err {
		return nil
	},
}
