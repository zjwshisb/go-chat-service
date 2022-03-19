package rpc

import (
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/spf13/viper"
	"log"
	"time"
	"ws/app/rpc/service/connection"
)

func addRegistryPlugin(s *server.Server) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@127.0.0.1:" + viper.GetString("Rpc.Port"),
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

func Serve() {
	if viper.GetBool("Rpc.Open") && viper.GetBool("App.Cluster") {
		s := server.NewServer()
		addRegistryPlugin(s)
		err := s.RegisterName("Connection", new(connection.Connection), "")
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			s.Serve("tcp", "127.0.0.1:"+viper.GetString("Rpc.port"))
		}()
	}
}
