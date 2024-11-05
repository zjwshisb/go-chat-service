package cmd

import (
	"context"
	"fmt"
	"os"

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
		sql, err := os.Open(pwd + (string) os.PathSeparator + "database.sql")
		fmt.Println(pwd)
		return nil
	},
}
