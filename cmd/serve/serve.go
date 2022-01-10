package serve

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"ws/app"
	"ws/app/cron"
	"ws/app/databases"
	"ws/app/file"
	mylog "ws/app/log"
	"ws/app/util"
)

func initCheck(cmd *cobra.Command, args []string) {
	if app.IsRunning() {
		log.Fatalln("serve is running")
	}
	workDir := util.GetWorkDir()
	if !util.DirExist(workDir) {
		panic(fmt.Sprintf("workdir:%s not exit", workDir))
	}
	storagePath := util.GetStoragePath()
	if !util.DirExist(storagePath) {
		err := util.MkDir(storagePath)
		if err != nil {
			panic(err)
		}
	}
}

func NewServeCommand() *cobra.Command {

	var cronFlag bool
	cmd := &cobra.Command{
		Use:                        "serve",
		Short: "start the server",
		PreRun: initCheck,
		Run: func(cmd *cobra.Command, args []string) {
			databases.MysqlSetup()
			databases.RedisSetup()
			file.Setup()
			mylog.Setup()
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
