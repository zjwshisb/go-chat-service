package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var Migrate = &gcmd.Command{
	Name:        "migrate",
	Brief:       "migrate database",
	Description: "migrate database",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		sql, err := os.ReadFile(fmt.Sprintf("%s%s%s", pwd, string(os.PathSeparator), "database.sql"))
		if err != nil {
			panic(err)
		}
		_, err := g.DB().GetCore().Exec(ctx, string(sql))
		if err != nil {
			panic(err)
		}
		fmt.Printf(result)
		return err
	},
}
