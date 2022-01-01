package serve

import (
	"fmt"
	"github.com/spf13/cobra"
	"ws/app"
	"ws/app/cron"
	"ws/app/databases"
	"ws/app/log"
)

func NewServeCommand() *cobra.Command {

	var cronFlag bool

	cmd := &cobra.Command{
		Use:                        "serve",
		Short: "start the server",
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
		CompletionOptions:          cobra.CompletionOptions{},
		SuggestionsMinimumDistance: 0,
		Run: func(cmd *cobra.Command, args []string) {
			databases.MysqlSetup()
			databases.RedisSetup()
			log.Setup()
			fmt.Println(cronFlag)
			if cronFlag {
				go cron.Run()
			}
			app.Start()
		},
	}
	flag := cmd.Flags()
	flag.BoolVar(&cronFlag, "cron", true, "run cron or not")
	return cmd
}
