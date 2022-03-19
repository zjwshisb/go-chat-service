package serve

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ws/app"
	"ws/app/cron"
	"ws/app/databases"
	"ws/app/file"
	_ "ws/app/http/requests"
	mylog "ws/app/log"
	"ws/app/routers"
	"ws/app/rpc"
	"ws/app/util"
	"ws/app/websocket"
	"ws/config"
)

func initCheck(cmd *cobra.Command, args []string) {
	if app.IsRunning() {
		log.Fatalln("serve is running")
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
			go rpc.Serve()
			websocket.SetupAdmin()
			websocket.SetupUser()
			if cronFlag {
				go cron.Run()
			}
			routers.Setup()
			srv := &http.Server{
				Addr:    viper.GetString("Http.Host") + ":" + viper.GetString("Http.Port"),
				Handler: routers.Router,
			}
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			go func() {
				err := srv.ListenAndServe()
				if err != nil {
					if err != http.ErrServerClosed {
						quit <- syscall.SIGINT
						log.Fatalln(err)
					}
				}
			}()
			defer func() {
				websocket.AdminManager.Destroy()
				websocket.UserManager.Destroy()
				cron.Stop()
			}()
			app.LogPid()
			<-quit
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer func() {
				cancel()
			}()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("Server Shutdown error :", err)
			}
			fmt.Println("exit forced")
		},
	}
	flag := cmd.Flags()
	flag.BoolVar(&cronFlag, "cron", true, "run cron or not")
	return cmd
}
