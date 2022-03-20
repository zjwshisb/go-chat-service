package serve

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ws/app/cron"
	"ws/app/databases"
	"ws/app/file"
	"ws/app/http"
	_ "ws/app/http/requests"
	"ws/app/http/websocket"
	mylog "ws/app/log"
	"ws/app/rpc"
	"ws/app/sys"
	"ws/app/util"
	"ws/config"
)

func initCheck(cmd *cobra.Command, args []string) {
	if sys.IsRunning() {
		log.Fatalln("service is running")
	}
	workDir := config.GetWorkDir()
	if !util.DirExist(workDir) {
		panic(fmt.Sprintf("workdir:%s not exit", workDir))
	}
	storagePath := config.GetStoragePath()
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
		Use:    "serve",
		Short:  "start the server",
		PreRun: initCheck,
		Run: func(cmd *cobra.Command, args []string) {
			databases.MysqlSetup()
			databases.RedisSetup()
			file.Setup()
			mylog.Setup()
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			htp := http.Serve(quit)
			if cronFlag {
				cn := cron.Serve()
				defer func() {
					cn.Stop()
				}()
			}
			if viper.GetBool("Rpc.open") {
				rc := rpc.Serve(quit)
				defer func() {
					_ = rc.Close()
				}()
			}
			defer func() {
				websocket.AdminManager.Destroy()
				websocket.UserManager.Destroy()
			}()
			sys.LogPid()
			<-quit
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer func() {
				cancel()
			}()
			if err := htp.Shutdown(ctx); err != nil {
				log.Fatal("Server Shutdown error :", err)
			}
			fmt.Println("exit forced")
		},
	}
	flag := cmd.Flags()
	flag.BoolVar(&cronFlag, "cron", true, "run cron or not")
	return cmd
}
