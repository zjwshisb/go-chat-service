package rpc

import (
	"fmt"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/spf13/viper"
	"log"
	"os"
	"syscall"
	"time"
	mlog "ws/app/log"
	"ws/app/rpc/service"
	"ws/app/util"
)

func addRegistryPlugin(s *server.Server) {
	ips := util.GetIPs()

	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + ips[0] + ":" + viper.GetString("Rpc.Port"),
		EtcdServers:    viper.GetStringSlice("Etcd.Host"),
		BasePath:       viper.GetString("Etcd.BasePath"),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}

func Serve(c chan os.Signal) *server.Server {
	s := server.NewServer()
	addRegistryPlugin(s)
	err := s.RegisterName("Connection", new(service.Connection), "")
	if err != nil {
		log.Fatal(err)
	}
	err = s.RegisterName("Message", new(service.Message), "")
	if err != nil {
		log.Fatal(err)
	}
	err = s.RegisterName("Admin", new(service.Admin), "")
	if err != nil {
		log.Fatal(err)
	}
	err = s.RegisterName("User", new(service.User), "")

	go func() {
		mlog.Log.WithField("a-type", "rpc").Info("start")
		err := s.Serve("tcp", ":"+viper.GetString("Rpc.port"))
		if err != nil {
			c <- syscall.SIGINT
			fmt.Println(err)
		}
	}()
	return s
}
