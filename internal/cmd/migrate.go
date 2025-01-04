package cmd

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/strutil"
	"os"
	"strings"

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
		sqls, err := os.ReadFile(fmt.Sprintf("%s%s%s", pwd, string(os.PathSeparator), "database.sql"))
		if err != nil {
			panic(err)
		}
		sqlArr := strings.Split(string(sqls), ";")
		tx, err := g.DB().Begin(ctx)
		if err != nil {
			panic(err)
		}
		for _, sql := range sqlArr {
			if strutil.Trim(sql) != "" {
				_, err = tx.Exec(sql)
				if err != nil {
					txErr := tx.Rollback()
					if txErr != nil {
						panic(txErr)
					}
					panic(err)
				}
			}
		}
		err = tx.Commit()
		return err
	},
}
