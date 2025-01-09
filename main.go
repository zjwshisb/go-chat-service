package main

import (
	"gf-chat/internal/cmd"
	_ "gf-chat/internal/cron"
	_ "gf-chat/internal/logic"
	_ "gf-chat/internal/packed"
	"log"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	err := cmd.Main.AddCommand(cmd.Http, cmd.Init, cmd.Migrate, cmd.Fake)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Main.Run(gctx.New())
}
