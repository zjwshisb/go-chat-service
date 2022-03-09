package client

import (
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	xclient "github.com/smallnest/rpcx/client"
	"github.com/spf13/viper"
	"log"
)

func NewClient(servicePath string) xclient.ServiceDiscovery {
	var d xclient.ServiceDiscovery
	var err error
	if viper.GetBool("App.Cluster") {
		d, err = etcdClient.NewEtcdV3Discovery(
			viper.GetString("Etcd.BasePath"),
			servicePath,
			viper.GetStringSlice("Etcd.Host"),
			false,
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		d, _ = xclient.NewPeer2PeerDiscovery("127.0.0.1:" + viper.GetString("Rpc.Port"), "")
	}
	return d
}
