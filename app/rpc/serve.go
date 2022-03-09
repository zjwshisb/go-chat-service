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
		ServiceAddress: "tcp@0.0.0.0:" + viper.GetString("Rpc.Port"),
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
	if viper.GetBool("Rpc.Open") {
		s := server.NewServer()
		addRegistryPlugin(s)
		s.RegisterName("Connection", new(connection.Connection), "")
		go s.Serve("tcp", "0.0.0.0:" + viper.GetString("Rpc.port"))
	}
}
