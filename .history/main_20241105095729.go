package main

import (
	"gf-chat/internal/cmd"
	_ "gf-chat/internal/cron"
	_ "gf-chat/internal/logic"
	_ "gf-chat/internal/packed"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	cmd.Main.AddCommand(cmd.Http, cmd.Init)
	cmd.Main.Run(gctx.New())
}
